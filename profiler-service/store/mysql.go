package store

import (
	"fmt"
	"log"
	"os"

	"github.com/xhermitx/gitpulse-tracker/profiler-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InsertData(topCandidates []models.RedisCandidate) error {

	fmt.Println("DSN: ", os.Getenv("DB_SERVER"))

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_SERVER")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to DB")
	}

	//INSERT THE CANDIDATES
	res := db.Create(topCandidates)

	if res.Error != nil {
		return res.Error
	}

	log.Println("Rows Affected: ", res.RowsAffected)

	return nil
}
