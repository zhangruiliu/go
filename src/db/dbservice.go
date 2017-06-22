package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Type 数据库类型
type Type uint8

const (
	// DBOracle oracle
	DBOracle Type = iota + 1
	// DBMysql mysql
	DBMysql
	DBsqlite3
)

// DB 表示数据库
type DB struct {
	*sql.DB
	Type Type
}

// NewDB 返回数据库
func NewDB(driver string, connection string) (retDB *DB, err error) {
	tp := Type(0)
	if strings.Contains(driver, "oci") {
		tp = DBOracle
	} else if strings.Contains(driver, "mysql") {
		tp = DBMysql
	} else if strings.Contains(driver, "sqlite3") {
		tp = DBsqlite3
	} else {
		return nil, errors.New("not support databse with driver " + driver)
	}
	var db *sql.DB
	db, err = sql.Open(driver, connection)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		for {
			time.Sleep(time.Hour)
			db.SetMaxOpenConns(100)

		}
	}()

	switch tp {
	case DBOracle:
		go func() {
			for {
				time.Sleep(time.Minute)
				if rows, e := db.Query("select 1 from dual"); e != nil {
					fmt.Println(e)
					db, err = sql.Open(driver, connection)
					if err != nil {
						fmt.Println(e)
					}
				} else {
					rows.Close()
				}

			}
		}()
	case DBMysql:
		go func() {
			for {
				time.Sleep(time.Minute)
				if rows, e := db.Query("select 1"); e != nil {
					fmt.Println(e)
					db, err = sql.Open(driver, connection)
					if err != nil {
						fmt.Println(e)
					}
				} else {
					rows.Close()
				}

			}
		}()
	case DBsqlite3:
		go func() {
			for {
				time.Sleep(time.Minute)
				if e := db.Ping(); e != nil {
					fmt.Println(e)
					db, err = sql.Open(driver, connection)
					if err != nil {
						fmt.Println(e)
					}
				} else {
					//no action
				}

			}
		}()
	}
	retDB = &DB{
		db,
		tp,
	}
	return
}
