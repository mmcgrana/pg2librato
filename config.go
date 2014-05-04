package main

import (
	"os"
	"strconv"
)

func MustGetenv(k string) string {
	s := os.Getenv(k)
	if s == "" {
		panic("Must set " + k)
	}
	return s
}

func DatabaseUrl() string {
	return MustGetenv("DATABASE_URL")
}

func QueryInterval() int {
	s := MustGetenv("QUERY_INTERVAL")
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	if i <= 0 {
		panic("Must provide QUERY_INTERVAL > 0")
	}
	return i
}

func RollbarToken() string {
	return os.Getenv(k)
}
