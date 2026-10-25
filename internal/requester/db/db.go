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
var insert_lock = &sync.Mutex{}

type Link struct {
	Url         string
	NextTime    time.Time
	Body        string
	RefreshTime time.Duration // seconds
}

func InitDb(url string) error {
	get_lock.Lock()
	defer get_lock.Unlock()
	_db, err := sql.Open("sqlite", url)
	db = _db
	if err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        url TEXT PRIMARY KEY,
        time TEXT,
        page TEXT,
        refresh_time INTEGER
    )`, TABLE_NAME))
	if err != nil {
		return err
	}
	return nil
}

func getDb() (*sql.DB, error) {
	return db, nil
}

// refresh_time in seconds
func Insert(url string, repeat_time time.Time, page string, refresh_time time.Duration) error {
	insert_lock.Lock()
	defer insert_lock.Unlock()

	query := fmt.Sprintf(
		`INSERT OR REPLACE INTO %s (url, time, page, refresh_time) VALUES (?, ?, ?, ?)`, TABLE_NAME)

	_, err := db.Exec(
		query,
		url,
		repeat_time.Format("2006-01-02 15:04:05"),
		page,
		refresh_time,
	)
	if err != nil {
		log.Print("Insert: ", err)
		return err
	}

	return nil
}
func Find() {

}
func UpdateTime(url string, repeat_time time.Time) error {
	insert_lock.Lock()
	defer insert_lock.Unlock()

	query := fmt.Sprintf("UPDATE %s SET time = ? WHERE url=?", TABLE_NAME)
	_, err := db.Exec(
		query,
		repeat_time.Format("2006-01-02 15:04:05"),
		url,
	)
	if err != nil {

		log.Print("UpdateTime: ", err)
		return err
	}
	return nil
}
func GetList() ([]Link, error) {

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
