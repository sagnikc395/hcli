package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
)

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

func request(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code : %d\n", resp.StatusCode)
	}

	return resp, nil
}

// parse the RSS response
func parseRSS(resp *http.Response) (*RSS, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rss := &RSS{}
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.Strict = false
	if err := decoder.Decode(rss); err != nil {
		return nil, err
	}

	return rss, nil
}

func printPosts(items []Item) {
	for i, item := range items {
		fmt.Printf("%d. %s\n", i+1, item.Title)
	}
}

func userPromptForPost(items []Item) {
	var postnum int
	fmt.Printf("Type post number to open, or 0 to quit:")
	_, err := fmt.Scanf("%d", &postnum)
	if err != nil {
		fmt.Printf("Invalid input\n")
		return
	}
	if postnum == 0 {
		return
	}

	if postnum < 1 || postnum > len(items) {
		fmt.Printf("Invalid post number")
		return
	}
	openURL(items[postnum-1].Link)
	userPromptForPost(items)
}

func openURL(url string) {
	cmd := initHCLI() + " \"" + url + "\""
	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		fmt.Printf("Error opening URL: %s\n", err)
	}
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
	resp, err := request("https://news.ycombinator.com/rss")
	if err != nil {
		fmt.Println("An error occured: ", err)
		return
	}
	rss, err := parseRSS(resp)
	if err != nil {
		fmt.Println("Error parsing RSS: ", err)
		return
	}

	printPosts(rss.Items)
	userPromptForPost(rss.Items)
}
