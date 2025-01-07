package services

import (
	"blizzflow/backend/domain/model"
	repository "blizzflow/backend/domain/repositories"
	"blizzflow/backend/infrastructure/database"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var (
	DB            *gorm.DB
	salesRepo     *repository.SaleRepository
	inventoryRepo *repository.InventoryRepository
	salesService  *SalesService
	testDBPath    = "test.db"
)

func TestSalesServiceSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Sales Service Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	os.Remove(testDBPath)
	database.InitDB(testDBPath)
	DB = database.DB

	salesRepo = repository.NewSaleRepository(DB)
	inventoryRepo = repository.NewInventoryRepository(DB)
	salesService = NewSalesService(salesRepo, inventoryRepo)
})

var _ = ginkgo.AfterSuite(func() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	os.Remove(testDBPath)
})

var _ = ginkgo.Describe("Sales Service", func() {
	var testInventory *model.Inventory

	ginkgo.BeforeEach(func() {
		DB.Exec("DELETE FROM sales")
		DB.Exec("DELETE FROM inventories")

		testInventory = &model.Inventory{
			Name:     "Test Product",
			Quantity: 100,
			Price:    10.0,
		}
		DB.Create(testInventory)
	})

	ginkgo.Context("CreateSale", func() {
		ginkgo.It("should create a sale successfully", func() {
			sale, err := salesService.CreateSale(testInventory.ID, 5)

			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(sale).NotTo(gomega.BeNil())
			gomega.Expect(sale.InventoryID).To(gomega.Equal(testInventory.ID))
			gomega.Expect(sale.Quantity).To(gomega.Equal(5))
			gomega.Expect(sale.TotalPrice).To(gomega.Equal(50.0))

			// Verify inventory was updated
			var updatedInventory model.Inventory
			DB.First(&updatedInventory, testInventory.ID)
			gomega.Expect(updatedInventory.Quantity).To(gomega.Equal(95))
		})

		ginkgo.It("should return error for invalid inventory ID", func() {
			sale, err := salesService.CreateSale(0, 5)

			gomega.Expect(err).To(gomega.Equal(ErrInvalidInventoryID))
			gomega.Expect(sale).To(gomega.BeNil())
		})

		ginkgo.It("should return error for invalid quantity", func() {
			sale, err := salesService.CreateSale(testInventory.ID, 0)

			gomega.Expect(err).To(gomega.Equal(ErrInvalidQuantity))
			gomega.Expect(sale).To(gomega.BeNil())
		})

		ginkgo.It("should return error for insufficient stock", func() {
			sale, err := salesService.CreateSale(testInventory.ID, 150)

			gomega.Expect(err).To(gomega.MatchError(gomega.ContainSubstring("insufficient stock")))
			gomega.Expect(sale).To(gomega.BeNil())
		})

		ginkgo.It("should return error for non-existent inventory", func() {
			sale, err := salesService.CreateSale(999, 5)

			gomega.Expect(err).To(gomega.MatchError(gomega.ContainSubstring("inventory not found")))
			gomega.Expect(sale).To(gomega.BeNil())
		})
	})
})
