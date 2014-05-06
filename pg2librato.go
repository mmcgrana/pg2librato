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
	monitorStop := make(chan bool)
	libratoStop := make(chan bool)
	postgresStop := make(chan bool)
	schedulerStop := make(chan bool)
	done := make(chan bool)

	go TrapStart(globalStop)
	go MonitorStart(queryTicks, metricBatches, monitorStop, done)
	go LibratoStart(libratoAuth, metricBatches, libratoStop, done)
	go PostgresStart(databaseUrl, queryTicks, metricBatches, postgresStop, done)
	go SchedulerStart(queryFiles, queryInterval, queryTicks, schedulerStop, done)

	Log("main.await")
	<-globalStop

	Log("main.stop")
	schedulerStop <- true
	<-done
	postgresStop <- true
	<-done
	libratoStop <- true
	<-done
	monitorStop <- true
	<-done

	Log("main.exit")
}
