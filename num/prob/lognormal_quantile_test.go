package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.9 (PPF exp(σΦ⁻¹(p)), the median is the scale m = e^μ): μ = 0, σ = 1, u = 0.5 → 1.
// Hand-checked.
func TestLognormalQuantile_sourceExample(t *testing.T) {
	if got, err := LognormalQuantile(0.5, 0, 1); err != nil || got != 1 {
		t.Errorf("LognormalQuantile(0.5, 0, 1) = (%v, %v), want (1, nil)", got, err)
	}
}

func TestLognormalQuantile_cases(t *testing.T) {
	for _, c := range []struct{ u, mu, sigma, want float64 }{{0.5, 2, 1, math.Exp(2)}, {0.975, 0, 1, math.Exp(1.959963984540054)}, {1e-300, -700, 10, 0}} {
		if got, err := LognormalQuantile(c.u, c.mu, c.sigma); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("LognormalQuantile(%v, %v, %v) = (%v, %v), want %v", c.u, c.mu, c.sigma, got, err, c.want)
		}
	}
}

func TestLognormalQuantile_errors(t *testing.T) {
	_, err := LognormalQuantile(0, 0, 1)
	requireErr(t, "u = 0", err, ErrOutOfRange)
	_, err = LognormalQuantile(0.5, 0, -1)
	requireErr(t, "sigma < 0", err, ErrOutOfRange)
	_, err = LognormalQuantile(0.999, 700, 10)
	requireErr(t, "quantile overflows", err, ErrOverflow)
	_, err = LognormalQuantile(0.5, inf, 1)
	requireErr(t, "mu Inf", err, ErrNotFinite)
}

// Oracle: LognormalCDF(LognormalQuantile(u)) = u for u in [0.01, 0.99]. Tolerance 1e-13 relative (the
// round trip passes through exp and log).
func TestLognormalQuantile_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 500 {
		u, mu, sigma := 0.01+0.98*rng.Float64(), (rng.Float64()-0.5)*4, 0.1+2*rng.Float64()
		x, err := LognormalQuantile(u, mu, sigma)
		if err != nil {
			t.Fatal(err)
		}
		if back, _ := LognormalCDF(x, mu, sigma); !within(back, u, 1e-13) {
			t.Fatalf("trial %d: LognormalCDF(LognormalQuantile(%v, %v, %v)) = %v", trial, u, mu, sigma, back)
		}
	}
}

func TestLognormalQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "LognormalQuantile", func() { sinkFloat, _ = LognormalQuantile(0.975, 0, 1) })
}

func BenchmarkLognormalQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LognormalQuantile(0.975, 0, 1)
	}
}
