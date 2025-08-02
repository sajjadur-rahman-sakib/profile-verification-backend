package routes

import (
	"main.go/handlers"
	"main.go/services"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	authService := services.NewAuthService()
	emailService := services.NewEmailService()
	faceService := services.NewFaceService()
	authHandler := handlers.NewAuthHandler(authService, emailService, faceService)

	api := e.Group("/api")
	{
		api.POST("/signup", authHandler.Signup)
		api.POST("/verify-otp", authHandler.VerifyOTP)
		api.POST("/upload-document", authHandler.UploadDocument)
		api.POST("/upload-selfie", authHandler.UploadSelfie)
		api.POST("/login", authHandler.Login)
	}
}
