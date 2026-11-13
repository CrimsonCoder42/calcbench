package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 21: F(λ) = 1 − e⁻¹ for every shape; small x keeps relative accuracy (expm1).
func TestWeibullCDF_sourceExample(t *testing.T) {
	for _, k := range []float64{0.5, 1, 2, 7} {
		if got, err := WeibullCDF(3, k, 3); err != nil || !within(got, 1-math.Exp(-1), 1e-15) {
			t.Errorf("k %v: (%v, %v)", k, got, err)
		}
	}
	// F(x) = 1 − e^(−x) = −expm1(−x) for k = λ = 1; at x = 1e-10 that is 1e-10 − 5e-21 = 9.9999999995e-11,
	// which differs from 1e-10 by 5e-11 relative, so the expectation is the expm1 value, not x.
	if got, _ := WeibullCDF(1e-10, 1, 1); !within(got, -math.Expm1(-1e-10), 1e-15) {
		t.Errorf("small x: %v", got)
	}
	if got, _ := WeibullCDF(-2, 1, 1); got != 0 {
		t.Errorf("x < 0: %v", got)
	}
	_, err := WeibullCDF(1, 1, 0)
	requireErr(t, "λ = 0", err, ErrOutOfRange)
	_, err = WeibullCDF(inf, 1, 1)
	requireErr(t, "x Inf", err, ErrNotFinite)
}

func BenchmarkWeibullCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = WeibullCDF(1.5, 2, 1)
	}
}
