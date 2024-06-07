package mysql

import "gorm.io/gorm"

type MySQLStore struct {
	db *gorm.DB
}

func NewMySQLStore(db *gorm.DB) *MySQLStore {
	return &MySQLStore{db: db}
}
