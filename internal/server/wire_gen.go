// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"golang_restapi/internal/config"
	"golang_restapi/internal/delivery/http/controller"
	"golang_restapi/internal/delivery/http/middleware"
	"golang_restapi/internal/delivery/http/route"
	"golang_restapi/internal/repository"
	"golang_restapi/internal/usecase"
)

// Injectors from injector.go:

func InitializeServer() *echo.Echo {
	authMiddleware := middleware.NewAuthMiddleware()
	db := config.NewDatabase()
	crypto := config.NewCrypto()
	validate := config.NewValidator()
	profileRepository := repository.NewProfileRepository(crypto)
	profileUseCase := usecase.NewProfileUseCase(db, crypto, validate, profileRepository)
	profileController := controller.NewProfileController(profileUseCase)
	profileNonPIIRepository := repository.NewProfileNonPIIRepository()
	profileNonPIIUseCase := usecase.NewProfileNonPIIUseCase(db, profileNonPIIRepository, validate)
	profileNonPIIController := controller.NewProfileNonPIIController(profileNonPIIUseCase)
	configRoute := route.NewRoute(authMiddleware, profileController, profileNonPIIController)
	echoEcho := config.NewApp(configRoute)
	return echoEcho
}

// injector.go:

var configSet = wire.NewSet(config.NewCrypto, config.NewValidator, config.NewDatabase)

var repositorySet = wire.NewSet(repository.NewProfileRepository, repository.NewProfileNonPIIRepository)

var usecaseSet = wire.NewSet(usecase.NewProfileUseCase, usecase.NewProfileNonPIIUseCase)

var controllerSet = wire.NewSet(controller.NewProfileController, controller.NewProfileNonPIIController)

var middlewareSet = wire.NewSet(middleware.NewAuthMiddleware)
