package store

import "github.com/xhermitx/gitpulse-tracker/auth-service/internal/models"

type Store interface {
	CreateRecruiter(Recruiter *models.Recruiter) error
	DeleteRecruiter(RecruiterId int) error
	UpdatePassword(NewPassword string, Recruiter *models.Recruiter) error
	ViewRecruiters() ([]models.Recruiter, error)
}
