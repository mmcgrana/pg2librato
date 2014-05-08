package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samuel/go-librato/librato"
)

func postgresPrep(db *sql.DB, timeout int, name string) error {
	_, err := db.Exec(fmt.Sprintf("set statement_timeout TO %d", timeout*1000))
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("set application_name TO '%s'", name))
	if err != nil {
		return err
	}
	return nil
}

func postgresScanMetric(rows *sql.Rows) (*librato.Metric, error) {
	var name string
	var nullSource sql.NullString
	var source string
	var value float64
	err := rows.Scan(&name, &nullSource, &value)
	if err != nil {
		return nil, err
	}
	if nullSource.Valid {
		source = nullSource.String
	}
	metric := &librato.Metric{
		Name:   name,
		Source: source,
		Value:  value,
	}
	return metric, nil
}

func postgresQuery(db *sql.DB, qf QueryFile, timeout int) ([]interface{}, error) {
	Log("postgres.query.start name=%s", qf.Name)
	err := postgresPrep(db, timeout, "pg2librato - "+qf.Name)
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
		metric, err := postgresScanMetric(rows)
		if err != nil {
			return nil, err
		}
		Log("postgres.result name=%s source=%s value=%f", metric.Name, metric.Source, metric.Value)
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func postgresWorkerStart(db *sql.DB, queryTicks <-chan QueryFile, queryTimeout int, metricBatches chan<- []interface{}) {
	Log("postgres.worker.start")
	for {
		queryFile := <-queryTicks
		metricBatch, err := postgresQuery(db, queryFile, queryTimeout)
		if err != nil {
			Error(err)
		}
		metricBatches <- metricBatch
	}
}

func PostgresStart(databaseUrl string, queryTicks <-chan QueryFile, queryTimeout int, metricBatches chan<- []interface{}) {
	Log("postgres.start")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		Error(err)
	}
	for w := 0; w < PostgresWorkers; w++ {
		go postgresWorkerStart(db, queryTicks, queryTimeout, metricBatches)
	}
}
