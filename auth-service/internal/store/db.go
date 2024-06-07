package store

import "github.com/xhermitx/gitpulse-tracker/auth-service/internal/models"

type Store interface {
	CreateRecruiter(Recruiter *models.Recruiter) error
	AuthenticateRecruiter(username string, password string) (*models.Recruiter, error)
}
