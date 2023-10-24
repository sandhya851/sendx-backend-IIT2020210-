 package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yourusername/web-crawler/storage"
	"golang.org/x/time/rate"
)

type API struct {
	limiter *rate.Limiter
}

func NewAPI() *API {
	// Allow 100 requests per second with burst of 1 request.
	return &API{
		limiter: rate.NewLimiter(rate.Limit(100), 1),
	}
}

func (a *API) CrawlHandler(w http.ResponseWriter, r *http.Request) {
	// Apply rate limiting
	if err := a.limiter.Wait(r.Context()); err != nil {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// For POST requests, retrieve the URL from the form data
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		url := r.Form.Get("url")
		// URL validation and sanitization can be added here if needed.
		crawlURL(w, r, url)
		return
	}

	// For GET requests, retrieve the URL from the query parameter
	url := r.URL.Query().Get("url")
	crawlURL(w, r, url)
}

func crawlURL(w http.ResponseWriter, r *http.Request, url string) {
	// Convert URL to a valid file name by replacing non-alphanumeric characters with underscores
	fileName := generateFileName(url)

	content, err := storage.GetCrawledData(fileName)
	if err == nil {
		// Data found in storage, return it
		fmt.Fprint(w, content)
		return
	}

	// Data not found, crawl in real-time
	crawledContent, err := fetchURLContent(url, r.Context())
	if err != nil {
		log.Println("Error fetching URL:", err)
		http.Error(w, "Failed to crawl the URL", http.StatusInternalServerError)
		return
	}

	// Save crawled data to storage
	err = storage.SaveCrawledData(fileName, crawledContent)
	if err != nil {
		log.Println("Error saving crawled data:", err)
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

func fetchURLContent(url string, ctx context.Context) (string, error) {
	// Implement logic to crawl the URL and fetch content with context handling
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
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

func TestCrawlHandler(t *testing.T) {
	// Create a form values map with the URL input
	form := url.Values{}
	form.Add("url", "https://example.com")

	// Encode the form values into the request body
	reqBody := form.Encode()

	// Create a new request to the /crawl endpoint with the form values in the body
	req, err := http.NewRequest("POST", "/crawl", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// Create an HTTP handler from the CrawlHandler function
	handler := http.HandlerFunc((&API{}).CrawlHandler)

	// Call the handler's ServeHTTP method to handle the request
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body (assuming your handler returns a specific response for success)
	expectedResponse := "Crawling successful"
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}
