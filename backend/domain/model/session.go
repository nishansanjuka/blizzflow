package model

import (
	"time"
)

// Session represents a user session in the application.
type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// NewSession creates a new session for a user.
func NewSession(userID uint) *Session {
	return &Session{
		UserID: userID,
	}
}
