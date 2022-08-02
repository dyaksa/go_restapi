package app

import (
	"github.com/julienschmidt/httprouter"
	"golang_restapi/exception"
	"golang_restapi/handler"
)

func SetupRouter(categoryHandler handler.CategoryHandler) *httprouter.Router {
	route := httprouter.New()

	route.GET("/categories", categoryHandler.FindAll)
	route.POST("/categories", categoryHandler.Create)
	route.GET("/categories/:categoryId", categoryHandler.FindOneByID)
	route.PUT("/categories/:categoryId", categoryHandler.Update)
	route.DELETE("/categories/:categoryId", categoryHandler.Delete)

	route.PanicHandler = exception.ErrorHandler
	return route
}
