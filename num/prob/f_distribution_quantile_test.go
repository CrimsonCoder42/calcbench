package prob

import (
	"math"
	"testing"
)

// d₁ = d₂ = 2: Q(u) = u/(1 − u). NIST e-Handbook §1.3.6.7.3 (α = 0.05): F(5, 10) → 3.326, F(1, 1) → 161.448.
func TestFDistributionQuantile_sourceExample(t *testing.T) {
	if got, err := FDistributionQuantile(0.5, 2, 2); err != nil || !within(got, 1, 1e-14) {
		t.Errorf("(0.5, 2, 2): (%v, %v)", got, err)
	}
	cases := []struct{ d1, d2, want float64 }{{5, 10, 3.326}, {1, 1, 161.448}, {10, 20, 2.348}}
	for _, c := range cases {
		if got, _ := FDistributionQuantile(0.95, c.d1, c.d2); math.Abs(got-c.want) > 5e-4 {
			t.Errorf("F(%v, %v): %v, want %v", c.d1, c.d2, got, c.want)
		}
	}
}

func TestFDistributionQuantile_roundTrip(t *testing.T) {
	for _, u := range []float64{1e-9, 0.2, 0.5, 0.8, 1 - 1e-9} {
		x, err := FDistributionQuantile(u, 3, 7)
		c, s, _ := fCDFSurvival(x, 3, 7)
		if err != nil || u <= 0.5 && !within(c, u, 1e-11) || u > 0.5 && !within(s, 1-u, 1e-9) {
			t.Errorf("u %v: x %v cdf %v surv %v err %v", u, x, c, s, err)
		}
	}
}

func TestFDistributionQuantile_errors(t *testing.T) {
	if got, err := FDistributionQuantile(0, 3, 7); err != nil || got != 0 {
		t.Errorf("u = 0: (%v, %v)", got, err)
	}
	_, err := FDistributionQuantile(1, 3, 7)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = FDistributionQuantile(-1, 3, 7)
	requireErr(t, "u < 0", err, ErrProbability)
	_, err = FDistributionQuantile(0.5, 3, nan)
	requireErr(t, "d₂ NaN", err, ErrNotFinite)
}

func BenchmarkFDistributionQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = FDistributionQuantile(0.95, 5, 10)
	}
}
