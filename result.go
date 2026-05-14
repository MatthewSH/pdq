package pdq

import "fmt"

type Result struct {
	Hash    Hash256
	Quality int
}

func (r Result) IsValid() bool {
	return !r.Hash.IsZero() && r.Quality > 0
}

func (r Result) String() string {
	return fmt.Sprintf("pdq{hash=%s quality=%d}", r.Hash.String(), r.Quality)
}
