package entity

import (
	"github.com/dyaksa/encryption-pii/crypto/types"
	"github.com/google/uuid"
)

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
}
