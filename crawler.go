package crawler

import (
	"errors"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	cache     map[string]time.Time
	cacheLock sync.RWMutex
}

type CrawlResult struct {
	URL     string
	Content string
	Error   error
}

func NewCrawler() *Crawler {
	return &Crawler{
		cache: make(map[string]time.Time),
	}
}

func (c *Crawler) Crawl(urls []string, isPayingCustomer bool, numWorkers int) map[string]CrawlResult {
	var wg sync.WaitGroup
	urlChan := make(chan string, len(urls))
	resultChan := make(chan CrawlResult, len(urls))

	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go c.crawlWorker(&wg, urlChan, resultChan, isPayingCustomer)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make(map[string]CrawlResult)
	for r := range resultChan {
		results[r.URL] = r
	}

	return results
}

func (c *Crawler) crawlWorker(wg *sync.WaitGroup, urlChan <-chan string, resultChan chan<- CrawlResult, isPayingCustomer bool) {
	defer wg.Done()

	for url := range urlChan {
		// Implement logic to crawl the URL here
		content, err := c.fetchAndParse(url)
		if err != nil {
			resultChan <- CrawlResult{URL: url, Error: err}
			continue
		}

		// Store content in resultChan
		resultChan <- CrawlResult{URL: url, Content: content, Error: nil}
	}
}

func (c *Crawler) fetchAndParse(url string) (string, error) {
	// Implement logic to fetch and parse HTML content using goquery or other libraries
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", err
	}

	// Extract content from the parsed HTML document
	content := ""
	doc.Find("body").Each(func(index int, element *goquery.Selection) {
		content += element.Text() + "\n"
	})

	if content == "" {
		return "", errors.New("no content found on the page")
	}

	return content, nil
}
