package web

import "github.com/google/uuid"

type ProfileRequest struct {
	ID    uuid.UUID `json:"id"`
	Nik   string    `json:"nik" validate:"required"`
	Name  string    `json:"name" validate:"required"`
	Phone string    `json:"phone" validate:"required"`
	Email string    `json:"email" validate:"required,email"`
}

type ProfileResponse struct {
	ID    uuid.UUID `json:"id"`
	Nik   string    `json:"nik"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Email string    `json:"email"`
}

type ProfileQueryParam struct {
	Page  string `query:"page"`
	Limit string `query:"limit"`
	Key   string `query:"key"`
	Value string `query:"value"`
}
