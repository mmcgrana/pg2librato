package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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
	fmt.Println("start")
	interval := QueryInterval()
	dbUrl := DatabaseUrl()

	fmt.Println("read")
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
	fmt.Println("%+v\n", sqlQueryfiles)

	fmt.Println("connect")
	db, err := sql.Open("postgres", dbUrl)
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

	fmt.Println("loop")
	for {
		fmt.Println("query")
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
				fmt.Println("1 col")
			case 2:
				fmt.Println("2 cols")
			default:
				panic("Must return 1 or 2 columns")
			}
		}

		fmt.Println("report")

		fmt.Println("wait")
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
