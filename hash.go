package pdq

import (
	"encoding/hex"

	"github.com/MatthewSH/pdq/internal"
)

type Hash256 [internal.HashSize]byte

func (h Hash256) Bytes() []byte {
	b := make([]byte, internal.HashSize)
	copy(b, h[:])

	return b
}

func (h Hash256) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash256) IsZero() bool {
	return h == Hash256{}
}
