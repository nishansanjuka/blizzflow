package config_test

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"

	"blizzflow/config"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Config Package", func() {
	var (
		tempConfigFile *os.File
		tempConfigPath string
	)

	ginkgo.BeforeEach(func() {
		var err error
		tempConfigFile, err = os.CreateTemp("", "test-config-*.json")
		gomega.Expect(err).To(gomega.BeNil())
		tempConfigPath = tempConfigFile.Name()
	})

	ginkgo.AfterEach(func() {
		if tempConfigPath != "" {
			os.Remove(tempConfigPath)
		}
	})

	ginkgo.Context("LoadConfig", func() {
		ginkgo.It("should load configuration successfully", func() {
			// Write valid config
			err := os.WriteFile(tempConfigPath, []byte(`{"config_data":"test_value"}`), 0644)
			gomega.Expect(err).To(gomega.BeNil())

			// Mock OpenFile
			originalOpenFile := config.OpenFile
			defer func() { config.OpenFile = originalOpenFile }()

			config.OpenFile = func(name string) (*os.File, error) {
				if name == "config/config.json" {
					return os.Open(tempConfigPath)
				}
				return originalOpenFile(name)
			}

			// Load config
			cfg := config.LoadConfig()

			// Assertions
			gomega.Expect(cfg).NotTo(gomega.BeNil())
			gomega.Expect(cfg.SomeConfig).To(gomega.Equal("test_value"))
		})

		ginkgo.It("should handle file open error", func() {
			// Mock OpenFile to return an error
			originalOpenFile := config.OpenFile
			defer func() { config.OpenFile = originalOpenFile }()

			config.OpenFile = func(name string) (*os.File, error) {
				return nil, os.ErrNotExist
			}

			// Expect a default config
			cfg := config.LoadConfig()
			gomega.Expect(cfg).NotTo(gomega.BeNil())
			gomega.Expect(cfg.SomeConfig).To(gomega.BeEmpty())
		})

		ginkgo.It("should handle JSON decode error", func() {
			// Write invalid JSON
			err := os.WriteFile(tempConfigPath, []byte(`{invalid json}`), 0644)
			gomega.Expect(err).To(gomega.BeNil())

			// Mock OpenFile and NewDecoder
			originalOpenFile := config.OpenFile
			originalNewDecoder := config.NewDecoder
			defer func() {
				config.OpenFile = originalOpenFile
				config.NewDecoder = originalNewDecoder
			}()

			config.OpenFile = func(name string) (*os.File, error) {
				if name == "config/config.json" {
					return os.Open(tempConfigPath)
				}
				return originalOpenFile(name)
			}

			// Simulate decode error
			config.NewDecoder = func(r io.Reader) *json.Decoder {
				decoder := json.NewDecoder(r)
				decoder.DisallowUnknownFields() // This will cause decoding to fail
				return decoder
			}

			// Load config
			cfg := config.LoadConfig()

			// Assertions
			gomega.Expect(cfg).NotTo(gomega.BeNil())
			gomega.Expect(cfg.SomeConfig).To(gomega.BeEmpty())
		})

		ginkgo.It("should handle read error", func() {
			// Create a mock reader that always returns an error
			mockReader := &errorReader{}

			// Mock OpenFile and NewDecoder
			originalOpenFile := config.OpenFile
			originalNewDecoder := config.NewDecoder
			defer func() {
				config.OpenFile = originalOpenFile
				config.NewDecoder = originalNewDecoder
			}()

			config.OpenFile = func(name string) (*os.File, error) {
				if name == "config/config.json" {
					return &os.File{}, nil
				}
				return originalOpenFile(name)
			}

			// Mock decoder to use error reader
			config.NewDecoder = func(r io.Reader) *json.Decoder {
				return json.NewDecoder(mockReader)
			}

			// Load config
			cfg := config.LoadConfig()

			// Assertions
			gomega.Expect(cfg).NotTo(gomega.BeNil())
			gomega.Expect(cfg.SomeConfig).To(gomega.BeEmpty())
		})
	})
})

// errorReader is a mock reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func TestConfig(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Config Suite")
}
