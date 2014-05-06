package main

func main() {
	Log("main.start")
	queryInterval := QueryInterval()
	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()
	queryFiles := ReadQueryFiles("./queries/*.sql")

	metricBatches := make(chan []interface{}, 10)
	libratoStop := make(chan bool)
	go LibratoStart(libratoAuth, metricBatches, libratoStop)

	queryTicks := make(chan QueryFile, 10)
	postgresStop := make(chan bool)
	go PostgresStart(databaseUrl, queryTicks, metricBatches, postgresStop)

	schedulerStop := make(chan bool)
	go SchedulerStart(queryFiles, queryInterval, queryTicks, schedulerStop)

	globalStop := make(chan bool)
	go TrapStart(globalStop)

	Log("main.await")
	stop := <-globalStop

	Log("main.stop")
	schedulerStop <- stop
	postgresStop <- stop
	libratoStop <- stop

	Log("main.exit")
}
