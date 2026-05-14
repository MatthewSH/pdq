package pdq

import (
	"fmt"
	"image"

	"github.com/MatthewSH/pdq/internal"
)

const (
	// DefaultMatchThreshold is the recommended Hamming distance threshold for
	// considering two PDQ hashes a match. Per the PDQ spec README:
	// "Distance Threshold to consider two hashes similar/matching: <=31"
	DefaultMatchThreshold = 31

	// DefaultQualityThreshold is the minimum quality score for a hash to be
	// considered reliable. Per the PDQ spec README:
	// "Quality Threshold where we recommend discarding hashes: <=49"
	DefaultQualityThreshold = 50
)

// Hash computes the PDQ perceptual hash of an image.
// Returns the primary hash, quality score (0–100), and all 8 dihedral hashes.
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

// packHash packs a [16]uint16 hash into Hash256 bytes.
// The PDQ reference (Python Hash256.__str__, C++ Hash256::format) outputs
// slots from index 15 down to 0, each as a big-endian 16-bit hex word.
// We match this by writing slot 15 into bytes 0-1, ..., slot 0 into bytes 30-31.
// Each uint16 is stored big-endian (high byte first) so that hex.EncodeToString
// produces the canonical PDQ hash string.
func packHash(hash [16]uint16) Hash256 {
	var out Hash256

	for i, v := range hash {
		pos := (15 - i) * 2
		out[pos] = byte(v >> 8) // high byte first (big-endian)
		out[pos+1] = byte(v)
	}

	return out
}
