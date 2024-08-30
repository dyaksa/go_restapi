package route

import (
	"golang_restapi/internal/delivery/http/controller"
	"golang_restapi/internal/delivery/http/middleware"

	"github.com/labstack/echo/v4"
)

type ConfigRoute struct {
	AuthMiddleware          *middleware.AuthMiddleware
	ProfileController       *controller.ProfileController
	ProfileNonPIIController *controller.ProfileNonPIIController
}

func NewRoute(
	authMiddleware *middleware.AuthMiddleware,
	profileController *controller.ProfileController,
	profileNonPIIController *controller.ProfileNonPIIController,
) *ConfigRoute {
	return &ConfigRoute{
		AuthMiddleware:          authMiddleware,
		ProfileController:       profileController,
		ProfileNonPIIController: profileNonPIIController,
	}
}

func (c *ConfigRoute) Setup(app *echo.Echo) {
	c.profileApiRoute(app)
	c.profileNonPIIApiRoute(app)
}

func (c *ConfigRoute) profileApiRoute(app *echo.Echo) {
	api := app.Group("/perf-test-go")
	{
		api.POST("/profile", echo.HandlerFunc(c.ProfileController.Create))
		api.PUT("/profile/:id", echo.HandlerFunc(c.ProfileController.Update))
		api.GET("/profile/:id", echo.HandlerFunc(c.ProfileController.FetchProfile))
		api.GET("/profile", echo.HandlerFunc(c.ProfileController.FindAll))
	}
}

func (c *ConfigRoute) profileNonPIIApiRoute(app *echo.Echo) {
	api := app.Group("/perf-test-go")
	{
		api.POST("/profile-not-pii", echo.HandlerFunc(c.ProfileNonPIIController.Create))
		api.GET("/profile-not-pii/:id", echo.HandlerFunc(c.ProfileNonPIIController.FetchProfile))
		api.GET("/profile-not-pii", echo.HandlerFunc(c.ProfileNonPIIController.FindAll))
		api.PUT("/profile-not-pii/:id", echo.HandlerFunc(c.ProfileNonPIIController.Update))
	}
}
