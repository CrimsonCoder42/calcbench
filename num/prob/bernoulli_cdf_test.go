package prob

import (
	"math"
	"testing"
)

// Ross §4.6 / Casella–Berger §3.2 (Bernoulli): p = 0.3, k = 0 → 0.7. Hand-checked.
func TestBernoulliCDF_sourceExample(t *testing.T) {
	if got, err := BernoulliCDF(0, 0.3); err != nil || got != 0.7 {
		t.Errorf("BernoulliCDF(0, 0.3) = (%v, %v), want (0.7, nil)", got, err)
	}
}

func TestBernoulliCDF_cases(t *testing.T) {
	for _, c := range []struct {
		k    int
		want float64
	}{{-1, 0}, {math.MinInt, 0}, {1, 1}, {math.MaxInt, 1}} {
		if got, err := BernoulliCDF(c.k, 0.3); err != nil || got != c.want {
			t.Errorf("BernoulliCDF(%d, 0.3) = (%v, %v), want %v", c.k, got, err, c.want)
		}
	}
}

func TestBernoulliCDF_errors(t *testing.T) {
	_, err := BernoulliCDF(0, inf)
	requireErr(t, "Inf", err, ErrNotFinite)
	_, err = BernoulliCDF(0, -0.1)
	requireErr(t, "p < 0", err, ErrProbability)
}

func TestBernoulliCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "BernoulliCDF", func() { sinkFloat, _ = BernoulliCDF(0, 0.3) })
}

func BenchmarkBernoulliCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BernoulliCDF(0, 0.3)
	}
}
