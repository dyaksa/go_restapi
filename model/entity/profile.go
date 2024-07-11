package entity

import (
	"github.com/dyaksa/encryption-pii/crypt/types"
	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID
	Nik       types.AEADString
	NikBidx   types.BIDXString
	Name      types.AEADString
	NameBidx  string
	Phone     types.AEADString
	PhoneBidx types.BIDXString
	Email     types.AEADString
	EmailBidx string
	DOB       types.AEADString
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
	Hash []string
}
