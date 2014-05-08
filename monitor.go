package main

import (
	"time"
)

func MonitorStart(queries chan QueryFile, metrics chan []interface{}, interval int) {
	Log("monitor.start")
	for {
		Log("monitor.tick queries=%d metrics=%d", len(queries), len(metrics))
		<-time.After(time.Duration(interval) * time.Second)
	}
}
