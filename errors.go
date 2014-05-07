package main

import (
	"github.com/stvp/rollbar"
)

func Error(e error) {
	rt := RollbarAccessToken()
	if rt != "" {
		Log("error.report.start")
		re := RollbarEnvironment()
		rollbar.Token = rt
		rollbar.Environment = re
		rollbar.Error("error", e)
		rollbar.Wait()
		Log("error.report.finish")
	}
	panic(e)
}
