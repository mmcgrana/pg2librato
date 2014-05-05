package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/samuel/go-librato/librato"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type Queryfile struct {
	Name string
	Sql  string
}

func main() {
	fmt.Println("config.read-env")
	queryInterval := QueryInterval()
	databaseUrl := DatabaseUrl()
	libratoAuth := LibratoAuth()

	fmt.Println("config.read-files")
	sqlPaths, err := filepath.Glob("./queries/*.sql")
	if err != nil {
		panic(err)
	}
	sqlQueryfiles := make([]Queryfile, len(sqlPaths))
	for i, path := range sqlPaths {
		sqlBytes, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		pathBase := filepath.Base(path)
		sqlQueryfiles[i] = Queryfile{
			Name: strings.TrimSuffix(pathBase, filepath.Ext(pathBase)),
			Sql:  string(sqlBytes),
		}
	}

	fmt.Println("postgres.connect")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SELECT 1")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			panic(err)
		}
		if count != 1 {
			panic("Couldn't connect to database")
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("librato.connect")
	lb := &librato.Client{libratoAuth[0], libratoAuth[1]}

	fmt.Println("reporter.loop")
	for {
		fmt.Println("postgres.query")
		for _, qf := range sqlQueryfiles {
			rows, err := db.Query(qf.Sql)
			if err != nil {
				panic(err)
			}
			cols, err := rows.Columns()
			if err != nil {
				panic(err)
			}

			switch len(cols) {
			case 1:
				present := rows.Next()
				if !present {
					panic("No row")
				}
				var value float64
				err = rows.Scan(&value)
				if err != nil {
					panic(err)
				}
				metric := librato.Metric{
					Name:  qf.Name,
					Value: value,
				}
				fmt.Printf("postgres.result name=%s value=%f\n", qf.Name, value)
				fmt.Println("librato.post")
				lb.PostMetrics(&librato.Metrics{
					Gauges: []interface{}{metric},
				})
			case 2:
				metrics := []interface{}{}
				for rows.Next() {
					var source string
					var value float64
					err = rows.Scan(&source, &value)
					if err != nil {
						panic(err)
					}
					fmt.Printf("postgres.result name=%s source=%s value=%f\n", qf.Name, source, value)
					metric := librato.Metric{
						Name:   qf.Name,
						Source: source,
						Value:  value,
					}
					metrics = append(metrics, metric)
				}
				fmt.Println("librato.post")
				lb.PostMetrics(&librato.Metrics{
					Gauges: metrics,
				})
			default:
				panic("Must return 1 or 2 columns")
			}
		}

		fmt.Println("reporter.wait")
		time.Sleep(time.Duration(queryInterval) * time.Second)
	}
}
