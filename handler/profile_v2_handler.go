package handler

import (
	"golang_restapi/helper"
	"golang_restapi/model/web"
	"golang_restapi/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ProfileV2 interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindOneByID(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindOneInv(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type ProfileV2Impl struct {
	profileServicev2 service.ProfileServiceV2
}

func NewProfileV2(profileService service.ProfileServiceV2) *ProfileV2Impl {
	return &ProfileV2Impl{profileServicev2: profileService}
}

func (h *ProfileV2Impl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var profileCreateRequest web.ProfileRequest
	helper.JSONDecoder(r, &profileCreateRequest)

	data := h.profileServicev2.Create(r.Context(), profileCreateRequest)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileV2Impl) FindOneByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	uid, err := uuid.Parse(id)
	helper.PanicIf(err)

	data := h.profileServicev2.FindByID(r.Context(), uid)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}

func (h *ProfileV2Impl) FindOneInv(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	uid, err := uuid.Parse(id)
	helper.PanicIf(err)

	data := h.profileServicev2.FindInvoice(r.Context(), uid)

	webResponse := web.ResponseCategory{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}
