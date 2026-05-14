package internal

import (
	"math/bits"

	"github.com/MatthewSH/pdq"
)

const (
	numRows = 64
	numCols = 64
	dims    = numRows * numCols
	scale   = float32(100.0 / 255.0)
)

func abs(x int) int {
	mask := x >> (bits.UintSize - 1)
	return (x ^ mask) - mask
}

func ComputeQuality(matrix []float32) (int, error) {
	if len(matrix) != dims {
		return 0, pdq.ErrMatrixLength
	}

	m := matrix[:dims]
	var gradientSum int

	// Vertical differences (63 x 64)
	for i := 0; i < numRows-1; i++ {
		row := m[i*numCols : i*numCols+numCols]
		next := m[(i+1)*numCols : (i+1)*numCols+numCols]
		for j := 0; j < numCols; j++ {
			gradientSum += abs(int((row[j] - next[j]) * scale))
		}
	}

	// Horizontal differences (64 x 63)
	for i := 0; i < numRows; i++ {
		row := m[i*numCols : i*numCols+numCols]
		for j := 0; j < numCols-1; j++ {
			gradientSum += abs(int((row[j] - row[j+1]) * scale))
		}
	}

	if quality := gradientSum / 90; quality < 100 {
		return quality, nil
	}

	return 100, nil
}
