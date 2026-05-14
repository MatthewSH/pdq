package internal

import "fmt"

const outSize = 64

// JaroszFilter applies a 2-pass filter to src and returns
// a downsampled 64x64 result.
func JaroszFilter(src []float32, numRows, numCols int) ([]float32, error) {
	if len(src) != numRows*numCols {
		return nil, fmt.Errorf(
			"pdq: jarosz: src length %d does not match %dx%d=%d",
			len(src), numRows, numCols, numRows*numCols,
		)
	}

	windowRows := jaroszWindowSize(numRows)
	windowCols := jaroszWindowSize(numCols)

	a := make([]float32, numRows*numCols)
	b := make([]float32, numRows*numCols)

	boxAlongRows(src, a, numRows, numCols, windowCols)
	boxAlongCols(a, b, numRows, numCols, windowRows)

	boxAlongRows(b, a, numRows, numCols, windowCols)
	boxAlongCols(a, b, numRows, numCols, windowRows)

	return decimate(b, numRows, numCols), nil
}

// jaroszWindowSize computes the 1D box-filter window size for a single pass.
func jaroszWindowSize(oldDimension int) int {
	return (oldDimension + 2*outSize - 1) / (2 * outSize)
}

// decimate samples the center pixel of each output block, matching the
// reference decimateFloat implementation:
//
//	ini = int(((outi + 0.5) * inNumRows) / outNumRows)
func decimate(src []float32, numRows, numCols int) []float32 {
	out := make([]float32, outSize*outSize)
	for outi := range outSize {
		ini := int((float64(outi) + 0.5) * float64(numRows) / float64(outSize))
		srcRow := src[ini*numCols:]
		outRow := out[outi*outSize:]
		for outj := range outSize {
			inj := int((float64(outj) + 0.5) * float64(numCols) / float64(outSize))
			outRow[outj] = srcRow[inj]
		}
	}
	return out
}

func boxAlongRows(src, dst []float32, numRows, numCols, windowSize int) {
	for y := range numRows {
		box1D(src[y*numCols:], dst[y*numCols:], numCols, 1, windowSize)
	}
}

func boxAlongCols(src, dst []float32, numRows, numCols, windowSize int) {
	for x := range numCols {
		box1D(src[x:], dst[x:], numRows, numCols, windowSize)
	}
}

// box1D implements the 4-phase sliding box filter from the reference.
// stride is 1 for rows, numCols for columns.
func box1D(in, out []float32, length, stride, fullWindowSize int) {
	halfWindowSize := (fullWindowSize + 2) / 2
	li, ri, oi, currentWindowSize := 0, 0, 0, 0

	var sum float32

	// accumulate without writing
	for range halfWindowSize - 1 {
		sum += in[ri]
		currentWindowSize++
		ri += stride
	}

	// write with growing window
	for range fullWindowSize - halfWindowSize + 1 {
		sum += in[ri]
		currentWindowSize++
		out[oi] = sum / float32(currentWindowSize)
		ri += stride
		oi += stride
	}

	// write with full window (add right, subtract left)
	for range length - fullWindowSize {
		sum += in[ri]
		sum -= in[li]
		out[oi] = sum / float32(currentWindowSize)
		li += stride
		ri += stride
		oi += stride
	}

	// write with shrinking window
	for range halfWindowSize - 1 {
		sum -= in[li]
		currentWindowSize--
		out[oi] = sum / float32(currentWindowSize)
		li += stride
		oi += stride
	}
}
