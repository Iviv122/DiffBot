package main

import (
	"fmt"
	"log"
	"time"

	DB "github.com/Iviv122/diffbot/internal/requester/db"
	_ "github.com/glebarez/go-sqlite"
)

const refresh_rate = time.Second // every second

func main() {
	now := time.Now()
	for {
		db, err := DB.GetDb()
		if err != nil {
			log.Fatal(err)
			continue
		}
		lines, err := DB.GetList(db)
		for _, item := range lines {
			fmt.Println(item.Url)
		}
		now.Add(time.Second)
		time.Sleep(time.Second)
	}

}
