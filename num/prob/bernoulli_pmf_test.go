package prob

import "testing"

// Ross §4.6 / Casella–Berger §3.2 (Bernoulli): p = 0.3: k = 1 → 0.3, k = 0 → 0.7, k = 2 → 0. Hand-checked.
func TestBernoulliPMF_sourceExample(t *testing.T) {
	for _, c := range []struct {
		k    int
		want float64
	}{{1, 0.3}, {0, 0.7}, {2, 0}, {-1, 0}} {
		if got, err := BernoulliPMF(c.k, 0.3); err != nil || got != c.want {
			t.Errorf("BernoulliPMF(%d, 0.3) = (%v, %v), want (%v, nil)", c.k, got, err, c.want)
		}
	}
}

func TestBernoulliPMF_cases(t *testing.T) {
	for _, c := range []struct {
		k       int
		p, want float64
	}{{0, 0, 1}, {1, 0, 0}, {1, 1, 1}, {0, 1, 0}} {
		if got, err := BernoulliPMF(c.k, c.p); err != nil || got != c.want {
			t.Errorf("BernoulliPMF(%d, %v) = (%v, %v), want %v", c.k, c.p, got, err, c.want)
		}
	}
}

func TestBernoulliPMF_errors(t *testing.T) {
	_, err := BernoulliPMF(1, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = BernoulliPMF(1, 1.5)
	requireErr(t, "p > 1", err, ErrProbability)
}

func TestBernoulliPMF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "BernoulliPMF", func() { sinkFloat, _ = BernoulliPMF(1, 0.3) })
}

func BenchmarkBernoulliPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BernoulliPMF(1, 0.3)
	}
}
