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

# Debug: List copied files
RUN echo "=== Listing /build directory ===" && ls -la
RUN echo "=== Listing /build/cmd ===" && ls -la cmd/ || echo "cmd/ not found"

# Force Go 1.23 in go.mod (remove toolchain directive)
RUN sed -i '/^toolchain/d' go.mod && \
    sed -i 's/^go 1\.24.*/go 1.23/' go.mod

# Sync dependencies
RUN go mod tidy

# Download dependencies
RUN go mod download

# Build main API binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o api ./cmd/api

# Build fetch command binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o fetch-demand ./cmd/fetch-demand-http
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o fetch-jepx ./cmd/fetch-jepx
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o fetch-reserve ./cmd/fetch-reserve

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies (including curl for healthcheck)
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/api .
COPY --from=builder /build/fetch-demand .
COPY --from=builder /build/fetch-jepx .
COPY --from=builder /build/fetch-reserve .

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Create data directories for fetched data
RUN mkdir -p public/data/jp/tokyo public/data/jp/kansai public/data/jp/jepx public/data/jp/system

# Set timezone to JST
ENV TZ=Asia/Tokyo

# Expose port
EXPOSE 8080

# Health check using curl
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/health || exit 1

# Run
CMD ["./api"]
