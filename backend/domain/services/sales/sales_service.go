package services

import (
	"blizzflow/backend/domain/model"
	repository "blizzflow/backend/domain/repositories"
	"fmt"
)

// Custom errors
var (
	ErrInvalidInventoryID = fmt.Errorf("invalid inventory ID")
	ErrInvalidQuantity    = fmt.Errorf("invalid quantity")
	ErrInsufficientStock  = fmt.Errorf("insufficient stock")
	ErrInventoryNotFound  = fmt.Errorf("inventory not found")
	ErrDatabaseOperation  = fmt.Errorf("database operation failed")
)

type SalesService struct {
	saleRepo      *repository.SaleRepository
	inventoryRepo *repository.InventoryRepository
}

func NewSalesService(saleRepo *repository.SaleRepository, invRepo *repository.InventoryRepository) *SalesService {
	return &SalesService{
		saleRepo:      saleRepo,
		inventoryRepo: invRepo,
	}
}

func (s *SalesService) CreateSale(inventoryID uint, quantity int) (*model.Sale, error) {
	if inventoryID == 0 {
		return nil, ErrInvalidInventoryID
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	inventory, err := s.inventoryRepo.GetByID(inventoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inventory: %w", ErrInventoryNotFound)
	}

	if inventory.Quantity < quantity {
		return nil, fmt.Errorf("requested quantity %d exceeds available stock %d: %w",
			quantity, inventory.Quantity, ErrInsufficientStock)
	}

	sale := &model.Sale{
		InventoryID: inventoryID,
		Quantity:    quantity,
		TotalPrice:  float64(quantity) * inventory.Price,
	}

	if err := s.saleRepo.Create(sale); err != nil {
		return nil, fmt.Errorf("failed to create sale record: %w", ErrDatabaseOperation)
	}

	// Update inventory quantity
	inventory.Quantity -= quantity
	if err := s.inventoryRepo.Update(inventory); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", ErrDatabaseOperation)
	}

	return sale, nil
}
