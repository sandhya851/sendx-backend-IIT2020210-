package main

import (
    "html/template"
    "log"
    "net/http"
    "github.com/yourusername/web-crawler/crawler"
    "github.com/yourusername/web-crawler/server"
    "github.com/yourusername/web-crawler/api"
)

func main() {
    // Parse the HTML template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        log.Fatal("Error parsing template:", err)
    }

    // Initialize Crawler, Server, and API instances
    c := crawler.NewCrawler()
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
    http.ListenAndServe(":8080", nil)
}
