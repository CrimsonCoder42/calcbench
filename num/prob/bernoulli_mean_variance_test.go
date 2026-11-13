package prob

import "testing"

// Ross §4.6 / Casella–Berger §3.2 (Bernoulli): p = 0.3 → (0.3, 0.21). Hand-checked; tolerance 1 ulp.
func TestBernoulliMeanVariance_sourceExample(t *testing.T) {
	m, v, err := BernoulliMeanVariance(0.3)
	if err != nil || m != 0.3 || !within(v, 0.21, 2.5e-16) {
		t.Errorf("BernoulliMeanVariance(0.3) = (%v, %v, %v), want (0.3, 0.21, nil)", m, v, err)
	}
}

func TestBernoulliMeanVariance_cases(t *testing.T) {
	for _, c := range []struct{ p, m, v float64 }{{0, 0, 0}, {1, 1, 0}, {0.5, 0.5, 0.25}} {
		if m, v, err := BernoulliMeanVariance(c.p); err != nil || m != c.m || v != c.v {
			t.Errorf("BernoulliMeanVariance(%v) = (%v, %v, %v), want (%v, %v)", c.p, m, v, err, c.m, c.v)
		}
	}
}

func TestBernoulliMeanVariance_errors(t *testing.T) {
	_, _, err := BernoulliMeanVariance(nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, _, err = BernoulliMeanVariance(2)
	requireErr(t, "p > 1", err, ErrProbability)
}

func TestBernoulliMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "BernoulliMeanVariance", func() { sinkFloat, _, _ = BernoulliMeanVariance(0.3) })
}

func BenchmarkBernoulliMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = BernoulliMeanVariance(0.3)
	}
}
