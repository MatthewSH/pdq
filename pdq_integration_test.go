//go:build integration

package pdq_test

import (
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"

	"github.com/MatthewSH/pdq"
)

type referenceVector struct {
	path string
	hash string
}

type qualityVector struct {
	path       string
	minQuality int
	maxQuality int
}

const tolerance = 16

func fixtureDir() string {
	if d := os.Getenv("TESTDATA_DIR"); d != "" {
		return d
	}
	return "testdata"
}

func TestMain(m *testing.M) {
	dir := fixtureDir()
	sentinel := dir + "/reg-test-input/dih/bridge-1-original.jpg"
	if _, err := os.Stat(sentinel); err != nil {
		fmt.Fprintf(os.Stderr, `
Integration test fixtures are missing.

Run the following to download them:

    go run ./internal/cmd/fetch-testdata

Then re-run:

    go test -tags integration ./...

`)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

var referenceVectors = []referenceVector{
	{"misc-images/c.png", "e64cc9d91e623842f8d1f1d9a398e78c9f199a3bd87924f2b7e11e0bf061b064"},
	{"misc-images/small.jpg", "0007001f003f003f007f00ff00ff00ff01ff01ff01ff03ff03ff03ff03ff03ff"},
	{"misc-images/wee.jpg", "6227401f601ff4ccafcc9fad4b0d95d371a2eb7265a3285234d228ca94deeb2d"},
	{"reg-test-input/labelme-subset/q0003.jpg", "54a977c221d14c1c43ba5e6e21d4a13989a3553f1462611cbb85fda7be83b677"},
	{"reg-test-input/labelme-subset/q0004.jpg", "992d44af36d69e6ca6b812585928bac11def254ef5398c6d07466c9abcc65b92"},
	{"reg-test-input/labelme-subset/q0122.jpg", "cfb2009ddd21c6dab0046a7745b5984757a8a4535b3377aea2591d32b33ff940"},
	{"reg-test-input/labelme-subset/q0291.jpg", "a0fe94f1e5cc1cc8dd855948498dc9243f7ca27336f036d7f212b74bc103c9a7"},
	{"reg-test-input/labelme-subset/q0746.jpg", "1049d96239e24d4dca2c55512b8bdb77425f4dbcf575a0a95555aaab5554aaaa"},
	{"reg-test-input/labelme-subset/q1050.jpg", "489db672e9190276d452aeab41eba20f02375fe4092d88defdf491a5c55c5f70"},
	{"reg-test-input/labelme-subset/q2821.jpg", "b150231ffae4710ffcf4f18bb574b109a576f14bb8543189f8743289f174b109"},
}

var dihedralVectors = []referenceVector{
	{"reg-test-input/dih/bridge-1-original.jpg", "d8f8f0cce0f4a84f0e370a22028f67f0b36e2ed596623e1d33e6b39c4e9c9b22"},
	{"reg-test-input/dih/bridge-2-rotate-90.jpg", "38a50efd71c83f429013d68d0ffffc52e34e0e15ada952a9d29684214aa9e5af"},
	{"reg-test-input/dih/bridge-3-rotate-180.jpg", "2dadda64b5a142e5d362209057da895ae63b8c7fc277b4b766b319361f893188"},
	{"reg-test-input/dih/bridge-4-rotate-270.jpg", "a5f0a457248995e8c9065c275aaa54d8b61ba4bdf8fcfc0387c32f8b0bfc4f05"},
	{"reg-test-input/dih/bridge-5-flipx.jpg", "d8f80f31e0f417b00e37f5dd028f980fb36ed12a9662c1e233e64c634e9c64dd"},
	{"reg-test-input/dih/bridge-6-flipy.jpg", "0dad259bb1a1bd18d362576556da32a1e63b7380c2374b4866b3c6c91b89ce77"},
	{"reg-test-input/dih/bridge-7-flip-plus-1.jpg", "f0a5e10271dcc0bd9c5309720fff018de34ef1e8ada9a956d2967ade1ea91a50"},
	{"reg-test-input/dih/bridge-8-flip-minus-1.jpg", "69f05aa8a4996a17c146a2da5aaaab07b61b5b60f8fc07fc83c3d0740bfcb0fa"},
}

var qualityVectors = []qualityVector{
	// C++ reports quality=0: tiny, nearly featureless image.
	{path: "misc-images/small.jpg", minQuality: 0, maxQuality: 10},
	{path: "misc-images/c.png", minQuality: 50, maxQuality: 100},
	{path: "misc-images/wee.jpg", minQuality: 50, maxQuality: 100},
	// C++ reports quality=3: low-gradient labelme image.
	{path: "reg-test-input/labelme-subset/q0003.jpg", minQuality: 0, maxQuality: 10},
	// C++ reports quality=4.
	{path: "reg-test-input/labelme-subset/q0004.jpg", minQuality: 0, maxQuality: 10},
	{path: "reg-test-input/labelme-subset/q0122.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/labelme-subset/q0291.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/labelme-subset/q0746.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/labelme-subset/q1050.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/labelme-subset/q2821.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-1-original.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-2-rotate-90.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-3-rotate-180.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-4-rotate-270.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-5-flipx.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-6-flipy.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-7-flip-plus-1.jpg", minQuality: 50, maxQuality: 100},
	{path: "reg-test-input/dih/bridge-8-flip-minus-1.jpg", minQuality: 50, maxQuality: 100},
}

func hammingDistance(a, b string) (int, error) {
	ab, err := hex.DecodeString(a)
	if err != nil {
		return 0, err
	}
	bb, err := hex.DecodeString(b)
	if err != nil {
		return 0, err
	}
	var dist int
	for i := range ab {
		x := ab[i] ^ bb[i]
		for x != 0 {
			dist += int(x & 1)
			x >>= 1
		}
	}
	return dist, nil
}

func hashFile(t *testing.T, path string) pdq.Result {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}

	result, err := pdq.Hash(img)
	if err != nil {
		t.Fatalf("hash %s: %v", path, err)
	}
	return result
}

func TestReferenceVectors(t *testing.T) {
	dir := fixtureDir()

	for _, tc := range referenceVectors {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			result := hashFile(t, dir+"/"+tc.path)
			got := result.Hash.String()
			dist, err := hammingDistance(got, tc.hash)
			if err != nil {
				t.Fatalf("hammingDistance: %v", err)
			}
			if dist > tolerance {
				t.Errorf("hamming %d > tolerance %d\n  got:  %s\n  want: %s", dist, tolerance, got, tc.hash)
			} else {
				t.Logf("hamming=%d  %s", dist, got)
			}
		})
	}
}

