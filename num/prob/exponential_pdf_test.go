package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.7 (f(x) = (1/β)e^(−x/β), λ = 1/β): λ = 2, x = 1 → 2e⁻² = 0.2706705664732254. Hand-checked;
// tolerance 2 ulps.
func TestExponentialPDF_sourceExample(t *testing.T) {
	if got, err := ExponentialPDF(1, 2); err != nil || !within(got, 0.2706705664732254, 5e-16) {
		t.Errorf("ExponentialPDF(1, 2) = (%v, %v), want (0.2706705664732254, nil)", got, err)
	}
}

func TestExponentialPDF_cases(t *testing.T) {
	for _, c := range []struct{ x, rate, want float64 }{{0, 2, 2}, {-1, 2, 0}, {1e300, 1, 0}, {0, math.MaxFloat64, math.MaxFloat64}, {1, 1e-300, 1e-300}} {
		if got, err := ExponentialPDF(c.x, c.rate); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("ExponentialPDF(%v, %v) = (%v, %v), want %v", c.x, c.rate, got, err, c.want)
		}
	}
}

func TestExponentialPDF_errors(t *testing.T) {
	_, err := ExponentialPDF(1, 0)
	requireErr(t, "rate = 0", err, ErrOutOfRange)
	_, err = ExponentialPDF(nan, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = ExponentialPDF(1, inf)
	requireErr(t, "rate Inf", err, ErrNotFinite)
}

func TestExponentialPDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ExponentialPDF", func() { sinkFloat, _ = ExponentialPDF(1, 2) })
}

func BenchmarkExponentialPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ExponentialPDF(1, 2)
	}
}
