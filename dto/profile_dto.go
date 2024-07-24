package dto

import (
	"github.com/google/uuid"
)

type FetchProfileDto struct {
	ID    uuid.UUID `json:"id"`
	Nik   string    `json:"nik"`
	Name  string    `json:"name"`
	Phone string    `json:"phone"`
	Email string    `json:"email"`
	DOB   string    `json:"dob"`
}

type ParamsListProfile struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
