package server

import (
	"encoding/json" 
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/web-crawler/crawler"
)

type Server struct {
	Crawler *crawler.Crawler
}

func NewServer(c *crawler.Crawler) *Server {
	return &Server{
		Crawler: c,
	}
}

func (s *Server) CrawlHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	payingCustomerParam := r.URL.Query().Get("paying_customer")
	isPayingCustomer, _ := strconv.ParseBool(payingCustomerParam)

	results := s.Crawler.Crawl([]string{url}, isPayingCustomer, 1)
	jsonResponse(w, http.StatusOK, results[url])
}


func (s *Server) SetWorkerCountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the desired number of workers
	var request struct {
		NumWorkers int `json:"num_workers"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Set the number of crawler workers
	s.Crawler.SetNumWorkers(request.NumWorkers)

	jsonResponse(w, http.StatusOK, map[string]string{"message": "Number of workers updated successfully"})
}

func (s *Server) SetCrawlSpeedHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the desired crawl speed (pages per hour per worker)
	var request struct {
		CrawlSpeed int `json:"crawl_speed"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Set the crawl speed per hour per worker
	s.Crawler.SetCrawlSpeed(request.CrawlSpeed)

	jsonResponse(w, http.StatusOK, map[string]string{"message": "Crawl speed updated successfully"})
}



func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/crawl", s.CrawlHandler).Methods("GET")
	r.HandleFunc("/set-worker-count", s.SetWorkerCountHandler).Methods("POST")
	r.HandleFunc("/set-crawl-speed", s.SetCrawlSpeedHandler).Methods("POST")
}
