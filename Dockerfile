# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o golog .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

# Create non-root user
RUN addgroup -g 1000 -S golog && \
    adduser -u 1000 -S golog -G golog

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/golog .

# Create data directory
RUN mkdir -p /app/data && chown -R golog:golog /app

# Switch to non-root user
USER golog

# Expose port
EXPOSE 8080

# Set environment variables
ENV HOST=0.0.0.0
ENV PORT=8080
ENV ENABLE_UI=true
ENV GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/sessions || exit 1

# Run the application
CMD ["./golog"]