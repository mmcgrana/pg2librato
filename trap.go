package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func TrapStart(stop chan<- bool) {
	Log("trap.start")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	Log("trap.listening")
	sig := <-sigs
	Log(fmt.Sprintf("trap.caught signal=%s", sig.String()))
	stop <- true
}
