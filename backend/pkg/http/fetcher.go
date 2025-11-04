// Package http provides HTTP utilities for fetching external data sources.
// Per AGENT_TECH_SPEC.md §9: retry with backoff, circuit-break on repeated failures.
package http

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// FetcherConfig holds configuration for HTTP fetching.
type FetcherConfig struct {
	MaxRetries     int           // Maximum number of retries (default: 3)
	InitialBackoff time.Duration // Initial backoff duration (default: 500ms)
	MaxBackoff     time.Duration // Maximum backoff duration (default: 30s)
	Timeout        time.Duration // HTTP request timeout (default: 30s)
	UserAgent      string        // User-Agent header
}

// DefaultConfig returns sensible defaults for fetching public data.
func DefaultConfig() FetcherConfig {
	return FetcherConfig{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     30 * time.Second,
		Timeout:        30 * time.Second,
		// Mimic modern Chrome browser to avoid bot detection
		UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// BrowserConfig returns configuration that mimics a real browser.
func BrowserConfig() FetcherConfig {
	return FetcherConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		Timeout:        45 * time.Second,
		// Latest Chrome on macOS
		UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// Fetcher performs HTTP requests with retry logic.
type Fetcher struct {
	config FetcherConfig
	client *http.Client
}

// NewFetcher creates a new HTTP fetcher with the given configuration.
func NewFetcher(config FetcherConfig) *Fetcher {
	return &Fetcher{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Fetch performs an HTTP GET request with exponential backoff retry.
// Returns the response body as an io.ReadCloser (caller must close).
func (f *Fetcher) Fetch(url string) (io.ReadCloser, error) {
	var lastErr error

	for attempt := 0; attempt <= f.config.MaxRetries; attempt++ {
		// Apply backoff before retry (skip on first attempt)
		if attempt > 0 {
			backoff := f.calculateBackoff(attempt)
			time.Sleep(backoff)
		}

		// Create request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set browser-like headers to avoid bot detection
		if f.config.UserAgent != "" {
			req.Header.Set("User-Agent", f.config.UserAgent)
		}
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Accept-Language", "ja-JP,ja;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("DNT", "1")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Cache-Control", "max-age=0")

		// Execute request
		resp, err := f.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed (attempt %d/%d): %w",
				attempt+1, f.config.MaxRetries+1, err)
			continue
		}

		// Check status code
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d from %s (attempt %d/%d)",
				resp.StatusCode, url, attempt+1, f.config.MaxRetries+1)
			continue
		}

		// Check for gzip encoding and decompress if needed
		body := resp.Body
		if strings.Contains(strings.ToLower(resp.Header.Get("Content-Encoding")), "gzip") {
			gzipReader, err := gzip.NewReader(resp.Body)
			if err != nil {
				resp.Body.Close()
				lastErr = fmt.Errorf("failed to create gzip reader: %w", err)
				continue
			}
			// Return a wrapper that closes both gzip reader and original body
			body = &gzipReadCloser{
				Reader: gzipReader,
				body:   resp.Body,
			}
		}

		// Success
		return body, nil
	}

	// All retries exhausted
	return nil, fmt.Errorf("fetch failed after %d attempts: %w",
		f.config.MaxRetries+1, lastErr)
}

// calculateBackoff computes exponential backoff duration.
// Formula: min(InitialBackoff * 2^attempt, MaxBackoff)
func (f *Fetcher) calculateBackoff(attempt int) time.Duration {
	backoff := float64(f.config.InitialBackoff) * math.Pow(2, float64(attempt-1))
	if backoff > float64(f.config.MaxBackoff) {
		backoff = float64(f.config.MaxBackoff)
	}
	return time.Duration(backoff)
}

// FetchFromZip downloads a ZIP file and extracts a specific file by name pattern.
// Returns the extracted file contents as an io.ReadCloser.
// Pattern can use filepath.Match syntax (e.g., "202511*.csv").
func (f *Fetcher) FetchFromZip(zipURL, filePattern string) (io.ReadCloser, error) {
	// Fetch the ZIP file
	zipBody, err := f.Fetch(zipURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ZIP: %w", err)
	}
	defer zipBody.Close()

	// Read ZIP contents into memory
	zipData, err := io.ReadAll(zipBody)
	if err != nil {
		return nil, fmt.Errorf("failed to read ZIP data: %w", err)
	}

	// Open ZIP archive
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("failed to open ZIP archive: %w", err)
	}

	// Find matching file in ZIP
	for _, file := range zipReader.File {
		matched, err := filepath.Match(filePattern, file.Name)
		if err != nil {
			continue
		}
		if matched {
			// Open the file
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s in ZIP: %w", file.Name, err)
			}

			// Read file contents into memory (since ZIP reader needs random access)
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s from ZIP: %w", file.Name, err)
			}

			// Return as ReadCloser
			return io.NopCloser(bytes.NewReader(data)), nil
		}
	}

	return nil, fmt.Errorf("no file matching pattern %q found in ZIP", filePattern)
}

// gzipReadCloser wraps a gzip.Reader and ensures both the gzip reader
// and underlying response body are properly closed.
type gzipReadCloser struct {
	*gzip.Reader
	body io.ReadCloser
}

// Close closes both the gzip reader and underlying response body.
func (g *gzipReadCloser) Close() error {
	// Close gzip reader first
	if err := g.Reader.Close(); err != nil {
		g.body.Close() // Still close body even if gzip close fails
		return err
	}
	// Close underlying response body
	return g.body.Close()
}
