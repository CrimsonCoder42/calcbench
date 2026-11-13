package prob

import "testing"

// NIST §1.3.6.7.1 (Φ(1.96) = 0.97500, so Q(1.96) = 0.02500); full precision from 150-digit decimal
// arithmetic for this test: Q(1.96) = 0.024997895148220435, Q(10) = 7.6198530241605255e-24 (1 − Φ(10)
// is 0 in float64). Tolerance 5e-14 relative (erfc's rounded argument costs about 2z² ulps).
func TestNormalSurvival_sourceExample(t *testing.T) {
	for _, c := range []struct{ x, want float64 }{{1.96, 0.024997895148220435}, {10, 7.6198530241605255e-24}} {
		if got, err := NormalSurvival(c.x, 0, 1); err != nil || !within(got, c.want, 5e-14) {
			t.Errorf("NormalSurvival(%v, 0, 1) = (%v, %v), want (%v, nil)", c.x, got, err, c.want)
		}
	}
}

func TestNormalSurvival_cases(t *testing.T) {
	for _, c := range []struct{ x, mean, sd, want float64 }{{0, 0, 1, 0.5}, {40, 0, 1, 0}, {-40, 0, 1, 1}, {8, 0, 1, 6.2209605742717839e-16}} {
		if got, err := NormalSurvival(c.x, c.mean, c.sd); err != nil || !within(got, c.want, 5e-14) {
			t.Errorf("NormalSurvival(%v, %v, %v) = (%v, %v), want %v", c.x, c.mean, c.sd, got, err, c.want)
		}
	}
}

func TestNormalSurvival_errors(t *testing.T) {
	_, err := NormalSurvival(0, 0, -1)
	requireErr(t, "sd < 0", err, ErrOutOfRange)
	_, err = NormalSurvival(0, nan, 1)
	requireErr(t, "NaN", err, ErrNotFinite)
}

func TestNormalSurvival_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NormalSurvival", func() { sinkFloat, _ = NormalSurvival(1.96, 0, 1) })
}

func BenchmarkNormalSurvival(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NormalSurvival(1.96, 0, 1)
	}
}
