package models

type Candidate struct {
	RedisCandidate
	Status bool
}

type RedisCandidate struct {
	JobID           uint
	Username        string
	Followers       int
	Contributions   int
	MostPopularRepo string
	RepoStars       int
}
