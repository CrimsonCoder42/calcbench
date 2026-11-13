package prob

import (
	"math"
	"testing"
)

// Hill (1977): F(μ) = ½ by symmetry; κ = 0 is uniform on [μ − π, μ + π]. Every quantile is checked by
// the defining property F(Q(u)) = u against VonMisesCDF (absolute, as the CDF's error is absolute).
func TestVonMisesQuantile(t *testing.T) {
	if got, err := VonMisesQuantile(0.5, 1.25, 3); err != nil || got != 1.25 {
		t.Errorf("u = ½: (%v, %v), want μ = 1.25", got, err)
	}
	if got, _ := VonMisesQuantile(0.75, 0, 0); !within(got, math.Pi/2, 1e-15) {
		t.Errorf("κ = 0, u = 0.75: %v, want π/2", got)
	}
	if got, _ := VonMisesQuantile(0, 0.5, 2); got != 0.5-math.Pi {
		t.Errorf("u = 0: %v", got)
	}
	if got, _ := VonMisesQuantile(1, 0.5, 2); got != 0.5+math.Pi {
		t.Errorf("u = 1: %v", got)
	}
	for _, kappa := range []float64{0.1, 1, 2, 10, 100, 1e4} {
		for _, u := range []float64{1e-6, 0.01, 0.1, 0.3, 0.5, 0.7, 0.9, 0.99, 1 - 1e-6} {
			x, err := VonMisesQuantile(u, 0.3, kappa)
			if err != nil {
				t.Errorf("κ = %v, u = %v: %v", kappa, u, err)
				continue
			}
			f, err := VonMisesCDF(x, 0.3, kappa)
			if err != nil || math.Abs(f-u) > 1e-13 {
				t.Errorf("κ = %v, u = %v: F(Q) = (%v, %v)", kappa, u, f, err)
			}
		}
	}
	_, err := VonMisesQuantile(-0.1, 0, 1)
	requireErr(t, "u < 0", err, ErrProbability)
	_, err = VonMisesQuantile(0.5, 0, -1)
	requireErr(t, "κ < 0", err, ErrOutOfRange)
	_, err = VonMisesQuantile(0.5, math.Inf(1), 1)
	requireErr(t, "μ Inf", err, ErrNotFinite)
	requireNoAlloc(t, "VonMisesQuantile", func() { sinkFloat, _ = VonMisesQuantile(0.9, 0, 4) })
}

func BenchmarkVonMisesQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = VonMisesQuantile(0.9, 0, 4)
	}
}
