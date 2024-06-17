package models

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
