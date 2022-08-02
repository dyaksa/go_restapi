package exception

import (
	"github.com/go-playground/validator/v10"
	"golang_restapi/helper"
	"golang_restapi/model/web"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if notFoundError(w, r, err) {
		return

	} else if validationError(w, r, err) {
		return
	} else {
		internalServerError(w, r, err)
	}
}

func validationError(w http.ResponseWriter, r *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		webResponse := web.ResponseCategory{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Data:   exception.Error(),
		}
		helper.JSONEncoder(w, webResponse)
		return true
	} else {
		return false
	}
}

func notFoundError(w http.ResponseWriter, r *http.Request, err interface{}) bool {

	exception, ok := err.(NotFoundError)
	if ok {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		webResponse := web.ResponseCategory{
			Code:   http.StatusNotFound,
			Status: "not found",
			Data:   exception.Error,
		}

		helper.JSONEncoder(w, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	webResponse := web.ResponseCategory{
		Code:   http.StatusInternalServerError,
		Status: "internal server error",
		Data:   err,
	}

	helper.JSONEncoder(w, webResponse)
}
