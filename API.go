package api

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/yourusername/web-crawler/storage"
)

type API struct {
	// You can add any necessary fields here
}

func NewAPI() *API {
	return &API{}
}

func (a *API) CrawlHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	content, err := storage.GetCrawledData(generateFileName(url))
	if err == nil {
		// Data found in storage, return it
		fmt.Fprint(w, content)
		return
	}

	// Data not found, crawl in real-time
	crawledContent, err := fetchURLContent(url)
	if err != nil {
		http.Error(w, "Failed to crawl the URL", http.StatusInternalServerError)
		return
	}

	// Save crawled data to storage
	err = storage.SaveCrawledData(generateFileName(url), crawledContent)
	if err != nil {
		http.Error(w, "Failed to save crawled data", http.StatusInternalServerError)
		return
	}

	// Return crawled content to the client
	fmt.Fprint(w, crawledContent)
}

func generateFileName(url string) string {
	// Convert URL to a valid file name by replacing non-alphanumeric characters with underscores
	return strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, url) + ".txt"
}

func fetchURLContent(url string) (string, error) {
	// Implement logic to crawl the URL and fetch content
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch URL: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract content from the parsed HTML document
	var content strings.Builder
	doc.Find("body").Each(func(index int, element *goquery.Selection) {
		content.WriteString(element.Text() + "\n")
	})

	return content.String(), nil
}
