package series

import (
	"errors"
	"math"
	"testing"
)

// Hand-worked (Hamilton 1994 §5.2: OLS of x_t on 1 and x_{t−1}): 1, 2, 4, 8 → y = (2, 4, 8) is exactly
// 2·(1, 2, 4), so intercept 0 and φ = 2. 3, 5, 6: y = (5, 6) on x = (3, 5) → φ = 1/2, c = 5 − 1.5 = 3.5.
// Tolerance 1e-14.
func TestAR1Fit_handWorked(t *testing.T) {
	tests := []struct {
		values []float64
		c, phi float64
	}{
		{[]float64{1, 2, 4, 8}, 0, 2},
		{[]float64{3, 5, 6}, 3.5, 0.5},
		{[]float64{-1, 2, -4, 8}, 0, -2},
	}
	for _, tc := range tests {
		c, phi, err := AR1Fit(tc.values)
		if err != nil || !closeRel(c, tc.c, 1e-14) || !closeRel(phi, tc.phi, 1e-14) {
			t.Errorf("%v: = (%v, %v, %v), want (%v, %v)", tc.values, c, phi, err, tc.c, tc.phi)
		}
	}
}

// AR(1) with c = 2, φ = 0.6 from the fixed seed, n = 20000: SE(φ̂) ≈ √((1 − φ²)/n) ≈ 0.0057, tolerance
// 0.03; the mean c/(1 − φ) = 5 is recovered within 0.2 (its SE is about σ/((1 − φ)√n) ≈ 0.018).
func TestAR1Fit_knownPhi(t *testing.T) {
	c, phi, err := AR1Fit(ar1Series(20000, 2, 0.6))
	if err != nil || math.Abs(phi-0.6) > 0.03 || math.Abs(c/(1-phi)-5) > 0.2 {
		t.Errorf("= (%v, %v, %v), want φ ≈ 0.6, mean ≈ 5", c, phi, err)
	}
}

func TestAR1Fit_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrTooFew},
		{"two", []float64{1, 2}, ErrTooFew},
		{"NaN", []float64{1, 2, math.NaN()}, ErrNotFinite},
		{"Inf", []float64{math.Inf(1), 2, 3}, ErrNotFinite},
		{"constant lag", []float64{4, 4, 9}, ErrZeroVariance},
		{"overflow (slope 1e600)", []float64{0, 1e-300, 1e300}, ErrOverflow},
	}
	for _, tc := range tests {
		if _, _, err := AR1Fit(tc.values); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

// Hand-worked: φ^h = ½ → h = ln 2/−ln φ. φ = 0.5 → 1; φ = 0.25 → 0.5; φ = 0.9 → 6.578813478960585
// (Python math.log(2)/-math.log(0.9)). Tolerance 1e-15.
func TestMeanReversionHalfLife_handWorked(t *testing.T) {
	tests := []struct{ phi, want float64 }{
		{0.5, 1},
		{0.25, 0.5},
		{0.9, 6.578813478960585},
		{math.Pow(0.5, 0.1), 10},
	}
	for _, tc := range tests {
		got, err := MeanReversionHalfLife(tc.phi)
		if err != nil || !closeRel(got, tc.want, 1e-14) {
			t.Errorf("φ=%v: = (%v, %v), want %v", tc.phi, got, err, tc.want)
		}
	}
	// Near the edges: tiny φ gives a short finite half-life; φ just below 1 a long finite one.
	if got, err := MeanReversionHalfLife(math.SmallestNonzeroFloat64); err != nil || !(got > 0 && got < 0.01) {
		t.Errorf("smallest φ: (%v, %v)", got, err)
	}
	if got, err := MeanReversionHalfLife(math.Nextafter(1, 0)); err != nil || !(got > 1e15 && !math.IsInf(got, 0)) {
		t.Errorf("φ just below 1: (%v, %v)", got, err)
	}
}

func TestMeanReversionHalfLife_errors(t *testing.T) {
	tests := []struct {
		name string
		phi  float64
		want error
	}{
		{"NaN", math.NaN(), ErrNotFinite},
		{"+Inf", math.Inf(1), ErrNotFinite},
		{"zero", 0, ErrOutOfRange},
		{"negative", -0.5, ErrOutOfRange},
		{"one (unit root)", 1, ErrOutOfRange},
		{"explosive", 1.2, ErrOutOfRange},
	}
	for _, tc := range tests {
		if _, err := MeanReversionHalfLife(tc.phi); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func BenchmarkAR1Fit(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		values := ar1Series(n, 1, 0.7)
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = AR1Fit(values)
			}
		})
	}
}

var sinkHalfLife float64

func BenchmarkMeanReversionHalfLife(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkHalfLife, _ = MeanReversionHalfLife(0.95)
	}
}
