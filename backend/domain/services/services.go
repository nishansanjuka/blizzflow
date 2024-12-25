package services

import (
	auth_service "blizzflow/backend/domain/services/auth"
	session_service "blizzflow/backend/domain/services/session"
	user_service "blizzflow/backend/domain/services/user"
)

// Export AuthService
type AuthService = auth_service.AuthService

var NewAuthService = auth_service.NewAuthService

// Export SessionService
type SessionService = session_service.SessionService

var NewSessionService = session_service.NewSessionService

// Export UserService
type UserService = user_service.UserService

var NewUserService = user_service.NewUserService
