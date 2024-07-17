package handler

import (
	"golang_restapi/helper"
	"golang_restapi/model/web"
	"golang_restapi/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type InvoiceHandler interface {
	FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Finding(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type InvoiceHandlerImpl struct {
	invService service.InvoiceService
}

func NewInvoiceHandler(invService service.InvoiceService) *InvoiceHandlerImpl {
	return &InvoiceHandlerImpl{invService: invService}
}

func (h *InvoiceHandlerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_, err := h.invService.FindAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *InvoiceHandlerImpl) Finding(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	queryValues := r.URL.Query()

	data, err := h.invService.Search(r.Context(), queryValues)
	helper.PanicIf(err)

	webResponse := web.ResponseCategory{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   data,
	}

	helper.JSONEncoder(w, webResponse)
}
