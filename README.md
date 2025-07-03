# RSS Reader Service

A high-performance HTTP service for parsing and aggregating RSS feeds. Built with Go, this service provides a RESTful API to fetch and parse multiple RSS feeds simultaneously, returning structured data in JSON format.

## Features

- **Multi-feed parsing**: Parse multiple RSS feeds in a single request
- **Timeout handling**: Built-in timeout protection (60 seconds) for RSS parsing operations
- **JSON API**: Clean RESTful API with JSON request/response format
- **Health monitoring**: Built-in health check endpoint
- **Error handling**: Comprehensive error handling with appropriate HTTP status codes
- **Go 1.24.4**: Built with the latest Go version for optimal performance

## Installation

### Prerequisites

- Go 1.24.4 or later
- Git

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/RssReaderProject/RssReaderService.git
cd RssReaderService
```

2. Build the service:
```bash
go build -o server cmd/server/main.go
```

3. Run the service:
```bash
./server
```

The service will start on port 8080 by default. You can customize the port by setting the `PORT` environment variable:

```bash
PORT=3000 ./server
```

## API Documentation

### Base URL
```
http://localhost:8080
```

### Endpoints

#### Health Check
```
GET /health
```
Returns a simple "OK" response to verify the service is running.

**Response:**
```
200 OK
OK
```

#### Parse RSS Feeds
```
POST /rss
```
Parses multiple RSS feeds and returns aggregated content.

**Request Body:**
```json
{
  "urls": [
    "https://example.com/feed.xml",
    "https://another-site.com/rss"
  ]
}
```

**Response:**
```json
{
  "items": [
    {
      "title": "Article Title",
      "source": "Source Name",
      "source_url": "https://source.com",
      "link": "https://source.com/article",
      "publish_date": "2024-01-15T10:30:00Z",
      "description": "Article description or content..."
    }
  ]
}
```

**Error Responses:**

- `400 Bad Request`: Invalid request body or no URLs provided
- `500 Internal Server Error`: RSS parsing error or response encoding error

## Usage Examples

### Using curl

```bash
# Parse a single RSS feed
curl -X POST http://localhost:8080/rss \
  -H "Content-Type: application/json" \
  -d '{"urls": ["https://rss.cnn.com/rss/edition.rss"]}'

# Parse multiple RSS feeds
curl -X POST http://localhost:8080/rss \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://rss.cnn.com/rss/edition.rss",
      "https://feeds.bbci.co.uk/news/rss.xml",
      "https://www.reddit.com/r/golang/.rss"
    ]
  }'
```

### Using JavaScript/Fetch

```javascript
const response = await fetch('http://localhost:8080/rss', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    urls: [
      'https://rss.cnn.com/rss/edition.rss',
      'https://feeds.bbci.co.uk/news/rss.xml'
    ]
  })
});

const data = await response.json();
console.log(data.items);
```

## Development

### Project Structure

```
RssReaderService/
├── cmd/
│   └── server/
│       ├── main.go          # Application entry point
│       └── main_test.go     # Server tests
├── internal/
│   ├── handler.go           # HTTP handlers
│   ├── handler_test.go      # Handler tests
│   └── models.go            # Data structures
├── go.mod                   # Go module file
├── go.sum                   # Go module checksums
└── README.md               # This file
```

### Running Tests

```bash
go test ./...
```

### Dependencies

The service uses the following main dependencies:

- `github.com/RssReaderProject/RssReader`: Core RSS parsing functionality
- `github.com/stretchr/testify`: Testing utilities

## Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)

### Timeouts

- RSS parsing timeout: 60 seconds (configurable in `internal/handler.go`)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions, please open an issue on the GitHub repository.