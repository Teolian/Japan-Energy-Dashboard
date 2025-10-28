// Package main provides a lightweight API server for data refresh operations.
// Usage: go run cmd/api/main.go
// Endpoint: POST /api/data/refresh
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	Date  string   `json:"date"`  // Optional, defaults to today
	Areas []string `json:"areas"` // Optional, defaults to ["tokyo", "kansai"]
}

type RefreshResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Results []DataFetchResult `json:"results"`
}

type DataFetchResult struct {
	Source   string `json:"source"`   // e.g. "tokyo-demand", "kansai-jepx"
	Status   string `json:"status"`   // "success", "error"
	FilePath string `json:"file_path,omitempty"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration"`
}

func main() {
	// Set Gin to release mode for production
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS configuration for frontend
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Data refresh endpoint
	router.POST("/api/data/refresh", handleRefresh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API server starting on http://localhost:%s", port)
	log.Printf("📊 Data refresh endpoint: POST /api/data/refresh")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleRefresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no body provided, use defaults
		req = RefreshRequest{
			Date:  time.Now().Format("2006-01-02"),
			Areas: []string{"tokyo", "kansai"},
		}
	}

	// Default date to today if not provided
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	// Default areas if not provided
	if len(req.Areas) == 0 {
		req.Areas = []string{"tokyo", "kansai"}
	}

	log.Printf("📥 Refresh request: date=%s, areas=%v", req.Date, req.Areas)

	var results []DataFetchResult

	// Fetch demand data for each area
	for _, area := range req.Areas {
		// Demand data
		demandResult := fetchDemand(area, req.Date)
		results = append(results, demandResult)

		// JEPX data
		jepxResult := fetchJEPX(area, req.Date)
		results = append(results, jepxResult)
	}

	// Fetch reserve data (system-wide)
	reserveResult := fetchReserve(req.Date)
	results = append(results, reserveResult)

	// Check if all succeeded
	allSuccess := true
	for _, result := range results {
		if result.Status != "success" {
			allSuccess = false
			break
		}
	}

	response := RefreshResponse{
		Success: allSuccess,
		Message: fmt.Sprintf("Fetched data for %s", req.Date),
		Results: results,
	}

	if allSuccess {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusPartialContent, response)
	}
}

func fetchDemand(area, date string) DataFetchResult {
	start := time.Now()
	source := fmt.Sprintf("%s-demand", area)

	cmd := exec.Command(
		"go", "run",
		"cmd/fetch-demand-http/main.go",
		"-area", area,
		"-date", date,
		"--use-http",
	)

	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		return DataFetchResult{
			Source:   source,
			Status:   "error",
			Error:    fmt.Sprintf("%v: %s", err, string(output)),
			Duration: duration.String(),
		}
	}

	// Parse output path from command output
	outputPath := filepath.Join("public", "data", "jp", area, fmt.Sprintf("demand-%s.json", date))

	return DataFetchResult{
		Source:   source,
		Status:   "success",
		FilePath: outputPath,
		Duration: duration.String(),
	}
}

func fetchJEPX(area, date string) DataFetchResult {
	start := time.Now()
	source := fmt.Sprintf("%s-jepx", area)

	cmd := exec.Command(
		"go", "run",
		"cmd/fetch-jepx/main.go",
		"-area", area,
		"-date", date,
	)

	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		return DataFetchResult{
			Source:   source,
			Status:   "error",
			Error:    fmt.Sprintf("%v: %s", err, string(output)),
			Duration: duration.String(),
		}
	}

	outputPath := filepath.Join("public", "data", "jp", "jepx", fmt.Sprintf("spot-%s-%s.json", area, date))

	return DataFetchResult{
		Source:   source,
		Status:   "success",
		FilePath: outputPath,
		Duration: duration.String(),
	}
}

func fetchReserve(date string) DataFetchResult {
	start := time.Now()
	source := "reserve"

	cmd := exec.Command(
		"go", "run",
		"cmd/fetch-reserve/main.go",
		"-date", date,
	)

	output, err := cmd.CombinedOutput()
	duration := time.Since(start)

	if err != nil {
		return DataFetchResult{
			Source:   source,
			Status:   "error",
			Error:    fmt.Sprintf("%v: %s", err, string(output)),
			Duration: duration.String(),
		}
	}

	outputPath := filepath.Join("public", "data", "jp", "system", fmt.Sprintf("reserve-%s.json", date))

	return DataFetchResult{
		Source:   source,
		Status:   "success",
		FilePath: outputPath,
		Duration: duration.String(),
	}
}
