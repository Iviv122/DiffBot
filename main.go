package main

import (
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	db, err := db_init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, "https://rezerwacje.duw.pl/status_kolejek/query.php?status", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "curl/8.0.0")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // close after main returns

	body, err := io.ReadAll(resp.Body)
	log.Printf("Status: %s", resp.Status)
	log.Printf("Final URL: %s", resp.Request.URL)
	log.Printf("Content-Type: %s", resp.Header.Get("Content-Type"))
	log.Printf("Body: %s", body)
	insert(db, "https://rezerwacje.duw.pl/status_kolejek/", time.Now().Format("2006-01-02 15:04:05"), string(body))
}
