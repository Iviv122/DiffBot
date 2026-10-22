package main

import (
	"log"
	"time"

	Requester "github.com/Iviv122/diffbot/internal/requester"
	DB "github.com/Iviv122/diffbot/internal/requester/db"
	_ "github.com/glebarez/go-sqlite"
)

const refresh_rate = time.Second           // every second
const next_refresh_rate = time.Second * 10 // every 5 seconds

func main() {
	if false {
		Requester.Request("https://rezerwacje.duw.pl/status_kolejek/query.php?status", time.Now().Add(next_refresh_rate))
	}
	now, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}
	for {

		db, err := DB.GetDb()
		if err != nil {
			log.Fatal(err)
			time.Sleep(refresh_rate)
			continue
		}

		lines, err := DB.GetList(db)
		if err != nil {
			log.Fatal(err)
			time.Sleep(refresh_rate)
			continue
		}

		for _, item := range lines {
			if item.NextTime.Before(now) {
				_, err := Requester.Request(item.Url, now.Add(next_refresh_rate))
				if err != nil {
					log.Println("Request failed:", err)
					continue
				}
				log.Default().Print("-> ", item.NextTime)
			}

		}
		now = now.Add(refresh_rate)
		time.Sleep(refresh_rate)
	}

}
