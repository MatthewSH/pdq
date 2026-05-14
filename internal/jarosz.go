package internal

import (
	"fmt"
)

const (
	outSize     = 64
	windowSize  = 16
	halfWindow  = windowSize / 2
	kernelWidth = 2*halfWindow + 1
	totalSize   = ImageSize * ImageSize
	block       = ImageSize / outSize
	center      = block / 2
	invKernel   = float32(1.0) / float32(kernelWidth)
	tileSize    = 32
)

// JaroszFilter applies a multi-stage Jarosz blur filter on the input image and returns the resulting filtered image.
func JaroszFilter(src []float32) ([]float32, error) {
	if len(src) != totalSize {
		return nil, fmt.Errorf(
			"source image dimensions do not match expected: %d vs %d",
			len(src), totalSize,
		)
	}

	a := make([]float32, totalSize)
	b := make([]float32, totalSize)

	boxFilterH(src, a, ImageSize)
	transpose(a, b, ImageSize)
	boxFilterH(b, a, ImageSize)
	transpose(a, b, ImageSize)

	boxFilterH(b, a, ImageSize)
	transpose(a, b, ImageSize)
	boxFilterH(b, a, ImageSize)
	transpose(a, b, ImageSize)

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

func boxFilterH(src, dest []float32, size int) {
	for y := range size {
		row := src[y*size : y*size+size]
		drow := dest[y*size : y*size+size]

		s := row[0] * float32(halfWindow+1)
		for j := 1; j <= halfWindow && j < size; j++ {
			s += row[j]
		}
		drow[0] = s * invKernel

		for x := 1; x < size; x++ {
			if jn := x + halfWindow; jn < size {
				s += row[jn]
			}
			if jp := x - halfWindow - 1; jp >= 0 {
				s -= row[jp]
			}
			drow[x] = s * invKernel
		}
	}
}

func transpose(src, dst []float32, size int) {
	for i := 0; i < size; i += tileSize {
		iEnd := min(i+tileSize, size)
		for j := 0; j < size; j += tileSize {
			jEnd := min(j+tileSize, size)
			for ii := i; ii < iEnd; ii++ {
				srcRow := src[ii*size:]
				for jj := j; jj < jEnd; jj++ {
					dst[jj*size+ii] = srcRow[jj]
				}
			}
		}
	}
}
