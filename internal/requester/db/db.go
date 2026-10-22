package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const TABLE_NAME = "pages"

type Link struct {
	Url      string
	NextTime time.Time
	Body     string
}

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
func Insert(db *sql.DB, url string, repeat_time int64, page string) error {
	time := time.Now().Add(time.Duration(repeat_time) * time.Second)
	query := fmt.Sprintf(`INSERT OR REPLACE INTO %s VALUES ('%s', '%s', '%s')`, TABLE_NAME, url, time, page)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func GetList(db *sql.DB) ([]Link, error) {
	rows, err := db.Query(fmt.Sprintf(`SElECT * FROM %s`, TABLE_NAME))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var links = make([]Link, 0)
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		var link Link
		var next_time string
		err := rows.Scan(&link.Url, &next_time, &link.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		link.NextTime, err = time.Parse("2006-01-02 15:04:05", next_time)
		if err != nil {
			log.Fatal(err)
			continue
		}
		links = append(links, link)

	}

	return links, err
}
