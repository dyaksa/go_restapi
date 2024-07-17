package entity

import (
	"github.com/dyaksa/encryption-pii/crypto/types"
	"github.com/google/uuid"
)

type Invoice struct {
	ID            uuid.UUID       `db:"id"`
	CustomerName  string          `db:"customer_name"`
	CustomerNameB types.AESChiper `db:"customer_name_b" txt_heap_table:"name_text_heap" bidx_col:"customer_name_bidx"`
	CustomerAddr  string          `db:"customer_address"`
	CustomerAddrB types.AESChiper `db:"customer_address_b" txt_heap_table:"address_text_heap" bidx_col:"customer_address_bidx"`
}

type FindInvoiceParams struct {
	Type    string
	Content string
}

type FindInvoiceNameRows struct {
	ID           uuid.UUID
	Number       string
	CustomerName types.AESChiper
}

type FindInvoiceByBIDXParams struct {
	ColumnHeap string
	Hash       []string
}
