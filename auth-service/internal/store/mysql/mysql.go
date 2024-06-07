package mysql

import (
	"errors"
	"fmt"
	"log"

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

func (m *MySQLStore) AuthenticateRecruiter(username string, password string) (*models.Recruiter, error) {

	var recruiter models.Recruiter

	res := m.db.First(&recruiter, "username = ?", username)
	if res.Error != nil {
		return nil, res.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(recruiter.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("incorrect username or password")
	}

	return &recruiter, nil
}
