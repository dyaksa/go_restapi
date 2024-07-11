package main

import (
	"golang_restapi/app"
	"golang_restapi/db"
	"golang_restapi/handler"
	"golang_restapi/helper"
	"golang_restapi/middleware"
	"golang_restapi/repository"
	"golang_restapi/service"
	"net/http"

	crypt "github.com/dyaksa/encryption-pii/go-encrypt"
	"github.com/go-playground/validator/v10"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	return &http.Server{
		Addr:    "localhost:7001",
		Handler: authMiddleware,
	}
}

func main() {
	sqlDB := db.DB()
	validate := validator.New()
	c, err := crypt.New()
	helper.PanicIf(err)

	categoryRepository := repository.NewCategoryRepository()
	profileRepository := repository.NewProfileRepository()
	textHeapRepository := repository.NewTextHeapRepository()

	categoryServiceImpl := service.NewCategoryService(sqlDB, categoryRepository, validate)
	profileServiceImpl := service.NewProfileService(sqlDB, profileRepository, textHeapRepository, validate, c)

	categoryHandlerImpl := handler.NewCategoryHandler(categoryServiceImpl)
	profileHandlerImpl := handler.NewProfileHandler(profileServiceImpl)

	router := app.SetupRouter(categoryHandlerImpl, profileHandlerImpl)

	authMiddleware := middleware.NewAuthMiddleware(router)
	server := NewServer(authMiddleware)
	err = server.ListenAndServe()
	helper.PanicIf(err)
}
