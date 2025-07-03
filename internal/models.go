package internal

import "time"

// RSSRequest represents the incoming request structure
type RSSRequest struct {
	URLs []string `json:"urls"`
}

// RSSResponse represents the response structure
type RSSResponse struct {
	Items []RSSItem `json:"items"`
}

// RSSItem represents a single RSS item in the response
type RSSItem struct {
	Title       string     `json:"title"`
	Source      string     `json:"source"`
	SourceURL   string     `json:"source_url"`
	Link        string     `json:"link"`
	PublishDate *time.Time `json:"publish_date"`
	Description string     `json:"description"`
}
