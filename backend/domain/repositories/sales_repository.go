package repository

import (
	"blizzflow/backend/domain/model"

	"gorm.io/gorm"
)

type SaleRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) Create(sale *model.Sale) error {
	return r.db.Create(sale).Error
}

func (r *SaleRepository) GetByID(id uint) (*model.Sale, error) {
	var sale model.Sale
	err := r.db.First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}
