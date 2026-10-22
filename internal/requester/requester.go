package requester

import (
	"io"
	"net/http"
	"sync"
	"time"

	DB "github.com/Iviv122/diffbot/internal/requester/db"
)

var lock = &sync.Mutex{}
var client *http.Client

func GetClient() *http.Client {
	if client == nil {
		client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}
	return client

}

func Request(link string) (*http.Response, error) {
	client = GetClient()

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

	//log.Printf("Status: %s", resp.Status)
	//log.Printf("Final URL: %s", resp.Request.URL)
	//log.Printf("Content-Type: %s", resp.Header.Get("Content-Type"))
	//log.Printf("Body: %s", body)

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close() // close after main returns
	DB.Insert(db, link, time.Now().Format("2006-01-02 15:04:05"), string(body))

	return resp, err
}
