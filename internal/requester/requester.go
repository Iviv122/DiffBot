package requester

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	DB "github.com/Iviv122/diffbot/internal/requester/db"

	"github.com/pmezard/go-difflib/difflib"
)

type ServerState int

func getClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
	}
}

// return if same and diff text
func AreSame(text1 string, text2 string) (bool, string) {
	if strings.Compare(text1, text2) == 0 {
		return true, ""
	}
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(text1),
		B:        difflib.SplitLines(text2),
		FromFile: "Old",
		ToFile:   "New",
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	return false, text
}

func Process(link string, next_time time.Time, repeat_time time.Duration) (string, error) {

	body, err := request(link)

	if err != nil {
		err = DB.Insert(link, next_time, "Bad fetch", repeat_time)
	}
	err = DB.Insert(link, next_time, string(body), repeat_time)
	if err != nil {
		log.Print("Request_Com:", err)
		return "", err
	}

	return string(body), err
}
func request(link string) (string, error) {
	client := getClient()
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		panic(err)
	}
	// more answers, also it is supposed to be moral
	req.Header.Set("User-Agent", "curl/8.0.0")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {

		log.Print("Request_ERR_START:", req.URL, resp, err)
		return "", err
	}
	return string(body), nil
}
