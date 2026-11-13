package signal

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

func equalFloats(t *testing.T, name string, got, want []float64, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: len %d, want %d (%v)", name, len(got), len(want), got)
	}
	for i := range want {
		if math.Abs(got[i]-want[i]) > tol {
			t.Fatalf("%s: got %v, want %v", name, got, want)
		}
	}
}

func TestFFTFreq_cases(t *testing.T) {
	// Bin k is k/(n·Δt) (Oppenheim & Schafer §8.1); bins above n/2 alias to k − n.
	got, err := FFTFreq(nil, 4, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	equalFloats(t, "n=4", got, []float64{0, 2.5, -5, -2.5}, 1e-15)
	got, _ = FFTFreq(nil, 5, 1)
	equalFloats(t, "n=5", got, []float64{0, 0.2, 0.4, -0.4, -0.2}, 1e-15)
	got, _ = FFTFreq(nil, 1, 2)
	equalFloats(t, "n=1", got, []float64{0}, 0)
	got, _ = FFTFreq(nil, 4, -1) // negative spacing flips the signs, as numpy does
	equalFloats(t, "negative d", got, []float64{0, -0.25, 0.5, 0.25}, 1e-15)
}

func TestRFFTFreq_cases(t *testing.T) {
	got, err := RFFTFreq(nil, 4, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	equalFloats(t, "n=4", got, []float64{0, 2.5, 5}, 1e-15)
	got, _ = RFFTFreq(nil, 5, 1)
	equalFloats(t, "n=5", got, []float64{0, 0.2, 0.4}, 1e-15)
	got, _ = RFFTFreq(nil, 1, 1)
	equalFloats(t, "n=1", got, []float64{0}, 0)
}

func TestFFTFreq_errors(t *testing.T) {
	for _, f := range []func([]float64, int, float64) ([]float64, error){FFTFreq, RFFTFreq} {
		if _, err := f(nil, 0, 1); !errors.Is(err, ErrDimension) {
			t.Fatalf("n=0: %v", err)
		}
		if _, err := f(nil, -3, 1); !errors.Is(err, ErrDimension) {
			t.Fatalf("n<0: %v", err)
		}
		for _, d := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			if _, err := f(nil, 4, d); !errors.Is(err, ErrNotFinite) {
				t.Fatalf("d=%v: %v", d, err)
			}
		}
		if _, err := f(nil, 4, 0); !errors.Is(err, ErrOutOfRange) {
			t.Fatalf("d=0: %v", err)
		}
		if _, err := f(nil, 4, 5e-324); !errors.Is(err, ErrOverflow) {
			t.Fatalf("subnormal d: %v", err)
		}
	}
}

func BenchmarkFFTFreq(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		dst := make([]float64, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = FFTFreq(dst, n, 0.001)
			}
		})
	}
}

func BenchmarkRFFTFreq(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		dst := make([]float64, n/2+1)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = RFFTFreq(dst, n, 0.001)
			}
		})
	}
}
