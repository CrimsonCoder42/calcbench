package prob

import (
	"math"
	"testing"
)

// I_x(2, 1) = x² (DLMF 8.17.5 family): u = 0.25 → 0.5.
func TestBetaQuantile_sourceExample(t *testing.T) {
	if got, err := BetaQuantile(0.25, 2, 1); err != nil || !within(got, 0.5, 1e-14) {
		t.Errorf("BetaQuantile(0.25, 2, 1) = (%v, %v), want 0.5", got, err)
	}
}

func TestBetaQuantile_roundTrip(t *testing.T) {
	for _, s := range [][2]float64{{0.5, 0.5}, {2, 3}, {30, 2}} {
		for _, u := range []float64{1e-8, 0.1, 0.5, 0.9, 1 - 1e-8} {
			x, err := BetaQuantile(u, s[0], s[1])
			c, sv, _ := betaCDFSurvival(x, s[0], s[1])
			if s == [2]float64{0.5, 0.5} && u > 0.5 {
				// Arcsine law: F(x) = (2/π)·asin(√x), so the true quantile is x = 1 − sin²(π(1−u)/2). At
				// u = 1 − 1e-8 that is 1 − 2.5e-16, about 2 ulps below 1: no double x has a survival within
				// 1e-9 of 1e-8 (the neighbours of 1 give 0 and ~6.7e-9), so the survival check cannot hold.
				// The documented property is checked instead: x within 4 ulps of the true quantile.
				sn := math.Sin(math.Pi * (1 - u) / 2)
				want := 1 - sn*sn
				if ulp := math.Nextafter(want, 2) - want; err != nil || math.Abs(x-want) > 4*ulp {
					t.Errorf("shapes %v u %v: x %v, want %v within 4 ulps, err %v", s, u, x, want, err)
				}
				continue
			}
			if err != nil || u <= 0.5 && !within(c, u, 1e-11) || u > 0.5 && !within(sv, 1-u, 1e-9) {
				t.Errorf("shapes %v u %v: x %v cdf %v surv %v err %v", s, u, x, c, sv, err)
			}
		}
	}
}

func TestBetaQuantile_errors(t *testing.T) {
	for _, u := range []float64{0, 1} {
		if got, err := BetaQuantile(u, 2, 3); err != nil || got != u {
			t.Errorf("u = %v: (%v, %v)", u, got, err)
		}
	}
	_, err := BetaQuantile(-0.1, 2, 3)
	requireErr(t, "u < 0", err, ErrProbability)
	_, err = BetaQuantile(0.5, 0, 3)
	requireErr(t, "α = 0", err, ErrOutOfRange)
	_, err = BetaQuantile(0.5, nan, 3)
	requireErr(t, "α NaN", err, ErrNotFinite)
}

func BenchmarkBetaQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BetaQuantile(0.3, 2, 3)
	}
}
