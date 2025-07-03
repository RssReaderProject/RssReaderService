package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	rssreader "github.com/RssReaderProject/RssReader"
)

// RegisterRoutes attaches all internal handlers to the provided mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /rss", HandlePostRSSParse)
}

// HandlePostRSSParse handles RSS parsing requests
func HandlePostRSSParse(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req RSSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if len(req.URLs) == 0 {
		http.Error(w, "No URLs provided", http.StatusBadRequest)
		return
	}

	// Create context with timeout for RSS parsing
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// Parse RSS feeds using the RssReader package
	rssItems, err := rssreader.Parse(ctx, req.URLs)
	if err != nil {
		http.Error(w, "Error parsing RSS feeds: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert rssreader.RssItem to internal.RssServiceItem
	items := make([]RssServiceItem, len(rssItems))
	for i, item := range rssItems {
		items[i] = RssServiceItem{
			Title:       item.Title,
			Source:      item.Source,
			SourceURL:   item.SourceURL,
			Link:        item.Link,
			PublishDate: &item.PublishDate,
			Description: item.Description,
			RssURL:      item.RssURL,
		}
	}

	// Create response
	response := RSSServiceResponse{
		Items: items,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
