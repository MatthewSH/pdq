package internal

import "errors"

var (
	ErrMatrixLength        = errors.New("pdq: matrix length is not the expected 64x64")
	ErrTorbenElementLength = errors.New("pdq: expected 256 elements for torben median")
	ErrQuantizeLength      = errors.New("pdq: expected 16x16 for DCT quantization")
	ErrDihedralDCTSize     = errors.New("pdq: expected 16x16 for DCT dihedral")
)
