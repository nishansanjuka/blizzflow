package services

import (
	repository "blizzflow/backend/domain/repositories"
	"blizzflow/backend/infrastructure/database"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

const testDBPath = "test.db"

var (
	DB               *gorm.DB
	inventoryRepo    *repository.InventoryRepository
	inventoryService *InventoryService
)

func TestInventoryServiceSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Inventory Service Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	os.Remove(testDBPath)
	database.InitDB(testDBPath)
	DB = database.DB
	inventoryRepo = repository.NewInventoryRepository(DB)
	inventoryService = NewInventoryService(inventoryRepo)
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

var _ = ginkgo.Describe("Inventory Service", func() {
	ginkgo.BeforeEach(func() {
		DB.Exec("DELETE FROM inventories")
	})

	ginkgo.It("should create inventory successfully", func() {
		inventory, err := inventoryService.CreateInventory("Test Item", 10, 99.99)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(inventory.Name).To(gomega.Equal("Test Item"))
		gomega.Expect(inventory.Quantity).To(gomega.Equal(10))
	})

	ginkgo.It("should fail with invalid inputs", func() {
		_, err := inventoryService.CreateInventory("", 10, 99.99)
		gomega.Expect(err).To(gomega.Equal(ErrInvalidInventoryName))

		_, err = inventoryService.CreateInventory("Test", -1, 99.99)
		gomega.Expect(err).To(gomega.Equal(ErrInvalidQuantity))
	})
})
