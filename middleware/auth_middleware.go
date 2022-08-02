package middleware

import (
	"golang_restapi/helper"
	"golang_restapi/model/web"
	"net/http"
)

type AuthMiddleware struct {
	handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{handler: handler}
}

func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if "secret" == r.Header.Get("X-Api-Key") {
		m.handler.ServeHTTP(w, r)
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		webResponse := web.ResponseCategory{
			Code:   http.StatusUnauthorized,
			Status: "unauthorized",
		}

		helper.JSONEncoder(w, webResponse)
	}
}
