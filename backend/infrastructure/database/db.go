package database

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(filename ...string) error {
	dbFile := "data.db"
	if len(filename) > 0 {
		dbFile = filename[0]
	}
	if DB != nil {
		return errors.New("database already initialized")
	}
	var err error
	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	return err
}

func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
