package auth_service

import (
	"blizzflow/backend/domain/model"
	repository "blizzflow/backend/domain/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
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
		return errors.New("username and password cannot be empty")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(passwordHash),
	}

	return s.userRepo.CreateUser(user)
}

func (s *AuthService) Login(username, password string) (*model.Session, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	session := &model.Session{
		UserID: user.ID,
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) SetSecurityQuestions(username string, questions map[string]string) error {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}

	// Delete existing questions if any
	if err := s.userRepo.DB.Where("user_id = ?", user.ID).Delete(&model.SecurityQuestion{}).Error; err != nil {
		return err
	}

	// Store new questions
	for question, answer := range questions {
		hashedAnswer, err := bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		securityQuestion := &model.SecurityQuestion{
			UserID:   user.ID,
			Question: question,
			Answer:   string(hashedAnswer),
		}
		err = s.securityQuestionsRepo.CreateSecurityQuestion(securityQuestion)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) verifySecurityAnswers(user *model.User, answers map[string]string) error {
	var questions []model.SecurityQuestion
	if err := s.userRepo.DB.Where("user_id = ?", user.ID).Find(&questions).Error; err != nil {
		return err
	}

	if len(questions) != len(answers) {
		return errors.New("incorrect number of answers provided")
	}

	for _, q := range questions {
		answer, exists := answers[q.Question]
		if !exists {
			return errors.New("missing answer for question")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(q.Answer), []byte(answer)); err != nil {
			return errors.New("incorrect answer provided")
		}
	}
	return nil
}

func (s *AuthService) RecoverPassword(username string, answers map[string]string, newPassword string) error {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}

	if err := s.verifySecurityAnswers(user, answers); err != nil {
		return err
	}

	// Hash and set new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.UpdateUser(user)
}
