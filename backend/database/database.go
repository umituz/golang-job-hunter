package database

import (
	"job-hunter-backend/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	
	// Connect to SQLite database
	DB, err = gorm.Open(sqlite.Open("jobs.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.Job{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}

// GetFilteredJobs returns jobs that contain target keywords
func GetFilteredJobs(limit int, offset int) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	// Get total count of filtered jobs
	DB.Model(&models.Job{}).Where("has_keywords = ?", true).Count(&total)

	// Get filtered jobs with pagination
	result := DB.Where("has_keywords = ?", true).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&jobs)

	return jobs, total, result.Error
}

// SearchJobs searches for jobs by keywords and location
func SearchJobs(keywords []string, location string, limit int, offset int) ([]models.Job, int64, error) {
	var jobs []models.Job
	var total int64

	query := DB.Model(&models.Job{}).Where("has_keywords = ?", true)

	// Add location filter if provided
	if location != "" {
		query = query.Where("LOWER(location) LIKE ?", "%"+location+"%")
	}

	// Add keyword search in title, description, company
	if len(keywords) > 0 {
		for _, keyword := range keywords {
			query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(company) LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
		}
	}

	// Get total count
	query.Count(&total)

	// Get jobs with pagination
	result := query.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&jobs)

	return jobs, total, result.Error
}

// CreateJob creates a new job entry
func CreateJob(job *models.Job) error {
	result := DB.Create(job)
	return result.Error
}

// GetJobByID returns a job by its ID
func GetJobByID(id uint) (*models.Job, error) {
	var job models.Job
	result := DB.First(&job, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &job, nil
}