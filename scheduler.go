package main

import (
	"time"
)

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
