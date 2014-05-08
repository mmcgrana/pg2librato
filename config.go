package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	MonitorInterval   = 5
	QueryTimeout      = 30
	MetricsBufferSize = 10
	QueriesBufferSize = 10
)

func MustGetenv(k string) string {
	s := os.Getenv(k)
	if s == "" {
		Error(errors.New("Must set " + k))
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
		Error(err)
	}
	if i <= 0 {
		Error(errors.New("Must provide QUERY_INTERVAL > 0"))
	}
	return i
}

func LibratoAuth() []string {
	s := MustGetenv("LIBRATO_AUTH")
	a := strings.Split(s, ":")
	if len(a) != 2 {
		Error(errors.New("Must provide LIBRATO_AUTH as email:token"))
	}
	return a
}

func RollbarAccessToken() string {
	return os.Getenv("ROLLBAR_ACCESS_TOKEN")
}

func RollbarEnvironment() string {
	return MustGetenv("ROLLBAR_ENVIRONMENT")
}

type QueryFile struct {
	Name string
	Sql  string
}

func QueryFiles() []QueryFile {
	sqlPaths, err := filepath.Glob("./queries/*.sql")
	if err != nil {
		Error(err)
	}
	queryFiles := make([]QueryFile, len(sqlPaths))
	for i, path := range sqlPaths {
		sqlBytes, err := ioutil.ReadFile(path)
		if err != nil {
			Error(err)
		}
		pathBase := filepath.Base(path)
		queryFiles[i] = QueryFile{
			Name: pathBase,
			Sql:  string(sqlBytes),
		}
	}
	return queryFiles
}
