package main

func main() {
	Log("main.start")

	queryInterval := QueryInterval()
	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()
	queryFiles := ReadQueryFiles("./queries/*.sql")

	metricBatches := make(chan []interface{}, 10)
	queryTicks := make(chan QueryFile, 10)
	libratoStop := make(chan bool)
	postgresStop := make(chan bool)
	schedulerStop := make(chan bool)
	globalStop := make(chan bool)

	go LibratoStart(libratoAuth, metricBatches, libratoStop)
	go PostgresStart(databaseUrl, queryTicks, metricBatches, postgresStop)
	go SchedulerStart(queryFiles, queryInterval, queryTicks, schedulerStop)
	go TrapStart(globalStop)

	Log("main.await")
	stop := <-globalStop

	Log("main.stop")
	schedulerStop <- stop
	postgresStop <- stop
	libratoStop <- stop

	Log("main.exit")
}
