package app

import (
	"golang_restapi/exception"
	"golang_restapi/handler"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(profileHandler handler.ProfileHandler) *httprouter.Router {
	route := httprouter.New()

	route.POST("/profile", profileHandler.Create)
	route.GET("/profile", profileHandler.FindAll)
	route.GET("/profile/:id", profileHandler.FetchProfile)
	route.PUT("/profile/:id", profileHandler.Update)

	route.PanicHandler = exception.ErrorHandler
	return route
}
