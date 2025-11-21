package config

import (
	"os"

	"github.com/joho/godotenv"
)

var configuration Config

type Config struct {
	DatabaseHost     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	DatabasePort     string
	GolangPort       string
	PythonPort       string
	JwtSecret        string
	SmtpHost         string
	SmtpPort         string
	SmtpUsername     string
	SmtpPassword     string
	SmtpEmail        string
	FaceService      string
	UploadDirectory  string
}

func loadConfig() {
	err := godotenv.Load("app.env")
	if err != nil {
		panic("Error loading app.env file")
	}

	configuration = Config{
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabaseUsername: os.Getenv("DATABASE_USERNAME"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		GolangPort:       os.Getenv("GOLANG_PORT"),
		PythonPort:       os.Getenv("PYTHON_PORT"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		SmtpHost:         os.Getenv("SMTP_HOST"),
		SmtpPort:         os.Getenv("SMTP_PORT"),
		SmtpUsername:     os.Getenv("SMTP_USERNAME"),
		SmtpPassword:     os.Getenv("SMTP_PASSWORD"),
		SmtpEmail:        os.Getenv("SMTP_EMAIL"),
		FaceService:      os.Getenv("FACE_SERVICE"),
		UploadDirectory:  os.Getenv("UPLOAD_DIRECTORY"),
	}
}

func GetConfig() Config {
	loadConfig()
	return configuration
}
