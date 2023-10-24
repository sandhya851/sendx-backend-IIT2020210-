package server

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	if url == "" {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL"})
		return
	}

	payingCustomerParam := r.URL.Query().Get("paying_customer")
	isPayingCustomer, err := strconv.ParseBool(payingCustomerParam)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid paying_customer parameter"})
		return
	}

	results := s.Crawler.Crawl([]string{url}, isPayingCustomer, 1)
	jsonResponse(w, http.StatusOK, results[url])
}

func (s *Server) SetWorkerCountHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		NumWorkers int `json:"num_workers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if request.NumWorkers <= 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid number of workers"})
		return
	}

	s.Crawler.SetNumWorkers(request.NumWorkers)
	jsonResponse(w, http.StatusOK, map[string]string{"message": "Number of workers updated successfully"})
}

func (s *Server) SetCrawlSpeedHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CrawlSpeed int `json:"crawl_speed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if request.CrawlSpeed <= 0 {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid crawl speed"})
		return
	}

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
