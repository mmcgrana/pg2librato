package main

func main() {
	Log("main.start")

	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()
	queryInterval := QueryInterval()
	queryTimeout := queryInterval
	queryFiles := ReadQueryFiles("./queries/*.sql")

	metricBatches := make(chan []interface{}, 10)
	queryTicks := make(chan QueryFile, 10)
	globalStop := make(chan bool)
	monitorStop := make(chan bool)
	libratoStop := make(chan bool)
	postgresStop := make(chan bool)
	schedulerStop := make(chan bool)

	go TrapStart(globalStop)
	go MonitorStart(queryTicks, metricBatches, MonitorInterval, monitorStop)
	go LibratoStart(libratoAuth, metricBatches, libratoStop)
	go PostgresStart(databaseUrl, queryTicks, queryTimeout, metricBatches, postgresStop)
	go SchedulerStart(queryFiles, queryInterval, queryTicks, schedulerStop)

	Log("main.await")
	<-globalStop

	Log("main.stop")
	schedulerStop <- true
	<-schedulerStop
	postgresStop <- true
	<-postgresStop
	libratoStop <- true
	<-libratoStop
	monitorStop <- true
	<-monitorStop

	Log("main.exit")
}
