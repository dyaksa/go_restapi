package main

import (
	validator2 "github.com/go-playground/validator/v10"
	"golang_restapi/app"
	db2 "golang_restapi/db"
	"golang_restapi/handler"
	"golang_restapi/helper"
	"golang_restapi/middleware"
	"golang_restapi/repository"
	"golang_restapi/service"
	"net/http"
)

func main() {

	var (
		db        = db2.DB()
		validator = validator2.New()
	)

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(db, categoryRepository, validator)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	route := app.SetupRouter(categoryHandler)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: middleware.NewAuthMiddleware(route),
	}

	err = server.ListenAndServe()
	helper.PanicIf(err)

}
