package internal

import (
	"math"
	"testing"
)

func TestJaroszFilter_WrongSize(t *testing.T) {
	_, err := JaroszFilter(make([]float32, 100))
	if err == nil {
		t.Fatal("expected error for wrong input size")
	}
}

func TestJaroszFilter_OutputSize(t *testing.T) {
	src := make([]float32, ImageSize*ImageSize)
	out, err := JaroszFilter(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 64*64 {
		t.Errorf("output length = %d, want %d", len(out), 64*64)
	}
}

func TestJaroszFilter_UniformInput(t *testing.T) {
	const value = 128.0
	src := make([]float32, ImageSize*ImageSize)
	for i := range src {
		src[i] = value
	}
	out, err := JaroszFilter(src)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range out {
		if math.Abs(float64(v-value)) > 1e-3 {
			t.Errorf("out[%d] = %v, want %v", i, v, value)
			break
		}
	}
}

func TestJaroszFilter_OutputInInputRange(t *testing.T) {
	src := make([]float32, ImageSize*ImageSize)
	for i := range src {
		src[i] = float32(i % 256)
	}
	out, err := JaroszFilter(src)
	if err != nil {
		t.Fatal(err)
	}
	const minVal, maxVal = 0.0, 255.0
	for i, v := range out {
		if v < minVal-1e-3 || v > maxVal+1e-3 {
			t.Errorf("out[%d] = %v out of range [%v, %v]", i, v, minVal, maxVal)
			break
		}
	}
}

func TestJaroszWindowSize(t *testing.T) {
	if got := jaroszWindowSize(ImageSize); got != 4 {
		t.Errorf("jaroszWindowSize(%d) = %d, want 4", ImageSize, got)
	}
}

func TestBox1D_WriteCount(t *testing.T) {
	in := make([]float32, 16)
	out := make([]float32, 16)
	for i := range in {
		in[i] = float32(i + 1)
	}
	box1D(in, out, 16, 1, 4)
	for i, v := range out {
		if v == 0 {
			t.Errorf("out[%d] was not written (zero)", i)
		}
	}
}
