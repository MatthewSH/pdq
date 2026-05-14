package internal

import (
	"fmt"

	"github.com/MatthewSH/pdq"
)

const (
	outSize     = 64
	windowSize  = 16
	halfWindow  = windowSize / 2
	kernelWidth = 2*halfWindow + 1
	block       = pdq.ImageSize / outSize
	center      = block / 2
)

func JaroszFilter(src []float32) ([]float32, error) {
	totalSize := pdq.ImageSize * pdq.ImageSize

	if len(src) != totalSize {
		return nil, fmt.Errorf("source image dimensions do not match expected: %d vs %d", len(src), totalSize)
	}

	temp := make([]float32, totalSize)
	dest := make([]float32, totalSize)

	boxFilterH(src, dest, pdq.ImageSize)
	boxFilterV(dest, temp, pdq.ImageSize)
	boxFilterH(temp, dest, pdq.ImageSize)
	boxFilterV(dest, temp, pdq.ImageSize)

	out := make([]float32, outSize*outSize)
	for y := 0; y < outSize; y++ {
		sy := y*block + center
		rowOffset := sy * pdq.ImageSize
		outOffset := y * outSize

		for x := 0; x < outSize; x++ {
			sx := x*block + center
			out[outOffset+x] = dest[rowOffset+sx]
		}
	}

	return out, nil
}

func boxFilterH(src, dest []float32, size int) {
	for y := 0; y < size; y++ {
		rowOffset := y * size
		row := src[rowOffset : rowOffset+size]

		s := row[0] * float32(halfWindow)

		for j := 1; j <= halfWindow && j < size; j++ {
			s += row[j]
		}

		s += row[0]
		dest[rowOffset] = s / float32(kernelWidth)

		for x := 1; x < size; x++ {
			jn := x + halfWindow
			jp := x - halfWindow - 1

			var lead, trail float32

			if jn < size {
				lead = row[jn]
			}

			if jp >= 0 {
				trail = row[jp]
			}

			s += lead - trail
			dest[rowOffset+x] = s / float32(kernelWidth)
		}
	}
}

func boxFilterV(src, dest []float32, size int) {
	for x := 0; x < size; x++ {
		s := src[x] * float32(halfWindow)

		for i := 1; i <= halfWindow && i < size; i++ {
			s += src[i*size+x]
		}

		s += src[x]
		dest[x] = s / float32(kernelWidth)

		for y := 1; y < size; y++ {
			in := y + halfWindow
			ip := y - halfWindow - 1

			var lead, trail float32

			if in < size {
				lead = src[in*size+x]
			}

			if ip >= 0 {
				trail = src[ip*size+x]
			}

			s += lead - trail
			dest[y*size+x] = s / kernelWidth
		}
	}
}
