package session_service

import (
	"blizzflow/backend/domain/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Custom errors
var (
	ErrSessionNotFound   = fmt.Errorf("session not found")
	ErrDatabaseOperation = fmt.Errorf("database operation failed")
	ErrInvalidSessionID  = fmt.Errorf("invalid session ID")
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) CreateSession(userID uint) (*model.Session, error) {
	session := &model.Session{
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", ErrDatabaseOperation)
	}
	return session, nil
}

func (s *SessionService) GetSession(sessionID uint) (*model.Session, error) {
	if sessionID == 0 {
		return nil, ErrInvalidSessionID
	}

	var session model.Session
	if err := s.db.First(&session, sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to retrieve session: %w", ErrDatabaseOperation)
	}
	return &session, nil
}

func (s *SessionService) DeleteSession(sessionID uint) error {
	if sessionID == 0 {
		return ErrInvalidSessionID
	}

	if err := s.db.Delete(&model.Session{}, sessionID).Error; err != nil {
		return fmt.Errorf("failed to delete session: %w", ErrDatabaseOperation)
	}
	return nil
}

func (s *SessionService) ValidateSession(sessionID uint) (bool, error) {
	if sessionID == 0 {
		return false, ErrInvalidSessionID
	}

	var session model.Session
	if err := s.db.First(&session, sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to validate session: %w", ErrDatabaseOperation)
	}
	return true, nil
}
