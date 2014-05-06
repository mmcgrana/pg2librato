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

	fmt.Println("postgres.connect.start")
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
	fmt.Println("postgres.connect.finish")

	fmt.Println("librato.connect.start")
	lb := &librato.Client{libratoAuth[0], libratoAuth[1]}
	fmt.Println("librato.connect.finish")

	fmt.Println("reporter.start")
	for {
		fmt.Println("reporter.loop.start")
		for _, qf := range sqlQueryfiles {
			fmt.Printf("postgres.query.start name=%s\n", qf.Name)
			rows, err := db.Query(qf.Sql)
			if err != nil {
				panic(err)
			}
			cols, err := rows.Columns()
			if err != nil {
				panic(err)
			}
			numCols := len(cols)
			if numCols != 3 {
				panic("Must return result set with exactly 3 rows")
			}
			metrics := []interface{}{}
			for rows.Next() {
				var name string
				var nullSource sql.NullString
				var source string
				var value float64
				err = rows.Scan(&name, &nullSource, &value)
				if err != nil {
					panic(err)
				}
				if nullSource.Valid {
					source = nullSource.String
				}
				fmt.Printf("postgres.result name=%s source=%s value=%f\n", name, source, value)
				metric := librato.Metric{
					Name:   name,
					Source: source,
					Value:  value,
				}
				metrics = append(metrics, metric)
			}
			fmt.Printf("postgres.query.finish name=%s\n", qf.Name)

			fmt.Printf("librato.post.start name=%s\n", qf.Name)
			err = lb.PostMetrics(&librato.Metrics{
				Gauges: metrics,
			})
			if err != nil {
				panic(err)
			}
			fmt.Printf("librato.post.finish name=%s\n", qf.Name)
		}

		fmt.Printf("reporter.loop.wait interval=%d\n", queryInterval)
		time.Sleep(time.Duration(queryInterval) * time.Second)
	}
}
