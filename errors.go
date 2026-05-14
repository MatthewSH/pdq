package pdq

import (
	"errors"
)

var (
	// ErrNilImage is returned when a nil image.Image is passed to Hash.
	ErrNilImage = errors.New("pdq: image must not be nil")
)
