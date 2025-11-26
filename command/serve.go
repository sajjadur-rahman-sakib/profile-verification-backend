package command

import (
	"verify/config"
	"verify/routes"

	"github.com/labstack/echo/v4"
)

func Serve() {
	e := echo.New()
	configuration := config.GetConfig()

	config.ConnectDatabase()
	config.UploadDirectory()

	routes.SetupRoutes(e)

	e.Logger.Fatal(e.Start(":" + configuration.GolangPort))
}
