package entity

import (
	"github.com/dyaksa/encryption-pii/crypt/types"
	"github.com/google/uuid"
)

type TextHeap struct {
	ID      uuid.UUID
	Content string
	Type    string
	Hash    types.HMACString
}

type FindTextHeapRow struct {
	ID      uuid.UUID
	Content string
	Hash    string
}

type HashExist struct {
	IsExist bool
}

type FindTextHeapByContentParams struct {
	Content string
}

type FindTextHeapByHashParams struct {
	Hash string
}
