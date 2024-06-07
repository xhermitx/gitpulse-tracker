package models

import "time"

type Recruiter struct {
	RecruiterId int    `gorm:"primary_key;AUTO_INCREMENT"`
	Username    string `gorm:"unique, not null"`
	Password    string `gorm:"unique, not null"`
	Email       string
	CreatedAt   time.Time
}
