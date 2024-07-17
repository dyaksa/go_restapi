package entity

import (
	"github.com/dyaksa/encryption-pii/crypt/types"
	"github.com/google/uuid"

	t "github.com/dyaksa/encryption-pii/crypto/types"
)

type ProfileV2 struct {
	ID    uuid.UUID   `db:"id"`
	Nik   t.AESChiper `db:"nik"`
	Name  t.AESChiper `db:"name" bidx_col:"name_bidx" txt_heap_table:"profile_text_heap"`
	Phone t.AESChiper `db:"phone"`
	Email t.AESChiper `db:"email" bidx_col:"email_bidx" txt_heap_table:"email_text_heap"`
	DOB   t.AESChiper `db:"dob"`
}

type FetchProfileRowV2 struct {
	ID    uuid.UUID
	Nik   t.AESChiper
	Name  t.AESChiper
	Phone t.AESChiper
	Email t.AESChiper
	DOB   t.AESChiper
}

type FetchInvoiceRow struct {
	ID                uuid.UUID
	CustomerName      string
	CustomerNameCrypt t.AESChiper
	CustomerAddress   string
	CustomerAddrCrypt t.AESChiper
}

type FetchInvoiceRowsCrypt struct {
	IDField      uuid.UUID
	CustomerName t.AESChiper
	CustomerAddr t.AESChiper
}

type Profile struct {
	ID    uuid.UUID        `db:"id"`
	Nik   types.AEADString `db:"nik"`
	Name  types.AEADString `db:"name" bidx_col:"name_bidx" txt_heap_table:"profile_text_heap"`
	Phone types.AEADString `db:"phone"`
	Email types.AEADString `db:"email" bidx_col:"email_bidx" txt_heap_table:"email_text_heap"`
	DOB   types.AEADString `db:"dob"`
}

type FetchProfileRow struct {
	ID    uuid.UUID
	Nik   types.AEADString
	Name  types.AEADString
	Phone types.AEADString
	Email types.AEADString
	DOB   types.AEADString
}

type FetchProfileParams struct {
	ID uuid.UUID
}

// type FindProfileByIDParams struct {
// 	ID []uuid.UUID
// }

// type FindProfileByNameParams struct {
// 	NameBidx types.BIDXString
// }

// type FindProfileByEmailParams struct {
// 	EmailBidx types.BIDXString
// }

type FindProfilesByNameRow struct {
	ID    uuid.UUID
	Nik   types.AEADString
	Name  types.AEADString
	Phone types.AEADString
	Email types.AEADString
	Dob   types.AEADString
}

type FindProfileByParams struct {
	Type    string
	Content string
}

type FindProfileByBIDXParams struct {
	ColumnHeap string
	Hash       []string
}
