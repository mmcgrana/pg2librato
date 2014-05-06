package main

func main() {
	Log("main.start")

	queryInterval := QueryInterval()
	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()
	queryFiles := ReadQueryFiles("./queries/*.sql")

	metricBatches := make(chan []interface{}, 10)
	queryTicks := make(chan QueryFile, 10)
	globalStop := make(chan bool)
	libratoStop := make(chan bool)
	libratoDone := make(chan bool)
	postgresStop := make(chan bool)
	postgresDone := make(chan bool)
	schedulerStop := make(chan bool)
	schedulerDone := make(chan bool)

	go TrapStart(globalStop)
	go LibratoStart(libratoAuth, metricBatches, libratoStop, libratoDone)
	go PostgresStart(databaseUrl, queryTicks, metricBatches, postgresStop, postgresDone)
	go SchedulerStart(queryFiles, queryInterval, queryTicks, schedulerStop, schedulerDone)

	Log("main.await")
	<-globalStop

	Log("main.stop")
	schedulerStop <- true
	<-schedulerDone
	postgresStop <- true
	<-postgresDone
	libratoStop <- true
	<-libratoDone

	Log("main.exit")
}
