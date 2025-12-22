# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies (if any)
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin ./cmd

# Stage 2: Runtime
FROM alpine:3.23.2

# Install CA certificates for HTTPS requests to external APIs
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app /app

# Copy credentials file if it exists (for local development)
# The asterisk makes it optional - won't fail if file doesn't exist
COPY credentials.json* /app/credentials.json

# Set default credentials path (can be overridden by environment variable)
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/bin"]
