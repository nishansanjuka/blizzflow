package main_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	. "blizzflow"

	"github.com/joho/godotenv"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

type MockConfig struct {
	SomeConfig string
}

// Mock configuration loader
func LoadMockConfig() *MockConfig {
	return &MockConfig{
		SomeConfig: "test-configuration",
	}
}

var _ = ginkgo.Describe("App Lifecycle", func() {
	var app *App

	ginkgo.BeforeEach(func() {
		app = NewApp()
	})

	ginkgo.Context("App Initialization", func() {
		ginkgo.It("should initialize the app correctly", func() {
			gomega.Expect(app).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Startup", func() {
		ginkgo.It("should assign the context during startup", func() {
			ctx := context.Background()
			app.Startup(ctx)
			gomega.Expect(app.Ctx).To(gomega.Equal(ctx))
		})
	})

	ginkgo.Context("DomReady", func() {
		ginkgo.It("should call DomReady without error", func() {
			ctx := context.Background()
			gomega.Expect(func() { app.DomReady(ctx) }).ShouldNot(gomega.Panic())
		})
	})

	ginkgo.Context("BeforeClose", func() {
		ginkgo.It("should return false by default", func() {
			ctx := context.Background()
			result := app.BeforeClose(ctx)
			gomega.Expect(result).To(gomega.BeFalse())
		})
	})

	ginkgo.Context("Shutdown", func() {
		ginkgo.It("should call Shutdown without error", func() {
			ctx := context.Background()
			gomega.Expect(func() { app.Shutdown(ctx) }).ShouldNot(gomega.Panic())
		})
	})

	ginkgo.Context("Greet Functionality", func() {
		ginkgo.It("should return the correct greeting message", func() {
			name := "World"
			expected := "Hello World, It's show time!"
			gomega.Expect(app.Greet(name)).To(gomega.Equal(expected))
		})
	})
})

var _ = ginkgo.Describe("Main Application Initialization", func() {
	var (
		app *App
	)

	ginkgo.BeforeEach(func() {
		// Create app and mock config
		app = NewApp()
	})

	ginkgo.Context("Wails Application Options", func() {
		var wailsOptions *options.App

		ginkgo.BeforeEach(func() {
			wailsOptions = &options.App{
				Title:             "blizzflow",
				Width:             1024,
				Height:            768,
				DisableResize:     false,
				Fullscreen:        false,
				Frameless:         false,
				StartHidden:       false,
				HideWindowOnClose: false,
				BackgroundColour: &options.RGBA{
					R: 255, G: 255, B: 255, A: 255,
				},
				AssetServer: &assetserver.Options{
					// Mock or provide test assets if needed
				},
				OnStartup:        app.Startup,
				OnDomReady:       app.DomReady,
				OnBeforeClose:    app.BeforeClose,
				OnShutdown:       app.Shutdown,
				WindowStartState: options.Normal,
				Bind:             []interface{}{app},
				LogLevel:         logger.DEBUG,

				// Platform Specific Options
				Windows: &windows.Options{
					WebviewIsTransparent: false,
					WindowIsTranslucent:  false,
					DisableWindowIcon:    false,
					ZoomFactor:           1.0,
				},
				Mac: &mac.Options{
					TitleBar: &mac.TitleBar{
						TitlebarAppearsTransparent: true,
						HideTitle:                  false,
						HideTitleBar:               false,
						FullSizeContent:            false,
						UseToolbar:                 false,
						HideToolbarSeparator:       true,
					},
					Appearance:           mac.NSAppearanceNameDarkAqua,
					WebviewIsTransparent: true,
					WindowIsTranslucent:  true,
				},
			}
		})

		ginkgo.It("should create valid Wails application options", func() {
			gomega.Expect(wailsOptions).NotTo(gomega.BeNil())
			gomega.Expect(wailsOptions.Title).To(gomega.Equal("blizzflow"))
			gomega.Expect(wailsOptions.Width).To(gomega.Equal(1024))
			gomega.Expect(wailsOptions.Height).To(gomega.Equal(768))
		})

		ginkgo.It("should have correct callback methods", func() {
			gomega.Expect(wailsOptions.OnStartup).NotTo(gomega.BeNil())
			gomega.Expect(wailsOptions.OnDomReady).NotTo(gomega.BeNil())
			gomega.Expect(wailsOptions.OnBeforeClose).NotTo(gomega.BeNil())
			gomega.Expect(wailsOptions.OnShutdown).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Application Lifecycle", func() {
		ginkgo.It("should create a new app instance", func() {
			gomega.Expect(app).NotTo(gomega.BeNil())
		})

		ginkgo.It("should handle startup without errors", func() {
			ctx := context.Background()
			gomega.Expect(func() { app.Startup(ctx) }).ShouldNot(gomega.Panic())
		})

		ginkgo.It("should handle DOM ready without errors", func() {
			ctx := context.Background()
			gomega.Expect(func() { app.DomReady(ctx) }).ShouldNot(gomega.Panic())
		})

		ginkgo.It("should handle before close correctly", func() {
			ctx := context.Background()
			result := app.BeforeClose(ctx)
			gomega.Expect(result).To(gomega.BeFalse())
		})

		ginkgo.It("should handle shutdown without errors", func() {
			ctx := context.Background()
			gomega.Expect(func() { app.Shutdown(ctx) }).ShouldNot(gomega.Panic())
		})
	})
})

var _ = ginkgo.Describe("Environment Configuration", func() {
	ginkgo.BeforeEach(func() {
		//attempt to create env
		err := os.WriteFile("./.env.test", []byte("WRITE_ENV=SUCCESS"), os.ModePerm)

		if err != nil {
			fmt.Println(err.Error())
		}
		// Attempt to load .env file, but don't fail if it doesn't exist
		godotenv.Load("./.env.test")
		fmt.Println("env.test loaded successfully")
	})

	ginkgo.Context("Environment Variables", func() {
		ginkgo.It("should set and retrieve environment variables", func() {
			// Set a test environment variable
			os.Setenv("TESTENV", "TESTPASS")

			// Retrieve the environment variable
			testEnv := os.Getenv("TESTENV")

			// Assert the environment variable is set correctly
			gomega.Expect(testEnv).To(gomega.Equal("TESTPASS"))
			gomega.Expect(os.Getenv("WRITE_ENV")).To(gomega.Equal("SUCCESS"))
		})
	})

	ginkgo.Context("Configuration Loading", func() {
		ginkgo.It("should have a valid configuration", func() {
			// Create a mock configuration
			cfg := &MockConfig{
				SomeConfig: "test-configuration",
			}

			// Assert the configuration is not nil
			gomega.Expect(cfg).NotTo(gomega.BeNil())

			// Assert the configuration has the expected value
			gomega.Expect(cfg.SomeConfig).To(gomega.Equal("test-configuration"))
		})
	})
})

func TestAppLifeCycle(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Component Test Suite")

	fileName := "./.env.test"

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("File does not exist:", fileName)
		return
	}

	// Attempt to delete the file
	err := os.Remove(fileName)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}
	fmt.Println("\nenv.test deleted successfully")
}
