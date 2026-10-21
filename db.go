package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

var db *sql.DB
var lock = &sync.Mutex{}

const TABLE_NAME = "pages"

func db_init() (*sql.DB, error) {
	_db, err := sql.Open("sqlite", "./foo.db")
	db = _db
	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        url TEXT PRIMARY KEY,
        time TEXT,
        page TEXT
    )`, TABLE_NAME))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func insert(db *sql.DB, url string, time string, page string) {
	query := fmt.Sprintf(`INSERT INTO %s VALUES ('%s', '%s', '%s')`, TABLE_NAME, url, time, page)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
