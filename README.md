# Image Analyzer API

Backend API for analyzing images using Google Cloud Vision AI.

## Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose (optional)
- Google Cloud Vision API credentials (required)

## Quick Start

### 1. Google Cloud Vision Setup

1. **Create a Google Cloud Project:**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select an existing one

2. **Enable the Vision API:**
   - Navigate to "APIs & Services" > "Library"
   - Search for "Cloud Vision API" and enable it

3. **Create Service Account Credentials:**
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "Service Account"
   - Fill in the details and create
   - Go to the service account's "Keys" tab
   - Click "Add Key" > "Create New Key" > "JSON"
   - Download the JSON file

4. **Configure credentials:**
   
   Save the credentials file outside the project:
   
   **Windows:**
   ```powershell
   # Create directory and save file
   mkdir C:\gcp
   # Move downloaded file to: C:\gcp\credentials.json
   
   # Set environment variable
   $env:GOOGLE_APPLICATION_CREDENTIALS="C:\gcp\credentials.json"
   ```
   
   **Mac/Linux:**
   ```bash
   # Create directory and save file
   mkdir -p ~/gcp
   # Move downloaded file to: ~/gcp/credentials.json
   
   # Set environment variable
   export GOOGLE_APPLICATION_CREDENTIALS="$HOME/gcp/credentials.json"
   ```

### 2. Environment Configuration

Create a `.env` file in the project root:

```env
AI_PROVIDER=google_vision
PORT=8080
ALLOWED_ORIGINS=http://localhost:4200,http://localhost:3000
MAX_FILE_SIZE=10485760
ALLOWED_FILE_TYPES=image/jpeg,image/png,image/gif,image/webp
AI_SERVICE_TIMEOUT=30s
HTTP_READ_TIMEOUT=10s
HTTP_WRITE_TIMEOUT=10s
HTTP_IDLE_TIMEOUT=60s
```

### 3. Run the Application

**Option A: With Docker (Recommended)**
```bash
docker-compose up --build
```

**Option B: Direct with Go**
```bash
go mod download
go run cmd/main.go
```

API will be available at: `http://localhost:8080`

## API Documentation

**Health Check:**
```
GET /health
```

**Analyze Image:**
```
POST /api/analyze
Content-Type: multipart/form-data

Body:
- image: file (JPEG, PNG, GIF, WEBP)

Response:
{
  "success": true,
  "data": {
    "tags": [{"label": "example", "confidence": 0.95}],
    "analyzed_at": "2025-12-22T10:00:00Z"
  }
}
```

## Features

- Image analysis using Google Cloud Vision API
- RESTful API with JSON responses
- File validation and size limits (10MB max)
- CORS support for frontend integration
- Health check endpoint

## Testing

```bash
go test ./...

# Run tests with coverage
go test -cover ./...
```

## License

MIT License
