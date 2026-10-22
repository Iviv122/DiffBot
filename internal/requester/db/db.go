package db

import (
	"database/sql"
	"fmt"
	"log"
)

const TABLE_NAME = "pages"

func GetDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./foo.db")
	if err != nil {
		return nil, err
	}
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

// TODO: Sanitize input
func Insert(db *sql.DB, url string, time string, page string) {
	query := fmt.Sprintf(`INSERT OR REPLACE INTO %s VALUES ('%s', '%s', '%s')`, TABLE_NAME, url, time, page)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
