package scrapers

import (
	"encoding/json"
	"fmt"
	"io"
	"job-hunter-backend/models"
	"net/http"
	"strings"
	"time"
)

type RemoteOKJob struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Company     string `json:"company"`
	CompanyLogo string `json:"company_logo"`
	Position    string `json:"position"`
	Tags        []string `json:"tags"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Remote      bool `json:"remote"`
	URL         string `json:"url"`
	Date        string `json:"date"`
	Salary      interface{} `json:"salary_min"`
}

type RemoteOKScraper struct {
	BaseURL string
}

func NewRemoteOKScraper() *RemoteOKScraper {
	return &RemoteOKScraper{
		BaseURL: "https://remoteok.io/api",
	}
}

func (r *RemoteOKScraper) ScrapeJobs() ([]models.Job, error) {
	resp, err := http.Get(r.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var remoteOKJobs []RemoteOKJob
	err = json.Unmarshal(body, &remoteOKJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var jobs []models.Job
	for _, remoteJob := range remoteOKJobs {
		// Skip the first item which is usually metadata
		if remoteJob.ID == "" {
			continue
		}

		// Convert salary interface to string
		salaryStr := ""
		if remoteJob.Salary != nil {
			salaryStr = fmt.Sprintf("%v", remoteJob.Salary)
		}

		// Convert RemoteOK job to our Job model
		job := models.Job{
			Title:       remoteJob.Position,
			Company:     remoteJob.Company,
			Description: cleanDescription(remoteJob.Description),
			Location:    remoteJob.Location,
			Salary:      salaryStr,
			URL:         fmt.Sprintf("https://remoteok.io/remote-jobs/%s", remoteJob.Slug),
			Source:      "RemoteOK",
			Remote:      remoteJob.Remote,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Add tags to description for better keyword matching
		if len(remoteJob.Tags) > 0 {
			job.Description += " Technologies: " + strings.Join(remoteJob.Tags, ", ")
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// cleanDescription removes HTML tags and cleans up the description
func cleanDescription(description string) string {
	// Remove common HTML tags
	description = strings.ReplaceAll(description, "<br>", "\n")
	description = strings.ReplaceAll(description, "<br/>", "\n")
	description = strings.ReplaceAll(description, "<p>", "")
	description = strings.ReplaceAll(description, "</p>", "\n")
	description = strings.ReplaceAll(description, "<div>", "")
	description = strings.ReplaceAll(description, "</div>", "\n")
	
	// Remove other HTML tags (basic cleanup)
	for strings.Contains(description, "<") && strings.Contains(description, ">") {
		start := strings.Index(description, "<")
		end := strings.Index(description, ">")
		if start < end {
			description = description[:start] + description[end+1:]
		} else {
			break
		}
	}

	// Clean up multiple newlines and spaces
	description = strings.ReplaceAll(description, "\n\n\n", "\n\n")
	description = strings.TrimSpace(description)

	return description
}

// ScrapeAndFilter scrapes jobs and returns only those with target keywords
func (r *RemoteOKScraper) ScrapeAndFilter() ([]models.Job, error) {
	allJobs, err := r.ScrapeJobs()
	if err != nil {
		return nil, err
	}

	var filteredJobs []models.Job
	for _, job := range allJobs {
		// CheckKeywords is called automatically in BeforeCreate hook
		// but we call it here to check before adding to filtered list
		job.CheckKeywords()
		if job.HasKeywords {
			filteredJobs = append(filteredJobs, job)
		}
	}

	return filteredJobs, nil
}