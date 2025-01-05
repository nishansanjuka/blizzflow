package repository

import (
	"blizzflow/backend/domain/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	if r.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var user model.User
	result := r.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}
