package session_service

import (
	"blizzflow/backend/domain/model"
	"errors"
	"time"

	"gorm.io/gorm"
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
		return nil, err
	}
	return session, nil
}

func (s *SessionService) GetSession(sessionID uint) (*model.Session, error) {
	var session model.Session
	if err := s.db.First(&session, sessionID).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionService) DeleteSession(sessionID uint) error {
	if err := s.db.Delete(&model.Session{}, sessionID).Error; err != nil {
		return err
	}
	return nil
}

func (s *SessionService) ValidateSession(sessionID uint) (bool, error) {
	var session model.Session
	if err := s.db.First(&session, sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
