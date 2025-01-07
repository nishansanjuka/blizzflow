package model

import "time"

type Sale struct {
	ID          uint      `gorm:"primaryKey"`
	InventoryID uint      `gorm:"not null"`
	Quantity    int       `gorm:"not null"`
	TotalPrice  float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
