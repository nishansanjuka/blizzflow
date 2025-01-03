package user_service

import (
	"blizzflow/backend/domain/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(username, password string) (*model.User, error) {
	// Check existing user
	var existingUser model.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *UserService) DeleteUser(userID uint) error {
	return s.db.Delete(&model.User{}, userID).Error
}

// func (s *UserService) ValidatePassword(user *model.User, password string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
// 	return err == nil
// }

// func (s *UserService) ChangePassword(userID uint, newPassword string) error {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	return s.db.Model(&model.User{}).Where("id = ?", userID).
// 		Update("password", string(hashedPassword)).Error
// }

// func (s *UserService) GetUserSessions(userID uint) ([]model.Session, error) {
// 	var sessions []model.Session
// 	if err := s.db.Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
// 		return nil, err
// 	}
// 	return sessions, nil
// }
