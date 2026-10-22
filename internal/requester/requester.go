package requester

import (
	"io"
	"net/http"
	"time"

	DB "github.com/Iviv122/diffbot/internal/requester/db"
)

func GetClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

func Request(link string, next_time int64) (string, error) {
	client := GetClient()
	db, err := DB.GetDb()
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "curl/8.0.0")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close() // close after main returns

	DB.Insert(db, link, 60, string(body))

	return string(body), err
}
