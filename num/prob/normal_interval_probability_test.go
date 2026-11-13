package prob

import "testing"

// NIST §1.3.6.7.1 (area 0 to 1.96 = 0.47500, so P(−1.96 ≤ Z ≤ 1.96) = 0.95000); full precision from
// 150-digit decimal arithmetic for this test: 0.95000420970355914, and P(8 ≤ Z ≤ 9) = 6.2198319858658304e-16
// (Φ(9) − Φ(8) is 0 in float64). Tolerance 5e-14 relative.
func TestNormalIntervalProbability_sourceExample(t *testing.T) {
	for _, c := range []struct{ lo, hi, want float64 }{{-1.96, 1.96, 0.95000420970355914}, {8, 9, 6.2198319858658304e-16}, {-9, -8, 6.2198319858658304e-16}} {
		if got, err := NormalIntervalProbability(c.lo, c.hi, 0, 1); err != nil || !within(got, c.want, 5e-14) {
			t.Errorf("NormalIntervalProbability(%v, %v, 0, 1) = (%v, %v), want (%v, nil)", c.lo, c.hi, got, err, c.want)
		}
	}
}

func TestNormalIntervalProbability_cases(t *testing.T) {
	for _, c := range []struct{ lo, hi, mean, sd, want float64 }{{3, 3, 0, 1, 0}, {-40, 40, 0, 1, 1}, {0, 40, 0, 1, 0.5}, {10, 12, 10, 2, 0.3413447460685429}} {
		if got, err := NormalIntervalProbability(c.lo, c.hi, c.mean, c.sd); err != nil || !within(got, c.want, 1e-15) {
			t.Errorf("NormalIntervalProbability(%v, %v, %v, %v) = (%v, %v), want %v", c.lo, c.hi, c.mean, c.sd, got, err, c.want)
		}
	}
}

func TestNormalIntervalProbability_errors(t *testing.T) {
	_, err := NormalIntervalProbability(1, 0, 0, 1)
	requireErr(t, "lo > hi", err, ErrOutOfRange)
	_, err = NormalIntervalProbability(0, 1, 0, 0)
	requireErr(t, "sd = 0", err, ErrOutOfRange)
	_, err = NormalIntervalProbability(-inf, 1, 0, 1)
	requireErr(t, "lo -Inf", err, ErrNotFinite)
	_, err = NormalIntervalProbability(0, nan, 0, 1)
	requireErr(t, "hi NaN", err, ErrNotFinite)
}

// Oracle: the interval equals NormalCDF(hi) − NormalCDF(lo) where that difference does not cancel
// (both z in [−3, 3]). Tolerance 1e-14 absolute.
func TestNormalIntervalProbability_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 500 {
		a, b := (rng.Float64()-0.5)*6, (rng.Float64()-0.5)*6
		lo, hi := min(a, b), max(a, b)
		pl, _ := NormalCDF(lo, 0, 1)
		ph, _ := NormalCDF(hi, 0, 1)
		if got, err := NormalIntervalProbability(lo, hi, 0, 1); err != nil || absDiff(got, ph-pl) > 1e-14 {
			t.Fatalf("trial %d: NormalIntervalProbability(%v, %v) = (%v, %v), want %v", trial, lo, hi, got, err, ph-pl)
		}
	}
}

func TestNormalIntervalProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NormalIntervalProbability", func() { sinkFloat, _ = NormalIntervalProbability(-1.96, 1.96, 0, 1) })
}

func BenchmarkNormalIntervalProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NormalIntervalProbability(-1.96, 1.96, 0, 1)
	}
}
