package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	Requester "github.com/Iviv122/diffbot/internal/requester"
	DB "github.com/Iviv122/diffbot/internal/requester/db"
	_ "github.com/glebarez/go-sqlite"
)

const refresh_rate = time.Second          // every second
const next_refresh_rate = time.Second * 5 // every 5 seconds

func main() {
	DB.GetDb()
	if true {
		go Populate("list1.sample")
	}
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
					_, err := Requester.Request(item.Url, now.Add(item.RefreshTime), item.RefreshTime)
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

func Populate(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Populate: ", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		a := strings.Split(scanner.Text(), " ")
		go func() {
			t, _ := strconv.Atoi(a[1])
			_, err := Requester.Request(a[0], time.Now().Add(time.Second*time.Duration(t)), time.Second*time.Duration(t))
			if err != nil {
				log.Print("Populate_Request: ", err)
			}
		}()
	}
	log.Println("ENDED POPULATION")
	if err := scanner.Err(); err != nil {
		log.Fatal("Populate_end: ", err)
		log.Fatal(err)
	}
}
