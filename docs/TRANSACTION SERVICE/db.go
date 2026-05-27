package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    connStr := "host=localhost port=5432 user=postgres password=postgres dbname=transactions_db sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    DB = db
    log.Println("Database connected")
}