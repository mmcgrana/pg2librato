package main

import (
	"time"
)

func MonitorStart(queryTicks chan QueryFile, metricBatches chan []interface{}, stop <-chan bool, done chan<- bool) {
	Log("monitor.start")
	for {
		Log("monitor.tick queries=%d metrics=%d", len(queryTicks), len(metricBatches))
		select {
		case <-stop:
			Log("monitor.exit")
			done <- true
			return
		case <-time.After(MonitorInterval * time.Second):
		}
	}
}
