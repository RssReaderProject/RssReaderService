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

// setupMockRSSServer creates a mock RSS server that serves a basic RSS XML
func setupMockRSSServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		rssXML := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test RSS Feed</title>
    <link>http://example.com</link>
    <description>A test RSS feed for testing purposes</description>
    <item>
      <title>Test Article 1</title>
      <link>http://example.com/article1</link>
      <description>This is the first test article description</description>
      <pubDate>Mon, 01 Jan 2024 12:00:00 GMT</pubDate>
    </item>
    <item>
      <title>Test Article 2</title>
      <link>http://example.com/article2</link>
      <description>This is the second test article description</description>
      <pubDate>Tue, 02 Jan 2024 12:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`
		if _, err := w.Write([]byte(rssXML)); err != nil {
			// In a test context, we can log the error but can't return it
			// since this is a handler function
			panic(err)
		}
	})
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

func TestHandlePostRSSParse_ValidRSSFeed(t *testing.T) {
	// Setup mock RSS server
	mockRSSServer := setupMockRSSServer()
	defer mockRSSServer.Close()

	// Setup test server
	server := setupTestServer()
	defer server.Close()

	// Create a request with the mock RSS server URL
	reqBody := RSSRequest{
		URLs: []string{mockRSSServer.URL},
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

	// Assert that we get a 200 OK status
	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK status")

	// Read and parse response body
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response RSSServiceResponse
	err = json.Unmarshal(bodyBytes, &response)
	require.NoError(t, err)

	// Assert that we got exactly 2 items back (as defined in our mock RSS)
	require.Equal(t, 2, len(response.Items), "Expected exactly 2 RSS items")

	// Assert that the first item has the expected structure and content
	if len(response.Items) > 0 {
		item := response.Items[0]
		require.Equal(t, "Test Article 1", item.Title, "Expected correct title")
		require.Equal(t, "Test RSS Feed", item.Source, "Expected correct source")
		require.Equal(t, mockRSSServer.URL, item.SourceURL, "Expected correct source URL")
		require.Equal(t, "http://example.com/article1", item.Link, "Expected correct link")
		require.Equal(t, "This is the first test article description", item.Description, "Expected correct description")
		require.NotNil(t, item.PublishDate, "Expected non-nil publish date")
	}

	// Assert that the second item has the expected structure and content
	if len(response.Items) > 1 {
		item := response.Items[1]
		require.Equal(t, "Test Article 2", item.Title, "Expected correct title")
		require.Equal(t, "Test RSS Feed", item.Source, "Expected correct source")
		require.Equal(t, mockRSSServer.URL, item.SourceURL, "Expected correct source URL")
		require.Equal(t, "http://example.com/article2", item.Link, "Expected correct link")
		require.Equal(t, "This is the second test article description", item.Description, "Expected correct description")
		require.NotNil(t, item.PublishDate, "Expected non-nil publish date")
	}
}
