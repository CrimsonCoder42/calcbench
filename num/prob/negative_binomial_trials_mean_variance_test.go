package prob

import "testing"

// Casella–Berger §3.2 (negative binomial, trials: mean r/p, variance r(1−p)/p²): r = 2, p = 0.5 → (4, 4).
// Hand-checked.
func TestNegativeBinomialTrialsMeanVariance_sourceExample(t *testing.T) {
	if m, v, err := NegativeBinomialTrialsMeanVariance(2, 0.5); err != nil || m != 4 || v != 4 {
		t.Errorf("NegativeBinomialTrialsMeanVariance(2, 0.5) = (%v, %v, %v), want (4, 4, nil)", m, v, err)
	}
}

func TestNegativeBinomialTrialsMeanVariance_cases(t *testing.T) {
	for _, c := range []struct {
		r       int
		p, m, v float64
	}{{1, 0.5, 2, 2}, {3, 1, 3, 0}, {10, 0.25, 40, 120}} {
		if m, v, err := NegativeBinomialTrialsMeanVariance(c.r, c.p); err != nil || m != c.m || v != c.v {
			t.Errorf("NegativeBinomialTrialsMeanVariance(%d, %v) = (%v, %v, %v), want (%v, %v)", c.r, c.p, m, v, err, c.m, c.v)
		}
	}
}

func TestNegativeBinomialTrialsMeanVariance_errors(t *testing.T) {
	_, _, err := NegativeBinomialTrialsMeanVariance(0, 0.5)
	requireErr(t, "r < 1", err, ErrOutOfRange)
	_, _, err = NegativeBinomialTrialsMeanVariance(2, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, _, err = NegativeBinomialTrialsMeanVariance(2, 1e-200)
	requireErr(t, "variance overflows", err, ErrOverflow)
	_, _, err = NegativeBinomialTrialsMeanVariance(2, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
}

func TestNegativeBinomialTrialsMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NegativeBinomialTrialsMeanVariance", func() { sinkFloat, _, _ = NegativeBinomialTrialsMeanVariance(2, 0.5) })
}

func BenchmarkNegativeBinomialTrialsMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = NegativeBinomialTrialsMeanVariance(2, 0.5)
	}
}
