package database

import (
	"blizzflow/backend/domain/model"
	"errors"
	"log"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(filename ...string) error {
	if DB != nil {
		return errors.New("database already initialized")
	}

	// Set database path in parent directory
	dbFile := filepath.Join("C:\\vs_code_repos\\blizzflow", "data.db")
	if len(filename) > 0 {
		dbFile = filename[0]
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate after successful connection
	return Migrate()
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

func Migrate() error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	err := DB.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.SecurityQuestion{},
		&model.License{},
		&model.Inventory{},
		&model.Sale{},
	)

	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
