package config

import "github.com/labstack/echo"

func NewEcho() *echo.Echo {
	return echo.New()
}
