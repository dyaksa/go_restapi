package web

import "github.com/google/uuid"

type ProfileRequest struct {
	ID    uuid.UUID `json:"id"`
	Nik   string    `json:"nik" validate:"required"`
	Name  string    `json:"name" validate:"required"`
	Phone string    `json:"phone" validate:"required"`
	Email string    `json:"email" validate:"required,email"`
}
