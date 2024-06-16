package models

type GitQuery struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

type GitResponse struct {
	Data struct {
		User struct {
			Login                   string `json:"login"`
			Typename                string `json:"__typename"`
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
			Followers struct {
				TotalCount int `json:"totalCount"`
			} `json:"followers"`
			Repositories struct {
				Nodes []struct {
					Name           string `json:"name"`
					StargazerCount int    `json:"stargazerCount"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type Job struct {
	JobID     uint
	Usernames []string
}

type Candidate struct {
	CandidateId     uint   `gorm:"not null; unique"`
	GithubId        string `gorm:"size: 255;not null"`
	Followers       uint
	Contributions   uint
	MostPopularRepo string `gorm:"size:255"`
	RepoStars       uint
	Score           uint
	JobId           uint `gorm:"not null; unique; index"`
	Status          bool
}
