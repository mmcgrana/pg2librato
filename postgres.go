package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samuel/go-librato/librato"
)

func postgresQuery(db *sql.DB, qf QueryFile, queryTimeout int) ([]interface{}, error) {
	Log("postgres.query.start name=%s", qf.Name)
	_, err := db.Exec(fmt.Sprintf("set application_name TO 'pg2librato - %s'", qf.Name))
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(fmt.Sprintf("set statement_timeout TO %d", queryTimeout*1000))
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(qf.Sql)
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	numCols := len(cols)
	if numCols != 3 {
		return nil, errors.New("Must return result set with exactly 3 rows")
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
			return nil, err
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
	return metrics, nil
}

func postgresWorkerStart(db *sql.DB, queryTicks <-chan QueryFile, queryTimeout int, metricBatches chan<- []interface{}, stop chan bool) {
	Log("postgres.worker.start")
	for {
		select {
		case queryFile := <-queryTicks:
			metricBatch, err := postgresQuery(db, queryFile, queryTimeout)
			if err != nil {
				Error(err)
			}
			metricBatches <- metricBatch
		default:
			select {
			case <-stop:
				Log("postgres.worker.exit")
				stop <- true
				return
			default:
			}
		}
	}
}

func PostgresStart(databaseUrl string, queryTicks <-chan QueryFile, queryTimeout int, metricBatches chan<- []interface{}, stop chan bool) {
	Log("postgres.start")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		Error(err)
	}

	postgresWorkerStops := make([]chan bool, PostgresWorkers)
	for w := 0; w < PostgresWorkers; w++ {
		postgresWorkerStops[w] = make(chan bool)
		go postgresWorkerStart(db, queryTicks, queryTimeout, metricBatches, postgresWorkerStops[w])
	}

	<-stop
	Log("postgres.stop")
	for w := 0; w < PostgresWorkers; w++ {
		postgresWorkerStops[w] <- true
		<-postgresWorkerStops[w]
	}
	err = db.Close()
	if err != nil {
		Error(err)
	}
	stop <- true

	Log("postgres.exit")
}
