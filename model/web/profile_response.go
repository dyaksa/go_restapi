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

type InvoiceResponse struct {
	ID              uuid.UUID `json:"id"`
	CustomerName    string    `json:"customer_name"`
	CustomerAddress string    `json:"customer_address"`
}

type InvoiceFormatterResponse struct {
	ID            uuid.UUID `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	CustomerName  string    `json:"customer_name"`
}
