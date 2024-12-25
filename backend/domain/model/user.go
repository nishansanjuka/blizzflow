package model

import (
	"gorm.io/gorm"
)

// User represents a user profile in the system.
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
}

// CreateUser creates a new user in the database.
func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ValidatePassword checks if the provided password matches the stored password hash.
func (u *User) ValidatePassword(password string) bool {
	// Implement password validation logic (e.g., using bcrypt)
	return false
}
