package internal

import "github.com/MatthewSH/pdq"

func Quantize(input []float32, median float32) ([16]uint16, error) {
	if len(input) != 16*16 {
		return [16]uint16{}, pdq.ErrQuantizeLength
	}

	var hash [16]uint16

	for i := range 16 {
		var bits uint16
		row := input[i*16 : (i+1)*16]

		for j, v := range row {
			if v > median {
				bits |= 1 << j
			}
		}

		hash[i] = bits
	}

	return hash, nil
}
