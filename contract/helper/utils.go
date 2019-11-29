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

// pads zeroes on front of a string until it's 32 bytes or 64 hex characters long
func PadTo32Bytes(s string) string {
	l := len(s)
	//TODO: check l > 64
	for {
		if l == 64 {
			return s
		} else {
			s = "0" + s
			l += 1
		}
	}
}
