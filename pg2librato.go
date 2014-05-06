package main

func main() {
	Log("main.read-env")
	queryInterval := QueryInterval()
	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()

	Log("main.read-query-files")
	queryFiles, err := ReadQueryFiles("./queries/*.sql")
	if err != nil {
		panic(err)
	}

	Log("main.librato-start")
	metricBatches := make(chan []interface{})
	libratoStop := make(chan bool)
	go LibratoStart(libratoAuth, metricBatches, libratoStop)

	Log("main.postgres-start")
	queryTicks := make(chan QueryFile)
	postgresStop := make(chan bool)
	go PostgresStart(databaseUrl, queryTicks, metricBatches, postgresStop)

	Log("main.scheduler-start")
	schedulerStop := make(chan bool)
	go SchedulerStart(queryFiles, queryInterval, queryTicks, schedulerStop)

	Log("main.trap-start")
	globalStop := make(chan bool)
	go TrapStart(globalStop)

	Log("main.await")
	stop := <-globalStop
	Log("main.scheduler-stop")
	schedulerStop <- stop
	Log("main.postgres-stop")
	postgresStop <- stop
	Log("main.librato-stop")
	libratoStop <- stop
	Log("main.exit")
}
