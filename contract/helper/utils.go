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

	padded := LeftPad32Bytes([]byte(s))
	return SecretHashPair{
		Secret:s,
		Hash: sha256.Sum256(padded[:]),
	}
}

// LeftPad32Bytes zero-pads slice to the left up to length 32.
func LeftPad32Bytes(slice []byte) [32]byte {
	var padded [32]byte
	if 32 <= len(slice) {
		return padded
	}

	copy(padded[32-len(slice):], slice)

	return padded
}