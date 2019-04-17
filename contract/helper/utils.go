package htlc

import (
	"crypto/sha256"
	"math/rand"
	"strconv"
)

type SecretHashPair struct {
	Secret 	string
	Hash	[32]byte
}

func NewSecretHashPair() SecretHashPair  {
	s := strconv.FormatUint(uint64(rand.Uint32()),10)

	return SecretHashPair{
		Secret:s,
		Hash: sha256.Sum256([]byte(s)),
	}
}