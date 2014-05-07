package main

import (
	"time"
)

func MonitorStart(queries chan QueryFile, metrics chan []interface{}, interval int, stop chan bool) {
	Log("monitor.start")
	for {
		Log("monitor.tick queries=%d metrics=%d", len(queries), len(metrics))
		select {
		case <-stop:
			Log("monitor.exit")
			stop <- true
			return
		case <-time.After(time.Duration(interval) * time.Second):
		}
	}
}
