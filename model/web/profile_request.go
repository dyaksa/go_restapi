package web

type ProfileRequest struct {
	Nik   string `json:"nik" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
