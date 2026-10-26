package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	Requester "github.com/Iviv122/diffbot/internal/requester"
	DB "github.com/Iviv122/diffbot/internal/requester/db"
	_ "github.com/glebarez/go-sqlite"
	"github.com/nuskey8/go-cliapp"
)

const refresh_rate = time.Second

func main() {
	log.SetOutput(os.Stdout)
	app := cliapp.Default()

	app.Add("run", "<db_path> starts process", func(db_url string) {
		DB.InitDb(db_url)
		run()
	})
	app.Add("populate", "<list_path> <db_path> Insert entries from list", func(list string, db_url string) {
		DB.InitDb(db_url)
		populate(list)
	})
	app.Add("fetch", "<db_path> manually fetch everything", func(db_url string) {
		DB.InitDb(db_url)
		fetch()
	})
	app.Add("initdb", "<db_path> create file and/or table", func(db_url string) {
		DB.InitDb(db_url)
	})

	err := app.Run()
	if err != nil {
		log.Print(err)
	}
}

func run() {
	now, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		panic(err)
	}
	for {

		lines, err := DB.GetList()
		if err != nil {
			log.Fatal(err)
			time.Sleep(refresh_rate)
			continue
		}

		for _, item := range lines {
			go func() {
				if item.NextTime.Before(now) {
					_, err := Requester.Process(item.Url, now.Add(item.RefreshTime), item.RefreshTime)
					if err != nil {
						log.Println("Request failed:", err)
					}
					log.Default().Print(item.Url, " -> ", item.NextTime)
				}
			}()

		}
		now = now.Add(refresh_rate)
		time.Sleep(refresh_rate)
	}
}

// time goes in second
func populate(filename string) {

	var wg sync.WaitGroup
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Populate: ", err)
	}
	defer f.Close()

	var i = 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		a := strings.Split(scanner.Text(), " ")
		wg.Go(func() {
			log.Print("Populate: ", a[0])
			if len(a) < 2 {
				log.Print("Error in line: ", i)
				return
			}
			t, _ := strconv.Atoi(a[1])

			wait_seconds := time.Second * time.Duration(t)

			_, err := Requester.Process(a[0], time.Now().Add(wait_seconds), wait_seconds)

			if err != nil {
				log.Print("Populate_Request: ", err)
				log.Default().Print()
			}
		})
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Populate_end: ", err)
		log.Fatal(err)
	}
	wg.Wait()
}

// fetch without update
func fetch() {
	var wg sync.WaitGroup
	lines, err := DB.GetList()
	if err != nil {
		log.Print(err)
		panic(err)
	}

	for _, item := range lines {
		wg.Go(func() {
			_, err := Requester.Process(item.Url, item.NextTime, item.RefreshTime)
			if err != nil {
				log.Println("Request failed:", err)
			}
			log.Print(item.Url, " -> ", item.NextTime)
		})
	}
	wg.Wait()
}
