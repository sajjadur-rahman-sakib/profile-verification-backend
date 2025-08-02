package main

import (
	"main.go/config"
	"main.go/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	config.ConnectDatabase()
	config.UploadDirectory()

	routes.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
