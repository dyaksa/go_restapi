package entity

import "github.com/google/uuid"

type ProfileNonPII struct {
	ID    uuid.UUID `db:"id"`
	Nik   string    `db:"nik"`
	Name  string    `db:"name"`
	Phone string    `db:"phone"`
	Email string    `db:"email"`
}
