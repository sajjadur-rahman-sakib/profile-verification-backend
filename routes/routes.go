package routes

import (
	"main.go/handlers"
	"main.go/middleware"
	"main.go/services"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.Static("/uploads", "uploads")

	authService := services.NewAuthService()
	emailService := services.NewEmailService()
	faceService := services.NewFaceService()
	ratingService := services.NewRatingService()
	messageService := services.NewMessageService()
	authHandler := handlers.NewAuthHandler(authService, emailService, faceService, ratingService)
	ratingHandler := handlers.NewRatingHandler(ratingService)
	messageHandler := handlers.NewMessageHandler(messageService)

	api := e.Group("/api")
	{
		api.POST("/user-login", authHandler.Login, middleware.JWTMiddleware)
		api.POST("/update-profile", authHandler.UpdateProfile, middleware.JWTMiddleware)
		api.POST("/change-password", authHandler.ChangePassword, middleware.JWTMiddleware)
		api.POST("/delete-account", authHandler.DeleteAccount, middleware.JWTMiddleware)

		api.POST("/user-signup", authHandler.Signup)
		api.POST("/resend-otp", authHandler.ResendOTP)
		api.POST("/verify-otp", authHandler.VerifyOTP)
		api.POST("/upload-document", authHandler.UploadDocument)
		api.POST("/upload-selfie", authHandler.UploadSelfie)

		api.POST("/forgot-password", authHandler.ForgotPassword)
		api.POST("/reset-password", authHandler.ResetPassword)
		api.POST("/search-profile", authHandler.SearchProfile)

		api.POST("/give-rating", ratingHandler.GiveRating, middleware.JWTMiddleware)
		api.POST("/user-ratings", ratingHandler.GetUserRatings, middleware.JWTMiddleware)

		api.POST("/send-message", messageHandler.SendMessage, middleware.JWTMiddleware)
		api.POST("/user-conversation", messageHandler.GetConversation, middleware.JWTMiddleware)
	}
}
