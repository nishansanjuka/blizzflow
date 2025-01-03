package repository

import (
	"blizzflow/backend/domain/model"

	"gorm.io/gorm"
)

type LicenseRepository struct {
	db *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) *LicenseRepository {
	return &LicenseRepository{db: db}
}

func (r *LicenseRepository) Create(license *model.License) error {
	return r.db.Create(license).Error
}

func (r *LicenseRepository) GetByKey(key string) (*model.License, error) {
	var license model.License
	err := r.db.Where("key = ?", key).First(&license).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}
