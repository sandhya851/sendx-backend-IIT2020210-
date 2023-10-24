// main.go
package main

import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"golang.org/x/net/html"

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


func main() {
	htmlContent := `<html><body><h1>Hello, World!</h1></body></html>`

	// Parse the HTML content
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}
	// Traverse the HTML parse tree
	traverseHTML(doc)
}
	func traverseHTML(node *html.Node) {
    if node.Type == html.ElementNode {
        // Do something with the element node, e.g., extract attributes or text content
        fmt.Println("Element:", node.Data)
    }

    // Recursive traversal for child nodes
    for c := node.FirstChild; c != nil; c = c.NextSibling {
        traverseHTML(c)
    }
}


}
