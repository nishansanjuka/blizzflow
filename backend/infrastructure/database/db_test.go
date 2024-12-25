package database

import (
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestDatabaseSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Database Test Suite")
}

const testDBPath = "test.db"

var _ = ginkgo.Describe("Database Connection", func() {
	ginkgo.BeforeEach(func() {
		DB = nil
		os.Remove(testDBPath)
	})

	ginkgo.AfterEach(func() {
		if DB != nil {
			sqlDB, err := DB.DB()
			if err == nil {
				sqlDB.Close()
			}
			DB = nil
		}
	})

	ginkgo.It("should initialize database successfully", func() {
		err := InitDB(testDBPath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(DB).NotTo(gomega.BeNil())
	})

	ginkgo.It("should handle database initialization errors", func() {
		err := InitDB(testDBPath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		err = InitDB(testDBPath)
		gomega.Expect(err).To(gomega.MatchError("database already initialized"))
	})

	ginkgo.It("should handle database closure errors", func() {
		// First create a valid DB
		err := InitDB(testDBPath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Close it properly
		err = CloseDB()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Try closing again
		err = CloseDB()
		gomega.Expect(err).NotTo(gomega.HaveOccurred()) // Should handle nil DB gracefully
	})
})

var _ = ginkgo.AfterSuite(func() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
		DB = nil
	}
	os.Remove(testDBPath)
})
