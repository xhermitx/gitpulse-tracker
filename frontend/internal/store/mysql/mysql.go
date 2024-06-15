package mysql

import (
	"errors"
	"fmt"
	"log"

	"github.com/xhermitx/gitpulse-tracker/frontend/internal/models"
	"gorm.io/gorm"
)

type MySQLStore struct {
	db *gorm.DB
}

func NewMySQLStore(db *gorm.DB) *MySQLStore {
	return &MySQLStore{db: db}
}

func (m MySQLStore) CreateJob(Job *models.Job) (*models.Job, error) {
	// CHECK IF THE JOB NAME ALREADY EXISTS
	res := m.db.First(&models.Job{}, "job_name = ? AND recruiter_id = ?", Job.JobName, Job.RecruiterId)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {

			//CREATE THE JOB
			res = m.db.Create(Job)

			if res.Error != nil {
				return nil, res.Error
			}

			log.Println("Rows Affected: ", res.RowsAffected)

			return Job, nil
		}

		return nil, res.Error
	}

	return nil, fmt.Errorf("job name already exists")
}

func (m MySQLStore) UpdateJob(Job *models.Job) (*models.Job, error) {
	return &models.Job{}, nil
}

func (m MySQLStore) DeleteJob(JobId uint) error {
	return nil
}

func (m MySQLStore) ListJobs(RecruiterId uint) ([]*models.Job, error) {

	var jobs []*models.Job

	res := m.db.Find(&jobs, "recruiter_id = ?", RecruiterId)
	if res.Error != nil {
		return nil, res.Error
	}

	return jobs, nil
}
