package auth_service

import (
	"blizzflow/backend/domain/model"
	repository "blizzflow/backend/domain/repositories"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Custom errors
var (
	ErrEmptyCredentials   = fmt.Errorf("username and password cannot be empty")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrSessionNotFound    = fmt.Errorf("session not found")
	ErrDatabaseOperation  = fmt.Errorf("database operation failed")
	ErrPasswordHash       = fmt.Errorf("password hashing failed")
	ErrSecurityQuestions  = fmt.Errorf("security questions validation failed")
	ErrInvalidAnswers     = fmt.Errorf("incorrect security answers provided")
)

type AuthService struct {
	userRepo              *repository.UserRepository
	sessionRepo           *repository.SessionRepository
	securityQuestionsRepo *repository.SecurityQuestionRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	securityQuestionsRepo *repository.SecurityQuestionRepository,
) *AuthService {
	return &AuthService{
		userRepo:              userRepo,
		sessionRepo:           sessionRepo,
		securityQuestionsRepo: securityQuestionsRepo,
	}
}

func (s *AuthService) Register(username, password string) error {
	if username == "" || password == "" {
		return ErrEmptyCredentials
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", ErrPasswordHash)
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(passwordHash),
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %w", ErrDatabaseOperation)
	}
	return nil
}

func (s *AuthService) Login(username, password string) (*model.Session, error) {
	if username == "" || password == "" {
		return nil, ErrEmptyCredentials
	}

	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", ErrUserNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	session := &model.Session{
		UserID: user.ID,
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", ErrDatabaseOperation)
	}

	return session, nil
}

func (s *AuthService) SetSecurityQuestions(username string, questions map[string]string) error {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", ErrUserNotFound)
	}

	// Delete existing questions if any
	if err := s.userRepo.DB.Where("user_id = ?", user.ID).Delete(&model.SecurityQuestion{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing questions: %w", ErrDatabaseOperation)
	}

	// Store new questions
	for question, answer := range questions {
		hashedAnswer, err := bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash answer: %w", ErrPasswordHash)
		}

		securityQuestion := &model.SecurityQuestion{
			UserID:   user.ID,
			Question: question,
			Answer:   string(hashedAnswer),
		}
		if err := s.securityQuestionsRepo.CreateSecurityQuestion(securityQuestion); err != nil {
			return fmt.Errorf("failed to create security question: %w", ErrDatabaseOperation)
		}
	}
	return nil
}

func (s *AuthService) verifySecurityAnswers(user *model.User, answers map[string]string) error {
	var questions []model.SecurityQuestion
	if err := s.userRepo.DB.Where("user_id = ?", user.ID).Find(&questions).Error; err != nil {
		return fmt.Errorf("failed to get security questions: %w", ErrDatabaseOperation)
	}

	if len(questions) != len(answers) {
		return ErrInvalidAnswers
	}

	for _, q := range questions {
		answer, exists := answers[q.Question]
		if !exists {
			return ErrSecurityQuestions
		}
		if err := bcrypt.CompareHashAndPassword([]byte(q.Answer), []byte(answer)); err != nil {
			return ErrInvalidAnswers
		}
	}
	return nil
}

func (s *AuthService) RecoverPassword(username string, answers map[string]string, newPassword string) error {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", ErrUserNotFound)
	}

	if err := s.verifySecurityAnswers(user, answers); err != nil {
		return err
	}

	// Hash and set new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", ErrPasswordHash)
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update user: %w", ErrDatabaseOperation)
	}
	return nil
}

func (s *AuthService) Logout(sessionID uint) error {
	// Check if session exists
	session, err := s.sessionRepo.GetSessionByID(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", ErrSessionNotFound)
	}
	if session == nil {
		return ErrSessionNotFound
	}

	// Delete the session
	if err := s.sessionRepo.DeleteSession(sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", ErrDatabaseOperation)
	}
	return nil
}
