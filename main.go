package main

import (
	"fmt"
	"golang_restapi/app"
	"golang_restapi/db"
	"golang_restapi/handler"
	"golang_restapi/helper"
	"golang_restapi/middleware"
	"golang_restapi/repository"
	"golang_restapi/service"
	"log"
	"net/http"
	"os"

	"github.com/dyaksa/encryption-pii/crypto"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: authMiddleware,
	}
	log.Println("Server running on port", os.Getenv("APP_PORT"))
	return srv
}

func main() {
	if err := godotenv.Load(); err != nil {
		helper.PanicIf(err)
	}

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
