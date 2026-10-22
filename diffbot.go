package main

import (
	"io"

	Requester "github.com/Iviv122/diffbot/internal/requester"
	_ "github.com/glebarez/go-sqlite"
)

func main() {

	resp, err := Requester.Request("https://rezerwacje.duw.pl/status_kolejek/query.php?status")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	print(body)
}
