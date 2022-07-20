package handler

import (
	"github.com/julienschmidt/httprouter"
	"golang_restapi/helper"
	"golang_restapi/model/web"
	"golang_restapi/service"
	"net/http"
	"strconv"
)

type CategoryHandler interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindOneByID(w http.ResponseWriter, r *http.Request, param httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, param httprouter.Params)
}

type CategoryHandlerImpl struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &CategoryHandlerImpl{categoryService: categoryService}
}

func (h *CategoryHandlerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var categoryCreateRequest web.CategoryCreateRequest
	helper.JSONDecoder(r, &categoryCreateRequest)

	data := h.categoryService.Create(r.Context(), categoryCreateRequest)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *CategoryHandlerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("categoryId")
	idInt, err := strconv.Atoi(id)
	helper.PanicIf(err)

	categoryUpdateRequest := web.CategoryUpdateRequest{
		ID: idInt,
	}
	helper.JSONDecoder(r, &categoryUpdateRequest)

	data := h.categoryService.Update(r.Context(), categoryUpdateRequest)
	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *CategoryHandlerImpl) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("categoryId")
	idInt, err := strconv.Atoi(id)
	helper.PanicIf(err)

	h.categoryService.Delete(r.Context(), idInt)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *CategoryHandlerImpl) FindOneByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("categoryId")
	idInt, err := strconv.Atoi(id)
	helper.PanicIf(err)

	data := h.categoryService.FindByID(r.Context(), idInt)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *CategoryHandlerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data := h.categoryService.FindAll(r.Context())

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}
