package internal

// DihedralHashes generates 8 quantized hash representations by applying
// dihedral transformations to a 16x16 DCT array.
//
// Sign conventions are taken directly from the C++ reference (pdqhashing.cpp):
//
//	orig      rot90     rot180    rot270
//	noxpose   xpose     noxpose   xpose
//	+ + + +   - + - +   + - + -   - - - -
//	+ + + +   - + - +   - + - +   + + + +
//	+ + + +   - + - +   + - + -   - - - -
//	+ + + +   - + - +   - + - +   + + + +
//
//	flipx     flipy     flipplus  flipminus
//	noxpose   noxpose   xpose     xpose
//	- - - -   - + - +   + + + +   + - + -
//	+ + + +   - + - +   + + + +   - + - +
//	- - - -   - + - +   + + + +   + - + -
//	+ + + +   - + - +   + + + +   - + - +
//
// For transposing transforms the C++ writes B[j][i] = ±A[i][j].
// Our closures receive (outRow, outCol); the mapping is outRow=j_cpp,
// outCol=i_cpp, so at(outCol, outRow) reads A[i][j].
func DihedralHashes(dct []float32) ([8][16]uint16, error) {
	if len(dct) != 16*16 {
		return [8][16]uint16{}, ErrDihedralDCTSize
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
		func(i, j int) float32 {
			return at(i, j)
		},

		// Rotate 90 CCW: transpose, negate even-row output
		func(i, j int) float32 {
			if i&1 != 0 {
				return at(j, i)
			}
			return -at(j, i)
		},

		// Rotate 180: no transpose, negate where (i+j) is odd
		func(i, j int) float32 {
			if (i+j)&1 != 0 {
				return -at(i, j)
			}
			return at(i, j)
		},

		// Rotate 270 CCW: transpose, negate even-col output
		func(i, j int) float32 {
			if j&1 != 0 {
				return at(j, i)
			}
			return -at(j, i)
		},

		// Flip X (top/bottom swap): no transpose, negate even-row output
		func(i, j int) float32 {
			if i&1 != 0 {
				return at(i, j)
			}
			return -at(i, j)
		},

		// Flip Y (left/right mirror): no transpose, negate even-col output
		func(i, j int) float32 {
			if j&1 != 0 {
				return at(i, j)
			}
			return -at(i, j)
		},

		// Flip +diagonal: transpose only, no sign change
		func(i, j int) float32 {
			return at(j, i)
		},

		// Flip -diagonal: transpose, negate where (i+j) is odd
		func(i, j int) float32 {
			if (i+j)&1 != 0 {
				return -at(j, i)
			}
			return at(j, i)
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
