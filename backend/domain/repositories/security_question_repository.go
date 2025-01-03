package repository

import (
	"blizzflow/backend/domain/model"
	"errors"

	"gorm.io/gorm"
)

type SecurityQuestionRepository struct {
	db *gorm.DB
}

func NewSecurityQuestionRepository(db *gorm.DB) *SecurityQuestionRepository {
	return &SecurityQuestionRepository{db: db}
}

func (r *SecurityQuestionRepository) CreateSecurityQuestion(question *model.SecurityQuestion) error {
	result := r.db.Create(question)
	return result.Error
}

func (r *SecurityQuestionRepository) GetSecurityQuestion(id uint) (*model.SecurityQuestion, error) {
	var question model.SecurityQuestion
	result := r.db.First(&question, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &question, result.Error
}

func (r *SecurityQuestionRepository) GetSecurityQuestionsByUserID(userID uint) ([]model.SecurityQuestion, error) {
	var questions []model.SecurityQuestion
	result := r.db.Where("user_id = ?", userID).Find(&questions)
	return questions, result.Error
}

func (r *SecurityQuestionRepository) UpdateSecurityQuestion(question *model.SecurityQuestion) error {
	result := r.db.Save(question)
	return result.Error
}

func (r *SecurityQuestionRepository) DeleteSecurityQuestion(id uint) error {
	result := r.db.Delete(&model.SecurityQuestion{}, id)
	return result.Error
}

func (r *SecurityQuestionRepository) DeleteUserSecurityQuestions(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&model.SecurityQuestion{})
	return result.Error
}

func (r *SecurityQuestionRepository) ValidateSecurityQuestion(userID uint, question string, answer string) (bool, error) {
	var securityQuestion model.SecurityQuestion
	result := r.db.Where("user_id = ? AND question = ?", userID, question).First(&securityQuestion)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return securityQuestion.Answer == answer, nil
}
