package handler

import (
	"errors"
	"golang_restapi/helper"
	"golang_restapi/model/web"
	"golang_restapi/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ProfileHandler interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FetchProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Finding(w http.ResponseWriter, r *http.Request, params httprouter.Params)
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

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		helper.PanicIf(errors.New("X-Tenant-ID is required"))
	}

	parseTenantID, err := uuid.Parse(tenantID)
	helper.PanicIf(err)

	data := h.profileService.Create(r.Context(), parseTenantID, profileCreateRequest)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileHandlerImpl) FetchProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := uuid.Parse(params.ByName("id"))
	helper.PanicIf(err)

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		helper.PanicIf(errors.New("X-Tenant-ID is required"))
	}

	parseTenantID, err := uuid.Parse(tenantID)
	helper.PanicIf(err)

	data := h.profileService.FetchProfile(r.Context(), parseTenantID, id)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileHandlerImpl) Finding(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	queryValues := r.URL.Query()

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		helper.PanicIf(errors.New("X-Tenant-ID is required"))
	}

	parseTenantID, err := uuid.Parse(tenantID)
	helper.PanicIf(err)

	data, err := h.profileService.FindByTextHeapContent(r.Context(), parseTenantID, queryValues)
	helper.PanicIf(err)

	webResponse := web.ResponseCategory{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileHandlerImpl) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := uuid.Parse(params.ByName("id"))
	helper.PanicIf(err)

	var profileUpdateRequest web.ProfileRequest
	helper.JSONDecoder(r, &profileUpdateRequest)

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		helper.PanicIf(errors.New("X-Tenant-ID is required"))
	}

	parseTenantID, err := uuid.Parse(tenantID)
	helper.PanicIf(err)

	data := h.profileService.Update(r.Context(), id, parseTenantID, profileUpdateRequest)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}
