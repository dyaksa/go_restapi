package main

import (
	validator2 "github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	db2 "golang_restapi/db"
	"golang_restapi/handler"
	"golang_restapi/helper"
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

	route := httprouter.New()

	route.GET("/categories", categoryHandler.FindAll)
	route.POST("/categories", categoryHandler.Create)
	route.GET("/categories/:categoryId", categoryHandler.FindOneByID)
	route.PUT("/categories/:categoryId", categoryHandler.Update)
	route.DELETE("/categories/:categoryId", categoryHandler.Delete)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: route,
	}

	err = server.ListenAndServe()
	helper.PanicIf(err)

}
