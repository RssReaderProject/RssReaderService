package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RssReaderProject/RssReaderService/internal"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	// Create a new request to the health endpoint
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new mux and register the health handler
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Error writing health response: %v", err)
		}
	})

	// Serve the request
	mux.ServeHTTP(rr, req)

	// Check the status code
	require.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body
	expected := "OK"
	require.Equal(t, expected, rr.Body.String(), "handler returned unexpected body")
}

func TestRSSEndpointExists(t *testing.T) {
	// Create a valid RSS request body
	requestBody := internal.RSSRequest{
		URLs: []string{"https://example.com/feed.xml"},
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	// Create a new request to the RSS endpoint
	req, err := http.NewRequest("POST", "/rss", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a new mux and register the RSS handler
	mux := http.NewServeMux()
	mux.HandleFunc("POST /rss", internal.HandlePostRSSParse)

	// Serve the request
	mux.ServeHTTP(rr, req)

	// Check that we don't get a 404 (Not Found)
	require.NotEqual(t, http.StatusNotFound, rr.Code, "RSS endpoint should not return 404")

	// Additional check: should return a successful status code (200 or 400 for bad request is acceptable)
	require.True(t, rr.Code == http.StatusOK || rr.Code == http.StatusBadRequest,
		"RSS endpoint should return either 200 (OK) or 400 (Bad Request), got %d", rr.Code)
}
