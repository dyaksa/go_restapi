package config

import (
	"golang_restapi/internal/delivery/http/route"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApp(route *route.ConfigRoute) *echo.Echo {
	var app = echo.New()

	app.Use(middleware.Recover())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	app.Use(middleware.CORS())

	route.Setup(app)

	return app
}
