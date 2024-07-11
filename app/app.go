package app

import (
	"golang_restapi/exception"
	"golang_restapi/handler"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(categoryHandler handler.CategoryHandler, profileHandler handler.ProfileHandler) *httprouter.Router {
	route := httprouter.New()

	route.GET("/categories", categoryHandler.FindAll)
	route.POST("/categories", categoryHandler.Create)
	route.GET("/categories/:categoryId", categoryHandler.FindOneByID)
	route.PUT("/categories/:categoryId", categoryHandler.Update)
	route.DELETE("/categories/:categoryId", categoryHandler.Delete)

	route.POST("/profile", profileHandler.Create)
	route.GET("/profile", profileHandler.Finding)
	route.GET("/profile/:id", profileHandler.FetchProfile)
	route.PUT("/profile/:id", profileHandler.Update)

	route.PanicHandler = exception.ErrorHandler
	return route
}