func TestDihedralVectors(t *testing.T) {
	dir := fixtureDir()

	for _, tc := range dihedralVectors {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			result := hashFile(t, dir+"/"+tc.path)
			got := result.Hash.String()
			dist, err := hammingDistance(got, tc.hash)
			if err != nil {
				t.Fatalf("hammingDistance: %v", err)
			}
			if dist > tolerance {
				t.Errorf("hamming %d > tolerance %d\n  got:  %s\n  want: %s", dist, tolerance, got, tc.hash)
			} else {
				t.Logf("hamming=%d  %s", dist, got)
			}
		})
	}
}

// TestQualityScores verifies quality scores against the C++ reference output.
// Each image has an expected range: genuinely low-quality images (small.jpg,
// q0003, q0004) are expected to score near zero, which is correct behaviour —
// the PDQ quality threshold is a consumer-side filter, not a lower bound on
// all possible inputs.
func TestQualityScores(t *testing.T) {
	dir := fixtureDir()

	for _, tc := range qualityVectors {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			result := hashFile(t, dir+"/"+tc.path)
			q := result.Quality
			if q < tc.minQuality || q > tc.maxQuality {
				t.Errorf("quality=%d, want in [%d, %d]", q, tc.minQuality, tc.maxQuality)
			} else {
				t.Logf("quality=%d (expected [%d, %d])", q, tc.minQuality, tc.maxQuality)
			}
		})
	}
}
