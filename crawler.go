package crawler

import (
	"errors"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Crawler represents a web crawler with caching mechanism.
type Crawler struct {
	cache     map[string]time.Time
	cacheLock sync.RWMutex
}

// CrawlResult represents the result of crawling a URL.
type CrawlResult struct {
	URL     string
	Content string
	Error   error
}

// NewCrawler creates a new instance of Crawler.
func NewCrawler() *Crawler {
	return &Crawler{
		cache: make(map[string]time.Time),
	}
}

// Crawl performs crawling for the given URLs using multiple workers.
func (c *Crawler) Crawl(urls []string, isPayingCustomer bool, numWorkers int) map[string]CrawlResult {
	var wg sync.WaitGroup
	urlChan := make(chan string, len(urls))
	resultChan := make(chan CrawlResult, len(urls))

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go c.crawlWorker(&wg, urlChan, resultChan, isPayingCustomer)
	}

	// Send URLs to be crawled
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from resultChan
	results := make(map[string]CrawlResult)
	for r := range resultChan {
		results[r.URL] = r
	}

	return results
}

// crawlWorker is a worker routine that crawls URLs.
func (c *Crawler) crawlWorker(wg *sync.WaitGroup, urlChan <-chan string, resultChan chan<- CrawlResult, isPayingCustomer bool) {
	defer wg.Done()

	for url := range urlChan {
		content, err := c.fetchAndParse(url)
		if err != nil {
			resultChan <- CrawlResult{URL: url, Error: err}
			continue
		}

		resultChan <- CrawlResult{URL: url, Content: content, Error: nil}
	}
}

// fetchAndParse fetches and parses HTML content from the given URL.
func (c *Crawler) fetchAndParse(url string) (string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}

	content := ""
	doc.Find("body").Each(func(index int, element *goquery.Selection) {
		content += element.Text() + "\n"
	})

	if content == "" {
		return "", errors.New("no content found on the page")
	}

	return content, nil
}
