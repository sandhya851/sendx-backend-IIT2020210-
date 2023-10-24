package main

import (
    "html/template"
    "log"
    "net/http"
)

func main() {
    // Parse the HTML template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        log.Fatal("Error parsing template:", err)
    }

    // Handle static files (CSS, JS, images, etc.)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Handle the root URL with the parsed template
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Execute the template, passing nil as data
        tmpl.Execute(w, nil)
    })

    // Start the HTTP server
    http.ListenAndServe(":8080", nil)
}
