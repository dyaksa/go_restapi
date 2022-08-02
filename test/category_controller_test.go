package test

import (
	"encoding/json"
	validator2 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"golang_restapi/app"
	db2 "golang_restapi/db"
	"golang_restapi/handler"
	"golang_restapi/middleware"
	"golang_restapi/repository"
	"golang_restapi/service"
	"io"
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

	var categoryResponse = map[string]interface{}{}
	body, _ := io.ReadAll(response.Body)

	_ = json.Unmarshal(body, &categoryResponse)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "OK", categoryResponse["status"])
	assert.Equal(t, 200, int(categoryResponse["code"].(float64)))
	assert.Equal(t, "new category", categoryResponse["data"].(map[string]interface{})["name"])
}
