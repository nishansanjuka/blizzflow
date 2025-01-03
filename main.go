package main

import (
	license_handler "blizzflow/backend/domain/handlers/license"
	repository "blizzflow/backend/domain/repositories"
	auth_service "blizzflow/backend/domain/services/auth"
	license_service "blizzflow/backend/domain/services/license"
	session_service "blizzflow/backend/domain/services/session"
	user_service "blizzflow/backend/domain/services/user"
	"embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gorm.io/gorm"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// Initialize repositories
	var db *gorm.DB // Initialize your database connection here
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	securityQuestionsRepo := repository.NewSecurityQuestionRepository(db)

	// Initialize services
	userService := user_service.NewUserService(db)
	sessionService := session_service.NewSessionService(db)
	authService := auth_service.NewAuthService(userRepo, sessionRepo, securityQuestionsRepo)
	licenseService := license_service.NewLicenseService(repository.NewLicenseRepository(db))

	// Initialize license handler
	appDir, _ := os.UserConfigDir()
	licensePath := filepath.Join(appDir, "blizzflow", "license.blizz")
	os.MkdirAll(filepath.Dir(licensePath), 0755)
	licenseHandler := license_handler.NewLicenseHandler(licensePath)

	app := application.New(application.Options{
		Name:        "blizzflow",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(userService),
			application.NewService(sessionService),
			application.NewService(authService),
			application.NewService(licenseService),
			application.NewService(licenseHandler),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Create a goroutine that emits an event containing the current time every second.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
