package models

import "time"

type Candidate struct {
	TopCandidates
	Status bool
}

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

var (
	STATUS_QUEUE      = "profiling_status_queue"
	USERNAME_QUEUE    = "username_queue"
	GITHUB_DATA_QUEUE = "github_data_queue"
)

type StatusQueue struct {
	JobId  uint
	Status bool
	Timer  time.Time
}
