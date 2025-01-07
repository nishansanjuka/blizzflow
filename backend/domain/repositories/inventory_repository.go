package repository

import (
	"blizzflow/backend/domain/model"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(inventory *model.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *InventoryRepository) Update(inventory *model.Inventory) error {
	return r.db.Save(inventory).Error
}

func (r *InventoryRepository) GetByID(id uint) (*model.Inventory, error) {
	var inventory model.Inventory
	err := r.db.First(&inventory, id).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}
