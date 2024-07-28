package store

import "github.com/xhermitx/gitpulse-tracker/backend/internal/models"

type Store interface {
	CreateJob(Job *models.Job) (*models.Job, error)
	DeleteJob(JobId uint) error
	UpdateJob(Job *models.Job) (*models.Job, error)
	GetJob(JobId uint) (*models.Job, error)
	ListJobs(RecruiterId uint) ([]*models.Job, error)
	ListCandidates(JobId uint) ([]*models.TopCandidates, error)
}
