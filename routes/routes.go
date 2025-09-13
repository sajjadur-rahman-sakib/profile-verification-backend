package routes

import (
	"main.go/handlers"
	"main.go/services"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.Static("/uploads", "uploads")

	authService := services.NewAuthService()
	emailService := services.NewEmailService()
	faceService := services.NewFaceService()
	authHandler := handlers.NewAuthHandler(authService, emailService, faceService)

	api := e.Group("/api")
	{
		api.POST("/user-signup", authHandler.Signup)
		api.POST("/verify-otp", authHandler.VerifyOTP)
		api.POST("/upload-document", authHandler.UploadDocument)
		api.POST("/upload-selfie", authHandler.UploadSelfie)
		api.POST("/user-login", authHandler.Login)
		api.POST("/delete-account", authHandler.DeleteAccount)
		api.POST("/change-password", authHandler.ChangePassword)
		api.POST("/forgot-password", authHandler.ForgotPassword)
		api.POST("/reset-password", authHandler.ResetPassword)
		api.POST("/update-profile", authHandler.UpdateProfile)
		api.POST("/search-profile", authHandler.SearchProfile)

		api.GET("/search-profile", authHandler.SearchProfile)
	}
}
