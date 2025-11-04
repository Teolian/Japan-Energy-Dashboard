# Multi-stage build for Japan Energy Dashboard API

# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Set Go toolchain to local to avoid version conflicts
ENV GOTOOLCHAIN=local

# Copy all source code first (needed for go mod tidy)
COPY . .

# Sync dependencies
RUN go mod tidy

# Download dependencies
RUN go mod download

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o api cmd/api/main.go

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies (including curl for healthcheck)
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/api .

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set timezone to JST
ENV TZ=Asia/Tokyo

# Expose port
EXPOSE 8080

# Health check using curl
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/health || exit 1

# Run
CMD ["./api"]
