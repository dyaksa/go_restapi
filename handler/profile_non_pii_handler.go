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

type ProfileNonPIIHandler interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FetchProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type ProfileNonPIIHandlerImpl struct {
	profileNonPIIService service.ProfileNonPIIService
}

func NewProfileNonPIIHandler(profileNonPii service.ProfileNonPIIService) *ProfileNonPIIHandlerImpl {
	return &ProfileNonPIIHandlerImpl{
		profileNonPIIService: profileNonPii,
	}
}

func (h *ProfileNonPIIHandlerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var profileCreateRequest web.ProfileRequest
	helper.JSONDecoder(r, &profileCreateRequest)

	data, err := h.profileNonPIIService.Create(r.Context(), profileCreateRequest)
	helper.PanicIf(err)

	webResponse := web.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileNonPIIHandlerImpl) FetchProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := uuid.Parse(params.ByName("id"))
	helper.PanicIf(err)

	data, err := h.profileNonPIIService.FetchProfile(r.Context(), id)
	helper.PanicIf(err)

	webResponse := web.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileNonPIIHandlerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	pagination := utils.Pagination{Page: page, Limit: limit}
	paramsListProfile := dto.ParamsListProfile{Key: key, Value: value}

	data, err := h.profileNonPIIService.FindAll(r.Context(), pagination, paramsListProfile)
	helper.PanicIf(err)

	webResponse := web.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
		Meta: web.Meta{
			Page: pagination.GetPage(),
			Size: pagination.GetLimit(),
		},
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileNonPIIHandlerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var profileUpdateRequest web.ProfileRequest
	helper.JSONDecoder(r, &profileUpdateRequest)

	id, err := uuid.Parse(params.ByName("id"))
	helper.PanicIf(err)

	data, err := h.profileNonPIIService.Update(r.Context(), id, profileUpdateRequest)
	helper.PanicIf(err)

	webResponse := web.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}
