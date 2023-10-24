package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/web-crawler/crawler"
	"github.com/yourusername/web-crawler/server"
	"github.com/yourusername/web-crawler/api"
)

func main() {
	// Initialize flags for command line options
	port := flag.Int("port", 8080, "HTTP server port")
	numWorkers := flag.Int("workers", 5, "Number of crawler workers")
	flag.Parse()

	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// Initialize logger
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize Crawler, Server, and API instances
	c := crawler.NewCrawler(*numWorkers)
	s := server.NewServer(c)
	a := api.NewAPI()

	// Handle static files (CSS, JS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle the root URL with the parsed template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Execute the template, passing nil as data
		tmpl.Execute(w, nil)
	})

	// Set up API routes
	http.HandleFunc("/crawl", a.CrawlHandler)
	// Add other API endpoints if needed
	http.HandleFunc("/set-worker-count", s.SetWorkerCountHandler)
	http.HandleFunc("/set-crawl-speed", s.SetCrawlSpeedHandler)

	// Start the HTTP server
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
		if err != nil {
			logger.Fatalf("Error starting server: %v", err)
		}
	}()

	logger.Printf("Server started on port %d\n", *port)

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Println("Shutting down server...")
}
