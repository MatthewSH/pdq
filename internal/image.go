package internal

import (
	"fmt"
	"image"

	"golang.org/x/image/draw"
)

// ResizeBilinear resizes the given image using bilinear interpolation to the specified square size (width and height).
func ResizeBilinear(src image.Image, size int) image.RGBA64Image {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)

	return dst
}

// ToLuminance converts an RGBA image to a grayscale luminance representation and verifies the input image dimensions.
func ToLuminance(src image.RGBA, size int) ([]float32, error) {
	bounds := src.Bounds()
	if (bounds.Dx() != size) || (bounds.Dy() != size) {
		return nil, fmt.Errorf("source image dimensions do not match expected: %dx%d vs %dx%d", bounds.Dx(), bounds.Dy(), size, size)
	}

	out := make([]float32, size*size)
	pixels := src.Pix
	stride := src.Stride

	for y := 0; y < size; y++ {
		rowOffset := y * stride
		outOffset := y * size

		for x := 0; x < size; x++ {
			i := rowOffset + x*4
			r := float32(pixels[i])
			g := float32(pixels[i+1])
			b := float32(pixels[i+2])

			out[outOffset+x] = 0.299*r + 0.587*g + 0.114*b
		}
	}

	return out, nil
}
