// Database-aware version of API server
// This file will replace main.go once Railway PostgreSQL is configured
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/teo/aversome/backend/internal/storage"
	"github.com/teo/aversome/backend/pkg/database"
)

var dbStorage *storage.DataStorage

func mainWithDB() {
	// Initialize database connection
	dbConfig := database.DefaultConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Println("📁 Falling back to filesystem mode")
		main() // Fall back to file-based version
		return
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize storage
	dbStorage = storage.NewDataStorage(db)

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://localhost:5174",
		"https://japan-energy-dashboard.vercel.app",
	}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	router.Use(cors.New(config))

	// Health check with DB status
	router.GET("/api/health", func(c *gin.Context) {
		dbHealthy := true
		if err := db.HealthCheck(); err != nil {
			dbHealthy = false
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"time":        time.Now().Format(time.RFC3339),
			"database":    dbHealthy,
			"storage":     "postgresql",
		})
	})

	// Data endpoints
	router.POST("/api/data/refresh", handleRefreshDB)
	router.GET("/api/demand/:area/:date", handleGetDemandDB)
	router.GET("/api/jepx/:area/:date", handleGetJEPXDB)
	router.GET("/api/reserve/:date", handleGetReserveDB)

	// Stats endpoint
	router.GET("/api/stats", func(c *gin.Context) {
		count, _ := dbStorage.GetDataCount()
		c.JSON(http.StatusOK, gin.H{
			"total_records": count,
			"storage_type":  "postgresql",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API server starting on http://localhost:%s", port)
	log.Printf("🗄️  Storage mode: PostgreSQL")
	log.Printf("📊 Data refresh endpoint: POST /api/data/refresh")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleRefreshDB(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = RefreshRequest{
			Date:  time.Now().Format("2006-01-02"),
			Areas: []string{"tokyo", "kansai"},
		}
	}

	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}
	if len(req.Areas) == 0 {
		req.Areas = []string{"tokyo", "kansai"}
	}

	log.Printf("📥 Refresh request: date=%s, areas=%v", req.Date, req.Areas)

	var results []DataFetchResult

	// Parse date
	targetDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Fetch demand data for each area
	for _, area := range req.Areas {
		// Demand data
		demandResult := fetchDemandDB(area, req.Date, targetDate)
		results = append(results, demandResult)

		// JEPX data
		jepxResult := fetchJEPXDB(area, req.Date, targetDate)
		results = append(results, jepxResult)
	}

	// Reserve data
	reserveResult := fetchReserveDB(req.Date, targetDate)
	results = append(results, reserveResult)

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

func fetchDemandDB(area, dateStr string, date time.Time) DataFetchResult {
	start := time.Now()
	source := fmt.Sprintf("%s-demand", area)

	// Execute fetch command
	cmd := exec.Command(
		"./fetch-demand",
		"-area", area,
		"-date", dateStr,
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

	// Read the generated JSON and save to database
	var jsonData interface{}
	if err := json.Unmarshal(output, &jsonData); err != nil {
		// Try reading from file if command didn't output JSON
		filePath := fmt.Sprintf("public/data/jp/%s/demand-%s.json", area, dateStr)
		fileData, readErr := os.ReadFile(filePath)
		if readErr == nil {
			json.Unmarshal(fileData, &jsonData)
		}
	}

	// Save to database
	if jsonData != nil {
		areaPtr := &area
		if err := dbStorage.SaveData("demand", areaPtr, date, jsonData); err != nil {
			log.Printf("⚠️  Failed to save to DB: %v", err)
		}
	}

	return DataFetchResult{
		Source:   source,
		Status:   "success",
		FilePath: fmt.Sprintf("database://demand/%s/%s", area, dateStr),
		Duration: duration.String(),
	}
}

func fetchJEPXDB(area, dateStr string, date time.Time) DataFetchResult {
	start := time.Now()
	source := fmt.Sprintf("%s-jepx", area)

	cmd := exec.Command(
		"./fetch-jepx",
		"-area", area,
		"-date", dateStr,
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

	// Save to database
	var jsonData interface{}
	if err := json.Unmarshal(output, &jsonData); err == nil {
		areaPtr := &area
		if err := dbStorage.SaveData("jepx", areaPtr, date, jsonData); err != nil {
			log.Printf("⚠️  Failed to save to DB: %v", err)
		}
	}

	return DataFetchResult{
		Source:   source,
		Status:   "success",
		FilePath: fmt.Sprintf("database://jepx/%s/%s", area, dateStr),
		Duration: duration.String(),
	}
}

func fetchReserveDB(dateStr string, date time.Time) DataFetchResult {
	start := time.Now()
	source := "reserve"

	cmd := exec.Command(
		"./fetch-reserve",
		"-date", dateStr,
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

	// Save to database (area is NULL for system-wide data)
	var jsonData interface{}
	if err := json.Unmarshal(output, &jsonData); err == nil {
		if err := dbStorage.SaveData("reserve", nil, date, jsonData); err != nil {
			log.Printf("⚠️  Failed to save to DB: %v", err)
		}
	}

	return DataFetchResult{
		Source:   source,
		Status:   "success",
		FilePath: fmt.Sprintf("database://reserve/%s", dateStr),
		Duration: duration.String(),
	}
}

func handleGetDemandDB(c *gin.Context) {
	area := c.Param("area")
	dateStr := c.Param("date")

	if area != "tokyo" && area != "kansai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Try to get from database
	areaPtr := &area
	data, err := dbStorage.GetData("demand", areaPtr, date)
	if err != nil {
		// Data not found - fetch fresh
		log.Printf("[GET /api/demand] Data not found in DB, fetching fresh for %s/%s", area, dateStr)
		result := fetchDemandDB(area, dateStr, date)

		if result.Status != "success" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch demand data"})
			return
		}

		// Try again to get from DB
		data, err = dbStorage.GetData("demand", areaPtr, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data fetch succeeded but read failed"})
			return
		}
	}

	c.Data(http.StatusOK, "application/json", data)
}

func handleGetJEPXDB(c *gin.Context) {
	area := c.Param("area")
	dateStr := c.Param("date")

	if area != "tokyo" && area != "kansai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	areaPtr := &area
	data, err := dbStorage.GetData("jepx", areaPtr, date)
	if err != nil {
		log.Printf("[GET /api/jepx] Data not found in DB, fetching fresh for %s/%s", area, dateStr)
		result := fetchJEPXDB(area, dateStr, date)

		if result.Status != "success" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch JEPX data"})
			return
		}

		data, err = dbStorage.GetData("jepx", areaPtr, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data fetch succeeded but read failed"})
			return
		}
	}

	c.Data(http.StatusOK, "application/json", data)
}

func handleGetReserveDB(c *gin.Context) {
	dateStr := c.Param("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	data, err := dbStorage.GetData("reserve", nil, date)
	if err != nil {
		log.Printf("[GET /api/reserve] Data not found in DB, fetching fresh for %s", dateStr)
		result := fetchReserveDB(dateStr, date)

		if result.Status != "success" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reserve data"})
			return
		}

		data, err = dbStorage.GetData("reserve", nil, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Data fetch succeeded but read failed"})
			return
		}
	}

	c.Data(http.StatusOK, "application/json", data)
}
