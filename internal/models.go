package internal

import "time"

// RSSRequest represents the incoming request structure
type RSSRequest struct {
	URLs []string `json:"urls"`
}

// RSSServiceResponse represents the response structure
type RSSServiceResponse struct {
	Items []RssServiceItem `json:"items"`
}

// RssServiceItem represents a single RSS item in the response
type RssServiceItem struct {
	Title       string     `json:"title"`
	Source      string     `json:"source"`
	SourceURL   string     `json:"source_url"`
	Link        string     `json:"link"`
	PublishDate *time.Time `json:"publish_date,omitempty,omitzero"`
	Description string     `json:"description"`
}
