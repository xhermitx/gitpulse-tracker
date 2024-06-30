package models

import "time"

type Recruiter struct {
	RecruiterId uint      `gorm:"primary_key;AUTO_INCREMENT"`
	Username    string    `gorm:"unique, not null"`
	Password    string    `gorm:"not null"`
	Email       string    `gorm:"unique"`
	Company     string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"type:datetime"`
}

type Credentials struct {
	Username string `json: string`
	Password string `json: string`
}
