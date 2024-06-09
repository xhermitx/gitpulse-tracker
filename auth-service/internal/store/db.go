package store

import "github.com/xhermitx/gitpulse-tracker/auth-service/internal/models"

type Store interface {
	CreateRecruiter(Recruiter *models.Recruiter) error
	AuthenticateRecruiter(credentilals *models.Credentials) (string, error)
	FindRecruiter(id int) (models.Recruiter, error)
}
