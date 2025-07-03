package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
