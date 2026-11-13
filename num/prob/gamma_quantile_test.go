package prob

import (
	"math"
	"testing"
)

// k = 1 is the exponential (JKB ch. 17 §1, ch. 19): Q(u) = −θ ln(1 − u); u = 0.5, θ = 1 → ln 2.
func TestGammaQuantile_sourceExample(t *testing.T) {
	if got, err := GammaQuantile(0.5, 1, 1); err != nil || !within(got, math.Ln2, 1e-14) {
		t.Errorf("GammaQuantile(0.5, 1, 1) = (%v, %v), want ln 2", got, err)
	}
	// Chi-square critical value via Gamma(ν/2, 2): NIST e-Handbook §1.3.6.7.4, ν = 3, 0.95 → 7.815.
	if got, _ := GammaQuantile(0.95, 1.5, 2); math.Abs(got-7.815) > 5e-4 {
		t.Errorf("chi-square 3 df 0.95 = %v, want 7.815", got)
	}
}

// CDF(Quantile(u)) = u over shapes and tails; tolerance 1e-12 relative to min(u, 1 − u).
func TestGammaQuantile_roundTrip(t *testing.T) {
	for _, k := range []float64{0.3, 1, 2.5, 40} {
		for _, u := range []float64{1e-10, 0.01, 0.3, 0.5, 0.9, 1 - 1e-9} {
			x, err := GammaQuantile(u, k, 1.5)
			if err != nil {
				t.Fatalf("k %v u %v: %v", k, u, err)
			}
			c, s, _ := gammaCDFSurvival(x, k, 1.5)
			if u <= 0.5 && !within(c, u, 1e-12) || u > 0.5 && !within(s, 1-u, 1e-11) {
				t.Errorf("k %v u %v: x %v, cdf %v, surv %v", k, u, x, c, s)
			}
		}
	}
}

func TestGammaQuantile_errors(t *testing.T) {
	if got, err := GammaQuantile(0, 2, 1); err != nil || got != 0 {
		t.Errorf("u = 0: (%v, %v)", got, err)
	}
	_, err := GammaQuantile(1, 2, 1)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = GammaQuantile(1.5, 2, 1)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = GammaQuantile(0.5, 0, 1)
	requireErr(t, "shape 0", err, ErrOutOfRange)
	_, err = GammaQuantile(nan, 1, 1)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = GammaQuantile(0.9, 1, 1e308)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkGammaQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GammaQuantile(0.9, 3, 1)
	}
}
