package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samuel/go-librato/librato"
)

func postgresQuery(db *sql.DB, qf QueryFile, queryTimeout int) []interface{} {
	Log("postgres.query.start name=%s", qf.Name)
	_, err := db.Exec(fmt.Sprintf("set application_name TO 'pg2librato - %s'", qf.Name))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(fmt.Sprintf("set statement_timeout TO %d", queryTimeout*1000))
	if err != nil {
		panic(err)
	}

	rows, err := db.Query(qf.Sql)
	if err != nil {
		panic(err)
	}
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	numCols := len(cols)
	if numCols != 3 {
		panic("Must return result set with exactly 3 rows")
	}
	Log("postgres.query.finish name=%s", qf.Name)
	metrics := []interface{}{}
	for rows.Next() {
		var name string
		var nullSource sql.NullString
		var source string
		var value float64
		err = rows.Scan(&name, &nullSource, &value)
		if err != nil {
			panic(err)
		}
		if nullSource.Valid {
			source = nullSource.String
		}
		Log("postgres.result name=%s source=%s value=%f", name, source, value)
		metric := librato.Metric{
			Name:   name,
			Source: source,
			Value:  value,
		}
		metrics = append(metrics, metric)
	}
	return metrics
}

func postgresWorkerStart(db *sql.DB, queryTicks <-chan QueryFile, queryTimeout int, metricBatches chan<- []interface{}, stop <-chan bool, done chan<- bool) {
	Log("postgres.worker.start")
	for {
		select {
		case queryFile := <-queryTicks:
			metricBatch := postgresQuery(db, queryFile, queryTimeout)
			metricBatches <- metricBatch
		default:
			select {
			case <-stop:
				Log("postgres.worker.exit")
				done <- true
				return
			default:
			}
		}
	}
}

func PostgresStart(databaseUrl string, queryTicks <-chan QueryFile, queryTimeout int, metricBatches chan<- []interface{}, stop <-chan bool, done chan<- bool) {
	Log("postgres.start")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		panic(err)
	}

	postgresWorkerStops := make([]chan bool, PostgresWorkers)
	postgresWorkerDones := make([]chan bool, PostgresWorkers)
	for w := 0; w < PostgresWorkers; w++ {
		postgresWorkerStops[w] = make(chan bool)
		postgresWorkerDones[w] = make(chan bool)
		go postgresWorkerStart(db, queryTicks, queryTimeout, metricBatches, postgresWorkerStops[w], postgresWorkerDones[w])
	}

	<-stop
	Log("postgres.stop")
	for w := 0; w < PostgresWorkers; w++ {
		postgresWorkerStops[w] <- true
		<-postgresWorkerDones[w]
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}
	done <- true

	Log("postgres.exit")
}
