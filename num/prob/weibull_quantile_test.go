package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 21, inverted: Q(1 − e⁻¹) = λ for every shape; round trip CDF(Q(u)) = u.
func TestWeibullQuantile_sourceExample(t *testing.T) {
	if got, err := WeibullQuantile(1-math.Exp(-1), 2, 1); err != nil || !within(got, 1, 1e-15) {
		t.Errorf("(1 − e⁻¹, 2, 1): (%v, %v)", got, err)
	}
	for _, u := range []float64{0, 1e-12, 0.3, 0.99} {
		x, _ := WeibullQuantile(u, 1.5, 2)
		if c, _ := WeibullCDF(x, 1.5, 2); !within(c, u, 1e-14) {
			t.Errorf("u %v: CDF(Q) = %v", u, c)
		}
	}
	_, err := WeibullQuantile(1, 2, 1)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = WeibullQuantile(-1, 2, 1)
	requireErr(t, "u < 0", err, ErrProbability)
	_, err = WeibullQuantile(0.5, 2, nan)
	requireErr(t, "λ NaN", err, ErrNotFinite)
	_, err = WeibullQuantile(0.99, 0.001, 1)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkWeibullQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = WeibullQuantile(0.9, 2, 1)
	}
}
