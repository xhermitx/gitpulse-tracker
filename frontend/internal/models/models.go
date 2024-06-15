package models

import (
	"time"

	"gorm.io/gorm"
)

// Recruiter represents a user in the system (recruiter or candidate)
type Recruiter struct {
	RecruiterId uint      `gorm:"primary_key;AUTO_INCREMENT"`
	Username    string    `gorm:"unique, not null"`
	Password    string    `gorm:"unique, not null"`
	Email       string    `gorm:"unique"`
	Company     string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"type:datetime"`
}

// Job represents a job posting by a recruiter
type Job struct {
	JobId       uint      `json:"job_id" gorm:"primaryKey;autoIncrement"`
	JobName     string    `json:"job_name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text; not null"`
	DriveLink   string    `json:"drive_link" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:dateTime"`
	RecruiterId uint      `json:"recruiter_id" gorm:"not null"`
}

// Candidates List
type CandidatesList struct {
	JobId    uint   `gorm:"not null"`
	GithubId string `gorm:"not null"`
}

// Top Candidates Data
type TopCandidates struct {
	CandidateId     uint   `gorm:"not null; unique"`
	GithubId        string `gorm:"size: 255;not null"`
	Followers       uint
	Contributions   uint
	MostPopularRepo string `gorm:"size:255"`
	RepoStars       uint
	Score           uint
	JobId           uint `gorm:"not null; unique; index"`
}

// ----------TO BE IMPLEMENTED----------------

// AuthToken represents an authentication token for a user
type AuthToken struct {
	RecruiterId uint   `gorm:"not null;index"`
	AuthToken   string `gorm:"size:255;not null;unique"`
	ExpiresAt   time.Time
}

// PasswordReset represents a password reset token for a user
type PasswordReset struct {
	gorm.Model
	RecruiterId uint   `gorm:"not null;index"`
	ResetToken  string `gorm:"size:255;not null;unique"`
	ExpiresAt   time.Time
}
