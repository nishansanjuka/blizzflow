package model

type SecurityQuestion struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"not null"`
	Question string `gorm:"not null"`
	Answer   string `gorm:"not null"`
}
