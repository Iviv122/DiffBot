package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

const TABLE_NAME = "pages"

var db *sql.DB
var get_lock = &sync.Mutex{}
var conn_lock = &sync.Mutex{}

type Link struct {
	Url         string
	NextTime    time.Time
	Body        string
	RefreshTime time.Duration
}

func GetDb() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	get_lock.Lock()
	defer get_lock.Unlock()
	_db, err := sql.Open("sqlite", "./foo.db")
	db = _db
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        url TEXT PRIMARY KEY,
        time TEXT,
        page TEXT,
        refresh_time INTEGER
    )`, TABLE_NAME))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Insert(url string, repeat_time time.Time, page string, repeat_wait time.Duration) error {
	if db == nil {
		GetDb()
	}
	conn_lock.Lock()
	defer conn_lock.Unlock()

	query := fmt.Sprintf(
		`INSERT OR REPLACE INTO %s (url, time, page, refresh_time) VALUES (?, ?, ?, ?)`, TABLE_NAME)

	_, err := db.Exec(
		query,
		url,
		repeat_time.Format("2006-01-02 15:04:05"),
		page,
		repeat_wait,
	)
	if err != nil {
		log.Print("Insert: ", err)
		return err
	}

	return nil
}
func UpdateTime(url string, repeat_time time.Time) error {
	if db == nil {
		GetDb()
	}
	conn_lock.Lock()
	defer conn_lock.Unlock()

	query := fmt.Sprintf("UPDATE %s SET time = ? WHERE url=?", TABLE_NAME)
	_, err := db.Exec(
		query,
		url,
		repeat_time.Format("2006-01-02 15:04:05"),
	)
	if err != nil {

		log.Print("UpdateTime: ", err)
		return err
	}
	return nil
}
func GetList() ([]Link, error) {
	if db == nil {
		GetDb()
	}
	conn_lock.Lock()
	defer conn_lock.Unlock()

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
		err := rows.Scan(&link.Url, &next_time, &link.Body, &link.RefreshTime)
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
