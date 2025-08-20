package main

import (
	"job-hunter-backend/database"
	"job-hunter-backend/models"
	"job-hunter-backend/scrapers"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDatabase()

	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:5000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Job Hunter API is running",
			"time":    time.Now(),
		})
	})

	// API routes group
	api := r.Group("/api/v1")
	{
		api.GET("/jobs", getJobs)
		api.POST("/search", searchJobs)
		api.GET("/jobs/:id", getJobByID)
		api.POST("/scrape", scrapeJobs)
		api.GET("/stats", getStats)
	}

	log.Println("Server starting on port 8081...")
	r.Run(":8081")
}

func getJobs(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	offset := (page - 1) * limit

	jobs, total, err := database.GetFilteredJobs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":    jobs,
		"total":   total,
		"page":    page,
		"limit":   limit,
		"pages":   (total + int64(limit) - 1) / int64(limit),
	})
}

func searchJobs(c *gin.Context) {
	var searchRequest struct {
		Keywords []string `json:"keywords"`
		Location string   `json:"location"`
		Page     int      `json:"page"`
		Limit    int      `json:"limit"`
	}

	if err := c.ShouldBindJSON(&searchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}
	if searchRequest.Limit < 1 || searchRequest.Limit > 100 {
		searchRequest.Limit = 20
	}

	offset := (searchRequest.Page - 1) * searchRequest.Limit

	jobs, total, err := database.SearchJobs(
		searchRequest.Keywords,
		searchRequest.Location,
		searchRequest.Limit,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":         jobs,
		"total":        total,
		"page":         searchRequest.Page,
		"limit":        searchRequest.Limit,
		"pages":        (total + int64(searchRequest.Limit) - 1) / int64(searchRequest.Limit),
		"searchParams": searchRequest,
	})
}

func getJobByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := database.GetJobByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job": job})
}

func scrapeJobs(c *gin.Context) {
	// Initialize RemoteOK scraper
	scraper := scrapers.NewRemoteOKScraper()

	// Scrape filtered jobs
	jobs, err := scraper.ScrapeAndFilter()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scrape jobs: " + err.Error(),
		})
		return
	}

	// Save jobs to database
	savedCount := 0
	errors := []string{}

	for _, job := range jobs {
		if err := database.CreateJob(&job); err != nil {
			// Skip duplicates (URL unique constraint)
			if !containsString(err.Error(), "UNIQUE constraint failed") {
				errors = append(errors, err.Error())
			}
		} else {
			savedCount++
		}
	}

	response := gin.H{
		"message":      "Scraping completed",
		"scraped":      len(jobs),
		"saved":        savedCount,
		"duplicates":   len(jobs) - savedCount - len(errors),
		"time":         time.Now(),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

func getStats(c *gin.Context) {
	var total int64
	var withKeywords int64

	database.DB.Model(&models.Job{}).Count(&total)
	database.DB.Model(&models.Job{}).Where("has_keywords = ?", true).Count(&withKeywords)

	c.JSON(http.StatusOK, gin.H{
		"totalJobs":          total,
		"jobsWithKeywords":   withKeywords,
		"keywordsPercentage": float64(withKeywords) / float64(total) * 100,
		"lastUpdated":        time.Now(),
	})
}

func containsString(slice string, item string) bool {
	return len(slice) > 0 && len(item) > 0 && 
		   (slice == item || 
		    (len(slice) > len(item) && slice[:len(item)] == item) ||
		    (len(slice) > len(item) && slice[len(slice)-len(item):] == item))
}