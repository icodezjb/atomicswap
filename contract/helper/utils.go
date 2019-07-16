package htlc

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"
)

type SecretHashPair struct {
	Secret uint32
	Hash   [32]byte
}

func NewSecretHashPair() SecretHashPair {
	rand.Seed(time.Now().UnixNano())

	var random = rand.Uint32()

	padded := LeftPad32Bytes(Uint32ToBytes(random))

	return SecretHashPair{
		Secret: random,
		Hash:   sha256.Sum256(padded[:]),
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

func Uint32ToBytes(u uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, u)
	return buf
}
