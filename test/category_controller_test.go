package test

import (
	validator2 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"golang_restapi/app"
	db2 "golang_restapi/db"
	"golang_restapi/handler"
	"golang_restapi/middleware"
	"golang_restapi/repository"
	"golang_restapi/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRouter() http.Handler {
	db := db2.DB()
	validator := validator2.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(db, categoryRepository, validator)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	router := app.SetupRouter(categoryHandler)
	return middleware.NewAuthMiddleware(router)
}

func TestCreateCategorySuccess(t *testing.T) {
	router := setupRouter()

	requestBody := strings.NewReader(`{"name": "new category"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/categories", requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-Api-Key", "secret")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()

	assert.Equal(t, 200, response.StatusCode)
}
