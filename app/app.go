package app

import (
	"golang_restapi/exception"
	"golang_restapi/handler"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(profileHandler handler.ProfileHandler, profileNonPIIHandler handler.ProfileNonPIIHandler) *httprouter.Router {
	route := httprouter.New()

	route.POST("/perf-test-go/profile", profileHandler.Create)
	route.GET("/perf-test-go/profile", profileHandler.FindAll)
	route.GET("/perf-test-go/profile/:id", profileHandler.FetchProfile)
	route.PUT("/perf-test-go/profile/:id", profileHandler.Update)

	route.POST("/perf-test-go/profile-not-pii", profileNonPIIHandler.Create)
	route.GET("/perf-test-go/profile-not-pii", profileNonPIIHandler.FindAll)
	route.GET("/perf-test-go/profile-not-pii/:id", profileNonPIIHandler.FetchProfile)
	route.PUT("/perf-test-go/profile-not-pii/:id", profileNonPIIHandler.Update)

	route.PanicHandler = exception.ErrorHandler
	return route
}
