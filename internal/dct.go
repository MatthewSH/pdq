package internal

import "math"

const (
	dctRows = 16
	dctCols = 64
)

var dctMatrix = func() [dctRows * dctCols]float32 {
	phaseStep := math.Pi / 128.0
	scale64 := 2.0 / float64(dctCols)
	scaleF := float32(math.Sqrt(scale64))

	var m [dctRows * dctCols]float32

	for i := range dctRows {
		fi := float64(i + 1)
		row := m[i*dctCols : (i+1)*dctCols]

		for j := range dctCols {
			row[j] = scaleF * float32(math.Cos(phaseStep*fi*float64(2*j+1)))
		}
	}

	return m
}()

// DCT64to16 performs a 2D Discrete Cosine Transform (DCT) to convert a 64x64 input into a reduced 16x16 output matrix.
func DCT64to16(input []float32) []float32 {
	out := make([]float32, dctRows*dctRows)

	var at [dctCols * dctCols]float32
	for i := range dctCols {
		for j, v := range input[i*dctCols : (i+1)*dctCols] {
			at[j*dctCols+i] = v
		}
	}

	var t [dctRows * dctCols]float32
	for i := range dctRows {
		dRow := dctMatrix[i*dctCols : (i+1)*dctCols]
		for j := range dctCols {
			atRow := at[j*dctCols : (j+1)*dctCols]
			var sum float32
			for k := range dctCols {
				sum += dRow[k] * atRow[k]
			}

			t[i*dctCols+j] = sum
		}
	}

	for i := range dctRows {
		tRow := t[i*dctCols : (i+1)*dctCols]
		for j := range dctRows {
			dRow := dctMatrix[j*dctCols : (j+1)*dctCols]
			var sum float32
			for k := range dctCols {
				sum += dRow[k] * tRow[k]
			}

			out[i*dctRows+j] = sum
		}
	}

	return out
}
