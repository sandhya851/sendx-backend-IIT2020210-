// main.go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/PuerkitoBio/goquery"
)

type CrawlCache struct {
	LastCrawled time.Time
	PageData   string
}

var cache map[string]CrawlCache

func main() {
	cache = make(map[string]CrawlCache)

	r := mux.NewRouter()
	r.HandleFunc("/crawl", crawlHandler).Methods("GET")
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}

func crawlHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	isPayingCustomer := r.URL.Query().Get("paying_customer")

	if cacheData, ok := cache[url]; ok && time.Since(cacheData.LastCrawled).Minutes() < 60 {
		// Page exists in cache, return cached page
		fmt.Fprintf(w, "Cached Page Data: %s", cacheData.PageData)
		return
	}

	// Page not in cache, crawl in real-time
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Fprintf(w, "Error fetching URL: %s", err.Error())
		return
	}

	pageData := doc.Text()
	cache[url] = CrawlCache{
		LastCrawled: time.Now(),
		PageData:   pageData,
	}

	fmt.Fprintf(w, "Real-time Crawled Page Data: %s", pageData)
}
