package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
)

const TABLE_NAME = "pages"

func db_init(db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
        url TEXT PRIMARY KEY,
        time
    )`, TABLE_NAME))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func insert(db *sql.DB, url string) {
	_, err := db.Exec(fmt.Sprintf(`INSERT INTO %s VALUES ('%s')`, TABLE_NAME, url))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {

	db, err := sql.Open("sqlite", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db_init(db)

	resp, err := http.Get("https://rezerwacje.duw.pl/status_kolejek/")
	insert(db, "https://rezerwacje.duw.pl/status_kolejek/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // close after main returns
	body, err := io.ReadAll(resp.Body)
	print(string(body))
}
