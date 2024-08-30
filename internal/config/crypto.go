package config

import "github.com/dyaksa/encryption-pii/crypto"

func NewCrypto() *crypto.Crypto {
	crypto, err := crypto.New(
		crypto.Aes256KeySize,
		crypto.WithInitHeapConnection(),
	)
	if err != nil {
		panic(err)
	}
	return crypto
}
