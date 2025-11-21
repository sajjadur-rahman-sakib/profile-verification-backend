package command

import (
	"github.com/labstack/echo/v4"
	"main.go/config"
	"main.go/routes"
)

func Serve() {
	e := echo.New()
	configuration := config.GetConfig()

	config.ConnectDatabase()
	config.UploadDirectory()

	routes.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":" + configuration.GolangPort))
}
