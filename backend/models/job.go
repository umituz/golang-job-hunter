package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Job struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Title       string    `json:"title" gorm:"not null"`
	Company     string    `json:"company" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Location    string    `json:"location"`
	Salary      string    `json:"salary"`
	URL         string    `json:"url" gorm:"unique;not null"`
	Source      string    `json:"source"` // RemoteOK, LinkedIn, etc.
	Remote      bool      `json:"remote" gorm:"default:false"`
	Keywords    string    `json:"keywords"`   // comma-separated keywords found
	HasKeywords bool      `json:"hasKeywords" gorm:"default:false"` // true if contains target keywords
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CheckKeywords checks if job contains target keywords (flutter, laravel, golang)
func (j *Job) CheckKeywords() {
	targetKeywords := []string{"flutter", "laravel", "golang", "go"}
	foundKeywords := []string{}
	
	searchText := strings.ToLower(j.Title + " " + j.Description + " " + j.Company)
	
	for _, keyword := range targetKeywords {
		if strings.Contains(searchText, strings.ToLower(keyword)) {
			foundKeywords = append(foundKeywords, keyword)
		}
	}
	
	if len(foundKeywords) > 0 {
		j.HasKeywords = true
		j.Keywords = strings.Join(foundKeywords, ",")
	} else {
		j.HasKeywords = false
		j.Keywords = ""
	}
}

// BeforeCreate hook to check keywords before saving
func (j *Job) BeforeCreate(tx *gorm.DB) error {
	j.CheckKeywords()
	return nil
}

// BeforeUpdate hook to check keywords before updating
func (j *Job) BeforeUpdate(tx *gorm.DB) error {
	j.CheckKeywords()
	return nil
}