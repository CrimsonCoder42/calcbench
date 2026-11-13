package series

import (
	"errors"
	"math"
	"testing"
)

// Hand-worked R/S (Hurst 1951; Mandelbrot & Wallis 1969): 1, 2, 3, 4.
// s = 2: chunks (1, 2) and (3, 4): deviations −0.5, 0.5, cumulative −0.5, 0 → R = 0.5, S = 0.5, R/S = 1.
// s = 4: deviations −1.5, −0.5, 0.5, 1.5, cumulative −1.5, −2, −1.5, 0 → R = 2, S = √1.25.
// Two points: H = (ln(2/√1.25) − ln 1)/(ln 4 − ln 2) = 2 − ½·log₂5. Tolerance 1e-14.
func TestHurstExponentRS_handWorked(t *testing.T) {
	h, err := HurstExponentRS([]float64{1, 2, 3, 4}, []int{2, 4})
	if err != nil || !closeRel(h, 0.8390359525563189, 1e-14) {
		t.Errorf("= (%v, %v), want 0.8390359525563189", h, err)
	}
	// Order of the chunk sizes does not matter.
	if h2, err := HurstExponentRS([]float64{1, 2, 3, 4}, []int{4, 2}); err != nil || !closeRel(h2, h, 1e-15) {
		t.Errorf("reversed sizes: (%v, %v), want %v", h2, err, h)
	}
}

// Constant chunks are skipped: 5, 5, 1, 2 with s = 2 uses only (1, 2) → R/S = 1, the same as above, so
// with s = 4 the slope differs only through the s = 4 point (deviations 1.75, 1.75, −2.25, −1.25,
// cumulative 1.75, 3.5, 1.25, 0 → R = 3.5, S = √(12.75/4)). H = ln(3.5/√3.1875)/ln 2.
func TestHurstExponentRS_skipsConstantChunk(t *testing.T) {
	want := math.Log(3.5/math.Sqrt(3.1875)) / math.Ln2
	h, err := HurstExponentRS([]float64{5, 5, 1, 2}, []int{2, 4})
	if err != nil || !closeRel(h, want, 1e-14) {
		t.Errorf("= (%v, %v), want %v", h, err, want)
	}
}

// Independent Gaussian noise: H ≈ 0.5, with the known upward small-sample bias of R/S (Anis & Lloyd 1976)
// pushing it to about 0.55–0.6 at these sizes. Bounds 0.4–0.75.
func TestHurstExponentRS_whiteNoise(t *testing.T) {
	h, err := HurstExponentRS(ar1Series(8192, 0, 0), []int{16, 32, 64, 128, 256, 512, 1024})
	if err != nil || h < 0.4 || h > 0.75 {
		t.Errorf("= (%v, %v), want ≈ 0.5–0.6", h, err)
	}
}

func TestHurstExponentRS_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		sizes  []int
		want   error
	}{
		{"NaN", []float64{1, math.NaN(), 3, 4}, []int{2, 4}, ErrNotFinite},
		{"one size", []float64{1, 2, 3, 4}, []int{2}, ErrTooFew},
		{"no sizes", []float64{1, 2, 3, 4}, nil, ErrTooFew},
		{"size 1", []float64{1, 2, 3, 4}, []int{1, 4}, ErrOutOfRange},
		{"size > n", []float64{1, 2, 3, 4}, []int{2, 5}, ErrOutOfRange},
		{"duplicate size", []float64{1, 2, 3, 4}, []int{2, 2}, ErrOutOfRange},
		{"empty values", nil, []int{2, 4}, ErrOutOfRange},
		{"constant", []float64{3, 3, 3, 3}, []int{2, 4}, ErrZeroVariance},
		{"overflow", []float64{-math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64}, []int{2, 4}, ErrOverflow},
	}
	for _, tc := range tests {
		if _, err := HurstExponentRS(tc.values, tc.sizes); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func BenchmarkHurstExponentRS(b *testing.B) {
	sizes := []int{16, 32, 64, 128, 256}
	for _, n := range []int{1000, 10000, 100000} {
		values := ar1Series(n, 0, 0.3)
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = HurstExponentRS(values, sizes)
			}
		})
	}
}
