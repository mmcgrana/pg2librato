package main

import (
	"time"
)

func SchedulerStart(queryFiles []QueryFile, queryInterval int, queryTicks chan<- QueryFile, stop <-chan bool) {
	Log("scheduler.start")
	for {
		Log("scheduler.tick")
		for _, queryFile := range queryFiles {
			queryTicks <- queryFile
		}

		select {
		case <-stop:
			Log("scheduler.stop")
			return
		case <-time.After(time.Duration(queryInterval) * time.Second):
		}
	}
}
