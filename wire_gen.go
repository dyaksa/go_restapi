// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"golang_restapi/app"
	"golang_restapi/db"
	"golang_restapi/handler"
	"golang_restapi/middleware"
	"golang_restapi/repository"
	"golang_restapi/service"
	"net/http"
)

// Injectors from injector.go:

func InitializedServer() *http.Server {
	sqlDB := db.DB()
	categoryRepository := repository.NewCategoryRepository()
	validate := validator.New()
	categoryServiceImpl := service.NewCategoryService(sqlDB, categoryRepository, validate)
	categoryHandlerImpl := handler.NewCategoryHandler(categoryServiceImpl)
	router := app.SetupRouter(categoryHandlerImpl)
	authMiddleware := middleware.NewAuthMiddleware(router)
	server := NewServer(authMiddleware)
	return server
}

// injector.go:

var categorySet = wire.NewSet(repository.NewCategoryRepository, wire.Bind(new(repository.Category), new(*repository.CategoryRepository)), service.NewCategoryService, wire.Bind(new(service.CategoryService), new(*service.CategoryServiceImpl)), handler.NewCategoryHandler, wire.Bind(new(handler.CategoryHandler), new(*handler.CategoryHandlerImpl)))