# PDQ
A Go implementation of the [PDQ perceptual hashing algorithm](https://github.com/facebook/ThreatExchange/tree/main/pdq) developed by Meta for detecting visually similar images.

PDQ produces a 256-bit fingerprint from an image. Two images that look similar to a human will have hashes with a low Hamming distance, despite being cropped, compressed, or edited in a minor way. 

## Installation
```bash
go get github.com/MatthewSH/pdq
```

## Quick Start
```go
package main

import (
    "fmt"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "os"

    "github.com/MatthewSH/pdq"
)

func main() {
    f, err := os.Open("photo.jpg")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    img, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    result, err := pdq.Hash(img)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Hash:    %s\n", result.Hash)
    fmt.Printf("Quality: %d\n", result.Quality)
}
```

## Comparing Images
PDQ hashes are compared using Hamming distance. Two images are considered a match when their distance is at or below `DefaultMatchThreshold` (31, recommended by Meta's reference implementation).

```go
r1, _ := pdq.Hash(img1)
r2, _ := pdq.Hash(img2)

if r1.Hash.Distance(r2.Hash) <= pdq.DefaultMatchThreshold {
    fmt.Println("Images are visually similar")
}
```

## Quality Scores
The `Result.Quality` field is a 0–100 score reflecting how much gradient information the image contained. Very small or nearly featureless images will score near zero and produce unreliable hashes. Discard results below `DefaultQualityThreshold` (50, recommended by Meta's reference implementation).

```go
result, err := pdq.Hash(img)
if err != nil || result.Quality < pdq.DefaultQualityThreshold {
    // hash is not reliable
}
```

## Dihedral Hashes
`Result.Dihedrals` contains all 8 rotational and reflective variants of the hash (original, rotate 90/180/270, flip-x, flip-y, flip +diagonal, flip −diagonal). These are useful for rotation-invariant matching.

```go
result, _ := pdq.Hash(img)
for i, h := range result.Dihedrals {
    fmt.Printf("dihedral[%d]: %v\n", i, h)
}
```

## API
### `Hash(img image.Image) (Result, error)`
Computes the PDQ perceptual hash of an image. Returns a `Result` containing the primary hash, quality score, and all 8 dihedral variant hashes.

### `Result`
```go
type Result struct {
    Hash      Hash256       // Primary perceptual hash
    Quality   int           // Quality score (0–100)
    Dihedrals [8][16]uint16 // All 8 dihedral variant hashes
}
```

### `Hash256`
A 256-bit PDQ perceptual hash stored as 32 bytes.

```go
type Hash256 [32]byte

func (h Hash256) String() string                  // 64-char lowercase hex
func (h Hash256) Bytes() []byte                   // Raw bytes
func (h Hash256) IsZero() bool                    // True if hash is empty
func (h Hash256) Distance(other Hash256) int      // Hamming distance (0–256)
```

### Constants
| Constant                  | Value | Description                                     |
|---------------------------|-------|-------------------------------------------------|
| `DefaultMatchThreshold`   | 31    | Max Hamming distance to consider images a match |
| `DefaultQualityThreshold` | 50    | Min quality score for a reliable hash           |

## Running Tests
Unit tests run without any fixtures:

```bash
go test
```

Integration tests compare output against known-good hashes from Meta's reference data. Download the fixtures first:

```bash
go run ./internal/cmd/fetch-testdata
go test -tags integration
```

## How It Works
The hashing pipeline follows the PDQ specification:

1. **Resize** – the image is downscaled to at most 512 × 512 using bilinear interpolation.
2. **Luminance** - RGB pixels are converted to grayscale using standard luma weights.
3. **Jarosz filter** – a separable box filter smooths and decimates the luma plane to 64 × 64.
4. **DCT** - a 2D Discrete Cosine Transform reduces the 64 × 64 matrix to a 16 × 16 block of low-frequency coefficients.
5. **Quantize** – each coefficient is thresholded against the median, producing the final 256-bit hash.

Quality is derived from the gradient energy of the 64 × 64 filtered image before the DCT step.

## License
This PDQ implementation is licensed under the [MIT](LICENSE) license.