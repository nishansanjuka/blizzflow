package services

import (
	"blizzflow/backend/domain/model"
	repository "blizzflow/backend/domain/repositories"
	"fmt"
)

// Custom errors
var (
	ErrInvalidInventoryName = fmt.Errorf("invalid inventory name")
	ErrInvalidQuantity      = fmt.Errorf("invalid quantity")
	ErrInvalidPrice         = fmt.Errorf("invalid price")
	ErrInventoryNotFound    = fmt.Errorf("inventory not found")
	ErrDatabaseOperation    = fmt.Errorf("database operation failed")
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{inventoryRepo: repo}
}

func (s *InventoryService) CreateInventory(name string, quantity int, price float64) (*model.Inventory, error) {
	if name == "" {
		return nil, ErrInvalidInventoryName
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if price < 0 {
		return nil, ErrInvalidPrice
	}

	inventory := &model.Inventory{
		Name:     name,
		Quantity: quantity,
		Price:    price,
	}
	if err := s.inventoryRepo.Create(inventory); err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", ErrDatabaseOperation)
	}
	return inventory, nil
}
