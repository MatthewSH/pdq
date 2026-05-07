package pdq

import "errors"

var (
	ErrNilImage       = errors.New("pdq: image must not be nil")
	ErrNotImplemented = errors.New("pdq: not implemented")
)
