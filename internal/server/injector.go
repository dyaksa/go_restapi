//go:build wireinject
// +build wireinject

package server

import (
	"golang_restapi/internal/config"
	"golang_restapi/internal/delivery/http/controller"
	"golang_restapi/internal/delivery/http/middleware"
	"golang_restapi/internal/delivery/http/route"
	"golang_restapi/internal/repository"
	"golang_restapi/internal/usecase"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

var configSet = wire.NewSet(config.NewCrypto, config.NewValidator, config.NewDatabase)
var repositorySet = wire.NewSet(repository.NewProfileRepository, repository.NewProfileNonPIIRepository)
var usecaseSet = wire.NewSet(usecase.NewProfileUseCase, usecase.NewProfileNonPIIUseCase)
var controllerSet = wire.NewSet(controller.NewProfileController, controller.NewProfileNonPIIController)
var middlewareSet = wire.NewSet(middleware.NewAuthMiddleware)

func InitializeServer() *echo.Echo {
	wire.Build(configSet, repositorySet, usecaseSet, controllerSet, middlewareSet, route.NewRoute, config.NewApp)
	return nil
}
