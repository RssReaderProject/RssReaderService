package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlePostRSSParse_NoURLs(t *testing.T) {
	// Create a request with no URLs
	reqBody := RSSRequest{
		URLs: []string{},
	}

	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/rss/parse", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	HandlePostRSSParse(rr, req)

	// Assert that we get a 400 Bad Request status
	require.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request status")

	// Assert that the response body contains the expected error message
	require.Contains(t, rr.Body.String(), "No URLs provided", "Expected error message about no URLs")
}

func TestHandlePostRSSParse_InvalidRequestBody(t *testing.T) {
	// Create a request with invalid JSON
	invalidJSON := []byte(`{"urls": "not an array"}`)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/rss/parse", bytes.NewBuffer(invalidJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	HandlePostRSSParse(rr, req)

	// Assert that we get a 400 Bad Request status
	require.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request status")

	// Assert that the response body contains the expected error message
	require.Contains(t, rr.Body.String(), "Invalid request body", "Expected error message about invalid request body")
}
