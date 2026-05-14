package pdq

import "encoding/hex"

type Hash256 [HashSize]byte

func (h Hash256) Bytes() []byte {
	b := make([]byte, HashSize)
	copy(b, h[:])

	return b
}

func (h Hash256) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash256) IsZero() bool {
	return h == Hash256{}
}
