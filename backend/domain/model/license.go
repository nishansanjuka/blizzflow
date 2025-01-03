package model

import "time"

type License struct {
	ID          uint      `gorm:"primaryKey"`
	Key         string    `gorm:"unique;not null"`
	Username    string    `gorm:"not null"`
	ExpiresAt   time.Time `gorm:"not null"`
	Fingerprint string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
