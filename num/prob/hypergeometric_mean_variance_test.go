package prob

import "testing"

// Ross §4.8 / Casella–Berger §3.2 (hypergeometric: mean nK/N, variance n(K/N)((N−K)/N)((N−n)/(N−1))):
// N = 52, K = 4, n = 5 → (5/13, 940/2873). Hand-checked with fractions; tolerance 4 ulps.
func TestHypergeometricMeanVariance_sourceExample(t *testing.T) {
	m, v, err := HypergeometricMeanVariance(52, 4, 5)
	if err != nil || !within(m, 5.0/13, 1e-15) || !within(v, 940.0/2873, 1e-15) {
		t.Errorf("HypergeometricMeanVariance(52, 4, 5) = (%v, %v, %v), want (5/13, 940/2873, nil)", m, v, err)
	}
}

func TestHypergeometricMeanVariance_cases(t *testing.T) {
	for _, c := range []struct {
		pop, succ, draws int
		m, v             float64
	}{{1, 1, 1, 1, 0}, {1, 0, 1, 0, 0}, {10, 5, 10, 5, 0}, {10, 5, 0, 0, 0}} {
		if m, v, err := HypergeometricMeanVariance(c.pop, c.succ, c.draws); err != nil || m != c.m || v != c.v {
			t.Errorf("HypergeometricMeanVariance(%d, %d, %d) = (%v, %v, %v), want (%v, %v)", c.pop, c.succ, c.draws, m, v, err, c.m, c.v)
		}
	}
}

func TestHypergeometricMeanVariance_errors(t *testing.T) {
	_, _, err := HypergeometricMeanVariance(0, 0, 0)
	requireErr(t, "population = 0", err, ErrOutOfRange)
	_, _, err = HypergeometricMeanVariance(10, 3, 11)
	requireErr(t, "draws > population", err, ErrOutOfRange)
}

func TestHypergeometricMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "HypergeometricMeanVariance", func() { sinkFloat, _, _ = HypergeometricMeanVariance(52, 4, 5) })
}

func BenchmarkHypergeometricMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = HypergeometricMeanVariance(52, 4, 5)
	}
}
