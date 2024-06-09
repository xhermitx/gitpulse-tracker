package mysql

import (
	"errors"
	"fmt"
	"log"

	"github.com/xhermitx/gitpulse-tracker/auth-service/internal/app/services"
	"github.com/xhermitx/gitpulse-tracker/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MySQLStore struct {
	db *gorm.DB
}

func NewMySQLStore(db *gorm.DB) *MySQLStore {
	return &MySQLStore{db: db}
}

func (m *MySQLStore) CreateRecruiter(Recruiter *models.Recruiter) error {

	// CHECK IF THE USER ALREADY EXISTS
	res := m.db.First(&models.Recruiter{}, "username = ? OR email = ?", Recruiter.Username, Recruiter.Email)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {

			hashedPass, err := bcrypt.GenerateFromPassword([]byte(Recruiter.Password), 10)
			if err != nil {
				return err
			}

			Recruiter.Password = string(hashedPass)

			//CREATE THE USER
			res = m.db.Create(Recruiter)

			if res.Error != nil {
				return res.Error
			}

			log.Println("Rows Affected: ", res.RowsAffected)

			return nil
		}

		return res.Error
	}

	return fmt.Errorf("user already exists")
}

func (m *MySQLStore) AuthenticateRecruiter(credentials *models.Credentials) (string, error) {

	var recruiter models.Recruiter

	res := m.db.Where("username = ?", credentials.Username).First(&recruiter)
	if res.Error != nil {
		log.Println("Error finding the username:", res.Error)
		return "", res.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(recruiter.Password), []byte(credentials.Password))
	if err != nil {
		return "", err
	}

	return services.JwtAuth(recruiter.RecruiterId), nil
}

func (m *MySQLStore) FindRecruiter(id int) (models.Recruiter, error) {
	var recruiter models.Recruiter

	res := m.db.First(&recruiter, id)
	if res.Error != nil {
		return recruiter, res.Error
	}

	return recruiter, nil
}
