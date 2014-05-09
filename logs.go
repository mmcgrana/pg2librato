package main

import (
	"log"
)

func Log(l string, t ...interface{}) {
	log.Printf(l, t...)
}
