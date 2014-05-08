package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samuel/go-librato/librato"
	"time"
)

func main() {
	Log("main.start")

	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()
	queryInterval := QueryInterval()
	queryFiles := QueryFiles()

	metricBatches := make(chan []interface{}, MetricsBufferSize)
	queryTicks := make(chan QueryFile, QueriesBufferSize)

	go monitorStart(queryTicks, metricBatches, MonitorInterval)
	go schedulerStart(queryFiles, queryInterval, queryTicks)
	go libratoStart(libratoAuth, metricBatches)
	go postgresStart(databaseUrl, queryTicks, QueryTimeout, metricBatches)

	<-make(chan bool)
}

func monitorStart(queries chan QueryFile, metrics chan []interface{}, interval int) {
	Log("monitor.start")
	for {
		Log("monitor.tick queries=%d metrics=%d", len(queries), len(metrics))
		<-time.After(time.Duration(interval) * time.Second)
	}
}

func schedulerStart(queryFiles []QueryFile, queryInterval int, queryTicks chan<- QueryFile) {
	Log("scheduler.start")
	for {
		Log("scheduler.tick")
		for _, queryFile := range queryFiles {
			select {
			case queryTicks <- queryFile:
			default:
				Error(errors.New("Queries channel full"))
			}
		}
		<-time.After(time.Duration(queryInterval) * time.Second)
	}
}

func libratoStart(libratoAuth []string, metricBatches <-chan []interface{}) {
	Log("librato.start")
	lb := &librato.Client{libratoAuth[0], libratoAuth[1]}
	for {
		metricBatch := <-metricBatches
		Log("librato.post.start")
		err := lb.PostMetrics(&librato.Metrics{
			Gauges: metricBatch,
		})
		if err != nil {
			Error(err)
		}
		Log("librato.post.finish")
	}
}

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

func postgresStart(databaseUrl string, queryTicks <-chan QueryFile, queryTimeout int, metricBatches chan<- []interface{}) {
	Log("postgres.start")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		Error(err)
	}
	for {
		queryFile := <-queryTicks
		metricBatch, err := postgresQuery(db, queryFile, queryTimeout)
		if err != nil {
			Error(err)
		}
		select {
		case metricBatches <- metricBatch:
		default:
			Error(errors.New("Metrics channel full"))
		}
	}
}
