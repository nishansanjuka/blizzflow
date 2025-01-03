package license_handler

import (
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestLicenseHandlerSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "License Handler Test Suite")
}

const testLicensePath = "test_license.dat"

var (
	handler *LicenseHandler
)

var _ = ginkgo.BeforeSuite(func() {
	os.Remove(testLicensePath)
	handler = NewLicenseHandler(testLicensePath)
})

var _ = ginkgo.AfterSuite(func() {
	os.Remove(testLicensePath)
})

var _ = ginkgo.Describe("License Handler", func() {
	ginkgo.BeforeEach(func() {
		// Clean up before each test
		os.Remove(testLicensePath)
	})

	ginkgo.Context("ReadLicense", func() {
		ginkgo.It("should successfully read a saved license", func() {
			// First save a license
			testLicense := "test-license-key-123"
			err := handler.SaveLicense(testLicense)
			gomega.Expect(err).To(gomega.BeNil())

			// Then try to read it
			readLicense, err := handler.ReadLicense()
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(readLicense).To(gomega.Equal(testLicense))
		})

		ginkgo.It("should handle non-existent file", func() {
			// Try to read without creating file
			_, err := handler.ReadLicense()
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err).To(gomega.BeAssignableToTypeOf(&os.PathError{}))
		})

		ginkgo.It("should handle corrupted file data", func() {
			// Write invalid base64 data
			err := os.WriteFile(testLicensePath, []byte("invalid-base64-data"), 0644)
			gomega.Expect(err).To(gomega.BeNil())

			// Try to read
			_, err = handler.ReadLicense()
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("should handle tampered encrypted data", func() {
			// First save valid license
			err := handler.SaveLicense("test-license")
			gomega.Expect(err).To(gomega.BeNil())

			// Read and modify the file content
			data, err := os.ReadFile(testLicensePath)
			gomega.Expect(err).To(gomega.BeNil())

			// Tamper with the data
			if len(data) > 0 {
				data[len(data)-1] ^= 0xFF // Flip last byte
				err = os.WriteFile(testLicensePath, data, 0644)
				gomega.Expect(err).To(gomega.BeNil())
			}

			// Try to read tampered data
			_, err = handler.ReadLicense()
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
