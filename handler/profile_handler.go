package handler

import (
	"golang_restapi/dto"
	"golang_restapi/helper"
	"golang_restapi/model/web"
	"golang_restapi/service"
	"golang_restapi/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ProfileHandler interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FetchProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type ProfileHandlerImpl struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandlerImpl {
	return &ProfileHandlerImpl{profileService: profileService}
}

func (h *ProfileHandlerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var profileCreateRequest web.ProfileRequest
	helper.JSONDecoder(r, &profileCreateRequest)

	data, err := h.profileService.Create(r.Context(), profileCreateRequest)
	helper.PanicIf(err)

	webResponse := web.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileHandlerImpl) FetchProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := uuid.Parse(params.ByName("id"))
	helper.PanicIf(err)

	data, err := h.profileService.FetchProfile(r.Context(), id)
	helper.PanicIf(err)

	webResponse := web.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileHandlerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	pagination := utils.Pagination{Page: page, Limit: limit}
	paramsListProfile := dto.ParamsListProfile{Key: key, Value: value}
	data, err := h.profileService.FindAll(r.Context(), pagination, paramsListProfile)
	helper.PanicIf(err)

	webResponse := web.Response{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   data,
		Meta: web.Meta{
			Page: pagination.GetPage(),
			Size: pagination.GetLimit(),
		},
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileHandlerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := uuid.Parse(params.ByName("id"))
	helper.PanicIf(err)

	var profileUpdateRequest web.ProfileRequest
	helper.JSONDecoder(r, &profileUpdateRequest)

	data, err := h.profileService.Update(r.Context(), id, profileUpdateRequest)
	helper.PanicIf(err)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}
