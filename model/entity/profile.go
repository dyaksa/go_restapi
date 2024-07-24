package entity

import (
	"github.com/dyaksa/encryption-pii/crypto/types"
	"github.com/google/uuid"
)

type ProfileV2 struct {
	ID    uuid.UUID       `db:"id"`
	Nik   types.AESCipher `db:"nik"`
	Name  types.AESCipher `db:"name" bidx_col:"name_bidx" txt_heap_table:"profile_text_heap"`
	Phone types.AESCipher `db:"phone"`
	Email types.AESCipher `db:"email" bidx_col:"email_bidx" txt_heap_table:"email_text_heap"`
	DOB   types.AESCipher `db:"dob"`
}

type FetchProfileRowV2 struct {
	ID    uuid.UUID
	Nik   types.AESCipher
	Name  types.AESCipher
	Phone types.AESCipher
	Email types.AESCipher
	DOB   types.AESCipher
}

type FetchInvoiceRow struct {
	ID                uuid.UUID
	CustomerName      string
	CustomerNameCrypt types.AESCipher
	CustomerAddress   string
	CustomerAddrCrypt types.AESCipher
}

type FetchInvoiceRowsCrypt struct {
	IDField      uuid.UUID
	CustomerName types.AESCipher
	CustomerAddr types.AESCipher
}

type Profile struct {
	ID        uuid.UUID       `db:"id"`
	Nik       types.AESCipher `db:"nik"`
	NikBidx   string          `db:"nik_bidx" txt_heap_table:"nik_text_heap"`
	Name      types.AESCipher `db:"name"`
	NameBidx  string          `db:"name_bidx" txt_heap_table:"name_text_heap"`
	Phone     types.AESCipher `db:"phone"`
	PhoneBidx string          `db:"phone_bidx" txt_heap_table:"phone_text_heap"`
	Email     types.AESCipher `db:"email"`
	EmailBidx string          `db:"email_bidx" txt_heap_table:"email_text_heap"`
	DOB       types.AESCipher `db:"dob"`
}

type FetchProfileRow struct {
	ID    uuid.UUID
	Nik   types.AESCipher
	Name  types.AESCipher
	Phone types.AESCipher
	Email types.AESCipher
	DOB   types.AESCipher
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
	Nik   types.AESCipher
	Name  types.AESCipher
	Phone types.AESCipher
	Email types.AESCipher
	Dob   types.AESCipher
}

type FindProfileByParams struct {
	Type    string
	Content string
}

type FindProfileByBIDXParams struct {
	ColumnHeap string
	Hash       []string
}
