package main

import (
	"time"
)

func main() {
	Log("main.start")

	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()
	queryInterval := QueryInterval()
	queryFiles := QueryFiles()

	metricBatches := make(chan []interface{}, 10)
	queryTicks := make(chan QueryFile, 10)

	go MonitorStart(queryTicks, metricBatches, MonitorInterval)
	go LibratoStart(libratoAuth, metricBatches)
	go PostgresStart(databaseUrl, queryTicks, QueryTimeout, metricBatches)
	go SchedulerStart(queryFiles, queryInterval, queryTicks)

	<-make(chan bool)
}

func MonitorStart(queries chan QueryFile, metrics chan []interface{}, interval int) {
	Log("monitor.start")
	for {
		Log("monitor.tick queries=%d metrics=%d", len(queries), len(metrics))
		<-time.After(time.Duration(interval) * time.Second)
	}
}

func SchedulerStart(queryFiles []QueryFile, queryInterval int, queryTicks chan<- QueryFile) {
	Log("scheduler.start")
	for {
		Log("scheduler.tick")
		for _, queryFile := range queryFiles {
			queryTicks <- queryFile
		}
		<-time.After(time.Duration(queryInterval) * time.Second)
	}
}
