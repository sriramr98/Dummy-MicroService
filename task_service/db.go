package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func InitDB(host string, user string, password string, dbName string) *sql.DB {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return db
}
