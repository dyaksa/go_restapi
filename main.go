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

	"github.com/dyaksa/encryption-pii/crypto"
	"github.com/go-playground/validator/v10"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	srv := &http.Server{
		Addr:    ":7001",
		Handler: authMiddleware,
	}
	return srv
}

func main() {
	sqlDB := db.DB()
	validate := validator.New()
	crypto, err := crypto.New(
		crypto.Aes256KeySize,
		crypto.WithInitHeapConnection(),
	)
	helper.PanicIf(err)

	profileRepository := repository.NewProfileRepository()

	profileServiceImpl := service.NewProfileService(sqlDB, profileRepository, validate, crypto)

	profileHandlerImpl := handler.NewProfileHandler(profileServiceImpl)

	router := app.SetupRouter(profileHandlerImpl)

	authMiddleware := middleware.NewAuthMiddleware(router)
	server := NewServer(authMiddleware)
	err = server.ListenAndServe()
	helper.PanicIf(err)
}
