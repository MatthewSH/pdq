package pdq

import (
	"encoding/hex"

	"github.com/MatthewSH/pdq/internal"
)

// Hash256 is a 256-bit PDQ perceptual hash stored as 32 bytes.
type Hash256 [internal.HashSize]byte

func (h Hash256) Bytes() []byte {
	b := make([]byte, internal.HashSize)
	copy(b, h[:])

	return b
}

// String returns the hash as a 64-character lowercase hex string.
func (h Hash256) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash256) IsZero() bool {
	return h == Hash256{}
}
