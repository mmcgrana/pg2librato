package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samuel/go-librato/librato"
)

func postgresConnect(databaseUrl string) *sql.DB {
	Log("postgres.connect.start")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT 1")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			panic(err)
		}
		if count != 1 {
			panic("Couldn't connect to database")
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	Log("postgres.connect.finish")
	return db
}

func postgresQuery(db *sql.DB, qf QueryFile) []interface{} {
	Log(fmt.Sprintf("postgres.query.start name=%s", qf.Name))
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
		Log(fmt.Sprintf("postgres.result name=%s source=%s value=%f", name, source, value))
		metric := librato.Metric{
			Name:   name,
			Source: source,
			Value:  value,
		}
		metrics = append(metrics, metric)
	}
	return metrics
}

func PostgresStart(databaseUrl string, queryTicks <-chan QueryFile, metricBatches chan<- []interface{}, stop <-chan bool) {
	Log("postgres.start")
	db := postgresConnect(databaseUrl)

	for {
		select {
		case <-stop:
			Log("postgres.stop")
			return
		default:
		}

		select {
		case queryFile := <-queryTicks:
			metricBatch := postgresQuery(db, queryFile)
			metricBatches <- metricBatch
		default:
		}
	}
}
