package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// setupTestServer creates a test server with registered routes
func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

func TestHandlePostRSSParse_NoURLs(t *testing.T) {
	// Setup test server
	server := setupTestServer()
	defer server.Close()

	// Create a request with no URLs
	reqBody := RSSRequest{
		URLs: []string{},
	}

	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Create HTTP request to the test server
	req, err := http.NewRequest("POST", server.URL+"/rss", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	}()

	// Assert that we get a 400 Bad Request status
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected 400 Bad Request status")

	// Read response body as plain text
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// Assert that the response contains the expected error message
	require.Contains(t, bodyString, "No URLs provided", "Expected error message about no URLs")
}

func TestHandlePostRSSParse_InvalidRequestBody(t *testing.T) {
	// Setup test server
	server := setupTestServer()
	defer server.Close()

	// Create a request with invalid JSON
	invalidJSON := []byte(`{"urls": "not an array"}`)

	// Create HTTP request to the test server
	req, err := http.NewRequest("POST", server.URL+"/rss", bytes.NewBuffer(invalidJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	}()

	// Assert that we get a 400 Bad Request status
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected 400 Bad Request status")

	// Read response body as plain text
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)

	// Assert that the response contains the expected error message
	require.Contains(t, bodyString, "Invalid request body", "Expected error message about invalid request body")
}
