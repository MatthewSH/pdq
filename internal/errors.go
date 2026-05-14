package internal

import "errors"

var (
	ErrNilImage            = errors.New("pdq: image must not be nil")
	ErrNotImplemented      = errors.New("pdq: not implemented")
	ErrDecodeFailed        = errors.New("pdq: decoding failed")
	ErrLowQuality          = errors.New("pdq: image is below quality threshold")
	ErrMatrixLength        = errors.New("pdq: matrix length is not the expected 64x64")
	ErrTorbenElementLength = errors.New("pdq: expected 256 elements for torben median")
	ErrQuantizeLength      = errors.New("pdq: expected 16x16 for DCT quantization")
	ErrDihedralDCTSize     = errors.New("pdq: expected 16x16 for DCT dihedral")
)
