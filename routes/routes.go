package routes

import (
	"verify/handlers"
	"verify/middleware"
	"verify/services"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.Static("/uploads", "uploads")

	// Initialize services
	emailService := services.NewEmailService()
	faceService := services.NewFaceService()
	ratingService := services.NewRatingService()
	messageService := services.NewMessageService()
	loginService := services.NewLoginService()
	signupService := services.NewSignupService()
	uploadService := services.NewUploadService()
	profileService := services.NewProfileService(uploadService)
	passwordService := services.NewPasswordService()
	otpService := services.NewOTPService()
	accountService := services.NewAccountService()
	tokenService := services.NewTokenService()

	// Initialize handlers
	loginHandler := handlers.NewLoginHandler(loginService, tokenService)
	signupHandler := handlers.NewSignupHandler(signupService, otpService, emailService, uploadService, tokenService)
	profileHandler := handlers.NewProfileHandler(profileService)
	passwordHandler := handlers.NewPasswordHandler(passwordService, otpService, emailService)
	otpHandler := handlers.NewOTPHandler(otpService, emailService)
	accountHandler := handlers.NewAccountHandler(accountService)
	uploadHandler := handlers.NewUploadHandler(uploadService, profileService, faceService)
	ratingHandler := handlers.NewRatingHandler(ratingService)
	messageHandler := handlers.NewMessageHandler(messageService)

	api := e.Group("/api")
	{
		// Login routes
		api.POST("/user-login", loginHandler.Login)

		// Signup routes
		api.POST("/user-signup", signupHandler.Signup)

		// Profile routes
		api.POST("/update-profile", profileHandler.UpdateProfile, middleware.JWTMiddleware)
		api.POST("/search-profile", profileHandler.SearchProfile, middleware.JWTMiddleware)

		// Password routes
		api.POST("/change-password", passwordHandler.ChangePassword, middleware.JWTMiddleware)
		api.POST("/forgot-password", passwordHandler.ForgotPassword)
		api.POST("/reset-password", passwordHandler.ResetPassword)

		// OTP routes
		api.POST("/resend-otp", otpHandler.ResendOTP, middleware.JWTMiddleware)
		api.POST("/verify-otp", otpHandler.VerifyOTP, middleware.JWTMiddleware)

		// Account routes
		api.POST("/delete-account", accountHandler.DeleteAccount, middleware.JWTMiddleware)

		// Upload routes
		api.POST("/upload-document", uploadHandler.UploadDocument, middleware.JWTMiddleware)
		api.POST("/upload-selfie", uploadHandler.UploadSelfie, middleware.JWTMiddleware)

		// Rating routes
		api.POST("/give-rating", ratingHandler.GiveRating, middleware.JWTMiddleware)
		api.POST("/user-ratings", ratingHandler.GetUserRatings, middleware.JWTMiddleware)

		// Message routes
		api.POST("/user-contacts", messageHandler.GetContacts, middleware.JWTMiddleware)
		api.POST("/send-message", messageHandler.SendMessage, middleware.JWTMiddleware)
		api.POST("/user-conversation", messageHandler.GetConversation, middleware.JWTMiddleware)
	}
}
