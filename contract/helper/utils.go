package htlc

import (
	"crypto/sha256"
	"math/rand"
)

type SecretHashPair struct {
	Secret [32]byte
	Hash   [32]byte
}

func NewSecretHashPair() *SecretHashPair {

	s := new(SecretHashPair)
	rand.Read(s.Secret[:])

	s.Hash = sha256.Sum256(s.Secret[:])

	return s
}
