package pdq_test

import (
	"errors"
	"image"
	"image/color"
	"testing"

	"github.com/MatthewSH/pdq"
)

func solidImage(size int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func gradientImage(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			v := uint8(x * 255 / (size - 1))
			img.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
		}
	}
	return img
}

func TestHash_NilImage(t *testing.T) {
	_, err := pdq.Hash(nil)
	if err == nil {
		t.Fatal("expected error for nil image, got nil")
	}
	if !errors.Is(err, pdq.ErrNilImage) {
		t.Errorf("errors.Is(err, ErrNilImage) = false; err = %v", err)
	}
}

func TestHash_SolidColorNoError(t *testing.T) {
	_, err := pdq.Hash(solidImage(64, color.RGBA{R: 128, G: 64, B: 32, A: 255}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHash_QualityInRange(t *testing.T) {
	result, err := pdq.Hash(gradientImage(256))
	if err != nil {
		t.Fatal(err)
	}
	if result.Quality < 0 || result.Quality > 100 {
		t.Errorf("quality = %d, want in [0, 100]", result.Quality)
	}
}

func TestHash_Deterministic(t *testing.T) {
	img := gradientImage(256)
	r1, err := pdq.Hash(img)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := pdq.Hash(img)
	if err != nil {
		t.Fatal(err)
	}
	if r1.Hash != r2.Hash {
		t.Errorf("same image produced different hashes:\n  %s\n  %s", r1.Hash, r2.Hash)
	}
	if r1.Quality != r2.Quality {
		t.Errorf("same image produced different quality: %d vs %d", r1.Quality, r2.Quality)
	}
}

func TestHash_EightDihedrals(t *testing.T) {
	result, err := pdq.Hash(gradientImage(256))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Dihedrals) != 8 {
		t.Errorf("got %d dihedrals, want 8", len(result.Dihedrals))
	}
}

func TestHash_PrimaryHashMatchesDihedral0(t *testing.T) {
	result, err := pdq.Hash(gradientImage(256))
	if err != nil {
		t.Fatal(err)
	}
	var repacked pdq.Hash256
	for i, v := range result.Dihedrals[0] {
		pos := (15 - i) * 2
		repacked[pos] = byte(v >> 8)
		repacked[pos+1] = byte(v)
	}
	if result.Hash != repacked {
		t.Errorf("Hash != packed Dihedrals[0]\n  Hash:     %s\n  repacked: %s", result.Hash, repacked)
	}
}

func TestHash_StringLength(t *testing.T) {
	result, err := pdq.Hash(gradientImage(64))
	if err != nil {
		t.Fatal(err)
	}
	if s := result.Hash.String(); len(s) != 64 {
		t.Errorf("hash string length = %d, want 64", len(s))
	}
}

func TestHash_DifferentImagesProduceDifferentHashes(t *testing.T) {
	r1, err := pdq.Hash(solidImage(64, color.RGBA{R: 0, G: 0, B: 0, A: 255}))
	if err != nil {
		t.Fatal(err)
	}
	r2, err := pdq.Hash(gradientImage(256))
	if err != nil {
		t.Fatal(err)
	}
	if r1.Hash == r2.Hash {
		t.Error("distinctly different images produced identical hashes")
	}
}

func TestHash_SimilarImagesCloseHamming(t *testing.T) {
	img1 := radialImage(256)
	img2 := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			c := img1.(*image.RGBA).RGBAAt(x, y)
			c.R = uint8(min(int(c.R)+5, 255))
			c.G = uint8(min(int(c.G)+5, 255))
			c.B = uint8(min(int(c.B)+5, 255))
			img2.SetRGBA(x, y, c)
		}
	}

	r1, err := pdq.Hash(img1)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := pdq.Hash(img2)
	if err != nil {
		t.Fatal(err)
	}

	var dist int
	for i, w := range r1.Dihedrals[0] {
		xor := w ^ r2.Dihedrals[0][i]
		for xor != 0 {
			dist += int(xor & 1)
			xor >>= 1
		}
	}

	t.Logf("hamming distance between similar images: %d", dist)
	if dist > pdq.DefaultMatchThreshold {
		t.Errorf("similar images have hamming distance %d > threshold %d", dist, pdq.DefaultMatchThreshold)
	}
}

func radialImage(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	cx, cy := float64(size)/2, float64(size)/2
	maxR := cx
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx, dy := float64(x)-cx, float64(y)-cy
			r := (dx*dx + dy*dy)
			v := uint8(255 * (1 - r/(maxR*maxR)))
			img.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
		}
	}
	return img
}
