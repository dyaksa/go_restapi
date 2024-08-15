package app

import (
	"golang_restapi/exception"
	"golang_restapi/handler"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(profileHandler handler.ProfileHandler) *httprouter.Router {
	route := httprouter.New()

	route.POST("/perf-test/go", profileHandler.Create)
	route.GET("/perf-test/go", profileHandler.FindAll)
	route.GET("/perf-test/go/:id", profileHandler.FetchProfile)
	route.PUT("/perf-test/go/:id", profileHandler.Update)

	route.PanicHandler = exception.ErrorHandler
	return route
}
