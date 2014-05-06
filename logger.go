package main

import (
	"fmt"
)

func Log(l string, t ...interface{}) {
	fmt.Printf(l+"\n", t...)
}
