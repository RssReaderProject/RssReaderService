package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RssReaderProject/RssReaderService/internal"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create router and register handlers
	mux := http.NewServeMux()

	// Register RSS parsing endpoint
	mux.HandleFunc("POST /rss", internal.HandlePostRSSParse)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Error writing health response: %v", err)
		}
	})

	// Start server
	log.Printf("Starting RSS Reader Service on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
