package internal

import "github.com/MatthewSH/pdq"

func DihedralHashes(dct []float32) ([8][16]uint16, error) {
	if len(dct) != 16*16 {
		return [8][16]uint16{}, pdq.ErrDihedralDCTSize
	}

	var buf [16 * 16]float32
	at := func(i, j int) float32 {
		return dct[i*16+j]
	}

	hash := func(get func(i, j int) float32) ([16]uint16, error) {
		for i := range 16 {
			for j := range 16 {
				buf[i*16+j] = get(i, j)
			}
		}
		median, err := TorbenMedian(buf[:])
		if err != nil {
			return [16]uint16{}, err
		}
		return Quantize(buf[:], median)
	}

	transforms := []func(i, j int) float32{
		// original
		func(i, j int) float32 {
			return at(i, j)
		},

		// rotate 90
		func(i, j int) float32 {
			return sign(i&1 == 1) * at(j, i)
		},

		// rotate 180
		func(i, j int) float32 {
			return sign((i+j)&1 == 0) * at(i, j)
		},

		// rotate 270
		func(i, j int) float32 {
			return sign(j&1 == 1) * at(j, i)
		},

		// flip X
		func(i, j int) float32 {
			return sign(i&1 == 1) * at(i, j)
		},

		// flip Y
		func(i, j int) float32 {
			return sign(j&1 == 1) * at(i, j)
		},

		// flip +diagonal
		func(i, j int) float32 {
			return at(j, i)
		},

		// flip -diagonal
		func(i, j int) float32 {
			return sign((i+j)&1 == 0) * at(j, i)
		},
	}

	var out [8][16]uint16
	for k, t := range transforms {
		h, err := hash(t)
		if err != nil {
			return [8][16]uint16{}, err
		}
		out[k] = h
	}

	return out, nil
}
func sign(positive bool) float32 {
	if positive {
		return 1
	}
	return -1
}
