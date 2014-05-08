package main

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
