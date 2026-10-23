package requester

import (
	"io"
	"log"
	"net/http"
	"time"

	DB "github.com/Iviv122/diffbot/internal/requester/db"
)

func getClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

func Request(link string, next_time time.Time) (string, error) {
	client := getClient()
	db, err := DB.GetDb()
	if err != nil {
		log.Print("Request_DB:", err)
		return "", err
	}
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "curl/8.0.0")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {

		err = DB.Insert(db, link, next_time, "Bad fetch")
		log.Print("Request_ERR:", req.URL, resp, err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close() // close after main returns

	err = DB.Insert(db, link, next_time, string(body))
	if err != nil {
		log.Print("Request_Com:", err)
		return "", err
	}

	return string(body), err
}
