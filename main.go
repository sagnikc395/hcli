package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"runtime"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"title"`
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

func request(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code : %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// parse the RSS response
func parseRSS(doc *goquery.Document) (*RSS, error) {
	rss := &RSS{}
	err := xml.Unmarshal([]byte(doc.Text()), rss)
	if err != nil {
		return nil, err
	}
	return rss, nil
}



func initHCLI() string {
	switch runtime.GOOS {
	case "windows":
		return "start"
	case "linux":
		return "xdg-open"
	case "darwin":
		return "open"
	default:
		return "open"
	}
}

func main() {

}
