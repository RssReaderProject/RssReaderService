# RssReaderService

A Go-based RSS feed parsing service that provides a REST API for parsing multiple RSS feeds simultaneously.

## Features

- Parse multiple RSS feeds in a single request
- RESTful API with JSON request/response
- Health check endpoint
- Docker support
- Multi-platform Docker images (AMD64, ARM64)

## API Endpoints

### POST /rss
Parse RSS feeds from provided URLs.

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
      "sourceURL": "https://source.com",
      "link": "https://source.com/article",
      "publishDate": "2024-01-01T12:00:00Z",
      "description": "Article description..."
    }
  ]
}
```

### GET /health
Health check endpoint that returns `200 OK`.

## Running the Service

### Using Docker (Recommended)

The service is available as a Docker image on GitHub Container Registry:

```bash
# Pull the latest image
docker pull ghcr.io/RssReaderProject/RssReaderService:latest

# Run the container
docker run -p 8080:8080 ghcr.io/RssReaderProject/RssReaderService:latest
```

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/RssReaderProject/RssReaderService.git
cd RssReaderService
```

2. Build the application:
```bash
go build -o server ./cmd/server
```

3. Run the server:
```bash
./server
```

The service will start on port 8080 by default. You can change the port by setting the `PORT` environment variable.

## Development

### Prerequisites
- Go 1.24 or later

### Running Tests
```bash
go test -v ./...
```

### Running Linter
```bash
golangci-lint run
```

### Building Docker Image Locally
```bash
docker build -t rss-reader-service .
docker run -p 8080:8080 rss-reader-service
```

## Docker Images

Docker images are automatically built and pushed to GitHub Container Registry on:
- Every push to the `main` branch
- Every tag push (e.g., `v1.0.0`)

Available tags:
- `latest` - Latest build from main branch
- `main` - Latest build from main branch
- `v*` - Semantic version tags
- `main-<sha>` - Build from specific commit

## CI/CD

The project uses GitHub Actions for continuous integration and deployment:

1. **CI Pipeline** (`ci.yml`):
   - Runs tests with Go 1.24
   - Performs linting with golangci-lint
   - Runs security checks with govulncheck

2. **Docker Build Pipeline** (`docker.yml`):
   - Triggers after successful CI completion
   - Builds multi-platform Docker images
   - Pushes to GitHub Container Registry
