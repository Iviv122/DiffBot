package requester

import (
	"io"
	"log"
	"net/http"
	"time"

	DB "github.com/Iviv122/diffbot/internal/requester/db"
)

type ServerState int

func getClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

func Request(link string, next_time time.Time, repeat_time time.Duration) (string, error) {
	client := getClient()
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "curl/8.0.0")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {

		log.Print("Request_ERR_START:", req.URL, resp, err)
		err = DB.Insert(link, next_time, "Bad fetch", repeat_time)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close() // close after main returns

	err = DB.Insert(link, next_time, string(body), repeat_time)
	if err != nil {
		log.Print("Request_Com:", err)
		return "", err
	}

	return string(body), err
}
