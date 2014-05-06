package main

func main() {
	Log("main.read-confg")
	queryInterval := QueryInterval()
	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()
	queryFiles := ReadQueryFiles("./queries/*.sql")

	Log("main.start")
	metricBatches := make(chan []interface{})
	libratoStop := make(chan bool)
	go LibratoStart(libratoAuth, metricBatches, libratoStop)

	queryTicks := make(chan QueryFile)
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
