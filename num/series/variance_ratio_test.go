package series

import (
	"errors"
	"math"
	"testing"
)

// Hand-worked from Lo & MacKinlay 1988 §1.1–1.2: logPrices 0, 1, 3, 2, 4, q = 2, T = 4: μ̂ = 1;
// one-period deviations 0, 1, −2, 1 → σ̂²_a = 6/3 = 2; two-period returns 3, 1, 1 minus 2 → 1, −1, −1,
// m = 2·3·(1 − 2/4) = 3 → σ̄²_c = 1; VR = 0.5; z = √4·(−0.5)/√(2·3·1/6) = −1. Tolerance 1e-15.
func TestVarianceRatio_handWorked(t *testing.T) {
	vr, z, err := VarianceRatio([]float64{0, 1, 3, 2, 4}, 2)
	if err != nil || !closeRel(vr, 0.5, 1e-15) || !closeRel(z, -1, 1e-15) {
		t.Errorf("= (%v, %v, %v), want (0.5, −1)", vr, z, err)
	}
}

// Shifting every log price by a constant (a price scale) changes nothing.
func TestVarianceRatio_shiftInvariant(t *testing.T) {
	vr, z, err := VarianceRatio([]float64{7, 8, 10, 9, 11}, 2)
	if err != nil || !closeRel(vr, 0.5, 1e-14) || !closeRel(z, -1, 1e-14) {
		t.Errorf("= (%v, %v, %v), want (0.5, −1)", vr, z, err)
	}
}

// A Gaussian random walk (cumulated white noise from the fixed seed): VR(q) ≈ 1 and |z| small. A
// mean-reverting AR(1) price (φ = 0.5) has VR(10) well below 1 (asymptotically (1/q)·(1 + 2Σ(1 − k/q)φ^k)
// of the returns' ACF; here bounded loosely by 0.6) and z far below 0.
func TestVarianceRatio_randomWalkAndMeanReversion(t *testing.T) {
	steps := ar1Series(20000, 0, 0)
	walk := make([]float64, len(steps)+1)
	for i, s := range steps {
		walk[i+1] = walk[i] + s
	}
	vr, z, err := VarianceRatio(walk, 10)
	if err != nil || math.Abs(vr-1) > 0.1 || math.Abs(z) > 4 {
		t.Errorf("random walk: (%v, %v, %v), want VR ≈ 1", vr, z, err)
	}
	vr, z, err = VarianceRatio(ar1Series(20001, 0, 0.5), 10)
	if err != nil || vr > 0.6 || z > -4 {
		t.Errorf("AR(1) level: (%v, %v, %v), want VR < 0.6, z ≪ 0", vr, z, err)
	}
}

func TestVarianceRatio_errors(t *testing.T) {
	tests := []struct {
		name   string
		prices []float64
		q      int
		want   error
	}{
		{"NaN", []float64{0, math.NaN(), 1, 2}, 2, ErrNotFinite},
		{"Inf", []float64{0, 1, math.Inf(1), 2}, 2, ErrNotFinite},
		{"q = 1", []float64{0, 1, 3, 2}, 1, ErrOutOfRange},
		{"T = q", []float64{0, 1, 3}, 2, ErrTooFew},
		{"empty", nil, 2, ErrTooFew},
		{"constant returns", []float64{0, 1, 2, 3, 4}, 2, ErrZeroVariance},
		{"overflowing differences", []float64{-math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64}, 2, ErrOverflow},
	}
	for _, tc := range tests {
		if _, _, err := VarianceRatio(tc.prices, tc.q); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func BenchmarkVarianceRatio(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		values := ar1Series(n, 0, 0.9)
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = VarianceRatio(values, 5)
			}
		})
	}
}
