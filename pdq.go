package pdq

import (
	"fmt"
	"image"

	"github.com/MatthewSH/pdq/internal"
)

// Hash computes the PDQ perceptual hash of an image.
// Returns the primary hash, quality score (0–100), and all 8 dihedral transform hashes.
func Hash(img image.Image) (Result, error) {
	if img == nil {
		return Result{}, ErrNilImage
	}

	resized := internal.ResizeBilinear(img, internal.ImageSize)

	luma, err := internal.ToLuminance(resized, internal.ImageSize)
	if err != nil {
		return Result{}, fmt.Errorf("pdq: luminance conversion failed: %w", err)
	}

	filtered, err := internal.JaroszFilter(luma)
	if err != nil {
		return Result{}, fmt.Errorf("pdq: jarosz filter failed: %w", err)
	}

	quality, err := internal.ComputeQuality(filtered)
	if err != nil {
		return Result{}, fmt.Errorf("pdq: quality computation failed: %w", err)
	}

	dct := internal.DCT64To16(filtered)
	dihedrals, err := internal.DihedralHashes(dct)
	if err != nil {
		return Result{}, fmt.Errorf("pdq: dihedral hashes failed: %w", err)
	}

	return Result{
		Hash:      packHash(dihedrals[0]),
		Quality:   quality,
		Dihedrals: dihedrals,
	}, nil
}

func packHash(hash [16]uint16) Hash256 {
	var out Hash256

	for i, v := range hash {
		out[i*2] = byte(v)        // low byte
		out[i*2+1] = byte(v >> 8) // high byte
	}

	return out
}
