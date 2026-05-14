package internal

import (
	"fmt"
	"image"

	"golang.org/x/image/draw"
)

// PrepareImage conditionally downscales src to at most ImageSize in either
// dimension, then returns the pixel data as a float32 luminance array along
// with the actual processing dimensions.
//
// Small images are processed at their native size. Only images larger than
// ImageSize in either dimension are downscaled.
func PrepareImage(src image.Image) (luma []float32, numRows, numCols int, err error) {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	var rgba *image.RGBA
	if w > ImageSize || h > ImageSize {
		rgba = ResizeBilinear(src, ImageSize)
		numRows, numCols = ImageSize, ImageSize
	} else {
		rgba = ToRGBA(src)
		numRows, numCols = h, w
	}

	luma, err = ToLuminance(rgba, numRows, numCols)
	return
}

// ResizeBilinear resizes the given image using bilinear interpolation to the specified square size (width and height).
func ResizeBilinear(src image.Image, size int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)

	return dst
}

// ToRGBA converts any image.Image to *image.RGBA without resizing.
func ToRGBA(src image.Image) *image.RGBA {
	if rgba, ok := src.(*image.RGBA); ok {
		return rgba
	}
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
	return dst
}

// ToLuminance converts an RGBA image to a grayscale luminance representation.
func ToLuminance(src *image.RGBA, numRows, numCols int) ([]float32, error) {
	bounds := src.Bounds()
	if bounds.Dx() != numCols || bounds.Dy() != numRows {
		return nil, fmt.Errorf("pdq: luminance: image %dx%d does not match expected %dx%d",
			bounds.Dx(), bounds.Dy(), numCols, numRows)
	}

	out := make([]float32, numRows*numCols)
	pixels := src.Pix
	stride := src.Stride

	for y := 0; y < numRows; y++ {
		rowOffset := y * stride
		outOffset := y * numCols

		for x := 0; x < numCols; x++ {
			i := rowOffset + x*4
			r := float32(pixels[i])
			g := float32(pixels[i+1])
			b := float32(pixels[i+2])
			out[outOffset+x] = 0.299*r + 0.587*g + 0.114*b
		}
	}
	return out, nil
}
