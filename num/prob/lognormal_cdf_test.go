package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.9 (F(x) = Φ(ln(x)/σ) in the standard form): μ = 0, σ = 1, x = 1 → 0.5 (hand-checked);
// μ = 1, σ = 0.5, x = 2 → 0.26970493073490953 (150-digit decimal arithmetic). Tolerance 1e-14.
func TestLognormalCDF_sourceExample(t *testing.T) {
	for _, c := range []struct{ x, mu, sigma, want float64 }{{1, 0, 1, 0.5}, {2, 1, 0.5, 0.26970493073490953}} {
		if got, err := LognormalCDF(c.x, c.mu, c.sigma); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("LognormalCDF(%v, %v, %v) = (%v, %v), want (%v, nil)", c.x, c.mu, c.sigma, got, err, c.want)
		}
	}
}

func TestLognormalCDF_cases(t *testing.T) {
	for _, c := range []struct{ x, mu, sigma, want float64 }{{0, 0, 1, 0}, {-3, 0, 1, 0}, {math.E, 1, 2, 0.5}, {math.MaxFloat64, 0, 1, 1}} {
		if got, err := LognormalCDF(c.x, c.mu, c.sigma); err != nil || !within(got, c.want, 1e-15) {
			t.Errorf("LognormalCDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.mu, c.sigma, got, err, c.want)
		}
	}
}

func TestLognormalCDF_errors(t *testing.T) {
	_, err := LognormalCDF(1, 0, -1)
	requireErr(t, "sigma < 0", err, ErrOutOfRange)
	_, err = LognormalCDF(inf, 0, 1)
	requireErr(t, "x Inf", err, ErrNotFinite)
}

func TestLognormalCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "LognormalCDF", func() { sinkFloat, _ = LognormalCDF(2, 1, 0.5) })
}

func BenchmarkLognormalCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LognormalCDF(2, 1, 0.5)
	}
}
