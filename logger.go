package main

import (
	"fmt"
	"log"
)

func Log(l string, t ...interface{}) {
	log.Println(fmt.Sprintf(l, t...))
}
