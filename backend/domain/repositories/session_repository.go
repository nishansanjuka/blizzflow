package repository

import (
	"blizzflow/backend/domain/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(session *model.Session) error {
	result := r.db.Create(session)
	return result.Error
}

func (r *SessionRepository) GetSession(id uint) (*model.Session, error) {
	var session model.Session
	result := r.db.First(&session, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &session, result.Error
}

func (r *SessionRepository) DeleteSession(id uint) error {
	result := r.db.Delete(&model.Session{}, id)
	return result.Error
}

func (r *SessionRepository) GetSessionByUserID(userID uint) (*model.Session, error) {
	var session model.Session
	result := r.db.Where("user_id = ?", userID).First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &session, result.Error
}

func (r *SessionRepository) CleanupExpiredSessions() error {
	expirationTime := time.Now().Add(-24 * time.Hour) // Example: 24 hours expiration
	result := r.db.Where("created_at < ?", expirationTime).Delete(&model.Session{})
	return result.Error
}
