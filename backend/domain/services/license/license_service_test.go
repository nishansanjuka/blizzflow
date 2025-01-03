package services

import (
	repository "blizzflow/backend/domain/repositories"
	"blizzflow/backend/infrastructure/database"
	"os"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestLicenseServiceSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "License Service Test Suite")
}

const testDBPath = "license_test.db"

var (
	DB             *gorm.DB
	licenseRepo    *repository.LicenseRepository
	licenseService *LicenseService
)

var _ = ginkgo.BeforeSuite(func() {
	os.Remove(testDBPath)
	database.InitDB(testDBPath)
	DB = database.DB

	licenseRepo = repository.NewLicenseRepository(DB)
	licenseService = NewLicenseService(licenseRepo)
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

var _ = ginkgo.Describe("License Service", func() {
	ginkgo.BeforeEach(func() {
		DB.Exec("DELETE FROM licenses")
	})

	ginkgo.It("should generate license successfully", func() {
		expiresAt := time.Now().Add(24 * time.Hour)
		license, err := licenseService.GenerateLicense("testuser", expiresAt)

		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(license).ToNot(gomega.BeNil())
		gomega.Expect(license.Username).To(gomega.Equal("testuser"))
		gomega.Expect(license.ExpiresAt.Unix()).To(gomega.Equal(expiresAt.Unix()))
	})

	ginkgo.It("should validate valid license", func() {
		expiresAt := time.Now().Add(24 * time.Hour)
		license, err := licenseService.GenerateLicense("testuser", expiresAt)
		gomega.Expect(err).To(gomega.BeNil())

		valid, err := licenseService.ValidateLicense(license.Key)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(valid).To(gomega.BeTrue())
	})

	ginkgo.It("should fail validation for expired license", func() {
		expiresAt := time.Now().Add(-24 * time.Hour) // expired
		license, err := licenseService.GenerateLicense("testuser", expiresAt)
		gomega.Expect(err).To(gomega.BeNil())

		_, err = licenseService.DecodeLicense(license.Key)
		gomega.Expect(err).To(gomega.Equal(ErrExpiredLicense))
	})

	ginkgo.It("should fail validation for invalid license format", func() {
		valid, err := licenseService.ValidateLicense("")
		gomega.Expect(err).To(gomega.Equal(ErrInvalidLicenseKey))
		gomega.Expect(valid).To(gomega.BeFalse())
	})

	ginkgo.It("should fail for empty username", func() {
		expiresAt := time.Now().Add(24 * time.Hour)
		_, err := licenseService.GenerateLicense("", expiresAt)
		gomega.Expect(err).To(gomega.Equal(ErrInvalidUsername))
	})
})
