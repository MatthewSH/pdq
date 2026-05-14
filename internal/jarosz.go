package internal

import (
	"fmt"
)

const (
	outSize   = 64
	totalSize = ImageSize * ImageSize
)

// JaroszFilter applies a 2-pass Jarosz blur on the input image
// and returns the downsampled 64x64 result.
// Input must be ImageSize x ImageSize (512x512).
func JaroszFilter(src []float32) ([]float32, error) {
	if len(src) != totalSize {
		return nil, fmt.Errorf(
			"source image dimensions do not match expected: %d vs %d",
			len(src), totalSize,
		)
	}

	windowRows := jaroszWindowSize(ImageSize)
	windowCols := jaroszWindowSize(ImageSize)

	a := make([]float32, totalSize)
	b := make([]float32, totalSize)

	boxAlongRows(src, a, ImageSize, ImageSize, windowRows)
	boxAlongCols(a, b, ImageSize, ImageSize, windowCols)

	boxAlongRows(b, a, ImageSize, ImageSize, windowRows)
	boxAlongCols(a, b, ImageSize, ImageSize, windowCols)

	block := ImageSize / outSize
	center := block / 2

	out := make([]float32, outSize*outSize)
	for y := range outSize {
		srcRow := b[(y*block+center)*ImageSize:]
		outRow := out[y*outSize:]
		for x := range outSize {
			outRow[x] = srcRow[x*block+center]
		}
	}

	return out, nil
}

// jaroszWindowSize computes the 1D box-filter window size for a single pass.
func jaroszWindowSize(oldDimension int) int {
	return (oldDimension + 2*outSize - 1) / (2 * outSize)
}

// boxAlongRows applies a 1D box filter horizontally across every row.
func boxAlongRows(src, dst []float32, numRows, numCols, windowSize int) {
	for y := range numRows {
		box1D(src[y*numCols:], dst[y*numCols:], numCols, 1, windowSize)
	}
}

// boxAlongCols applies a 1D box filter vertically down every column.
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
