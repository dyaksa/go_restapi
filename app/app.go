package app

import (
	"golang_restapi/exception"
	"golang_restapi/handler"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(categoryHandler handler.CategoryHandler, profileHandler handler.ProfileHandler, profileHandlerv2 handler.ProfileV2, invHandler handler.InvoiceHandlerImpl) *httprouter.Router {
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

	route.POST("/v2/profile", profileHandlerv2.Create)
	route.GET("/v2/profile/:id", profileHandlerv2.FindOneByID)

	route.GET("/invoice/:id", profileHandlerv2.FindOneInv)
	route.GET("/invoices", invHandler.FindAll)
	route.GET("/invoices/search", invHandler.Finding)

	route.PanicHandler = exception.ErrorHandler
	return route
}
