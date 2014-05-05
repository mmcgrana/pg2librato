package main

import (
	"os"
	"strconv"
	"strings"
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

func LibratoAuth() []string {
	s := MustGetenv("LIBRATO_AUTH")
	a := strings.Split(s, ":")
	if len(a) != 2 {
		panic("Must provide LIBRATO_AUTH as email:token")
	}
	return a
}

func RollbarToken() string {
	return os.Getenv("ROLLBAR_TOKEN")
}
