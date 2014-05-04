package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("starting")
	interval := QueryInterval()
	dbUrl := DatabaseUrl()
	for {
		fmt.Println("querying" + dbUrl)
		fmt.Println("reporting")
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
