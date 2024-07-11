package web

import "github.com/google/uuid"

type ProfileResponse struct {
	ID    uuid.UUID `json:"id"`
	Nik   string    `json:"nik"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Email string    `json:"email"`
	DOB   string    `json:"dob"`
}
