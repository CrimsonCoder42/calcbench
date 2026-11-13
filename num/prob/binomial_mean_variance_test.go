package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.18 (mean np, standard deviation √(np(1−p))): n = 10, p = 0.5 → (5, 2.5). Hand-checked.
func TestBinomialMeanVariance_sourceExample(t *testing.T) {
	if m, v, err := BinomialMeanVariance(10, 0.5); err != nil || m != 5 || v != 2.5 {
		t.Errorf("BinomialMeanVariance(10, 0.5) = (%v, %v, %v), want (5, 2.5, nil)", m, v, err)
	}
}

func TestBinomialMeanVariance_cases(t *testing.T) {
	for _, c := range []struct {
		n       int
		p, m, v float64
	}{{0, 0.5, 0, 0}, {10, 0, 0, 0}, {10, 1, 10, 0}, {math.MaxInt, 0.5, float64(math.MaxInt) / 2, float64(math.MaxInt) / 4}} {
		if m, v, err := BinomialMeanVariance(c.n, c.p); err != nil || m != c.m || v != c.v {
			t.Errorf("BinomialMeanVariance(%d, %v) = (%v, %v, %v), want (%v, %v)", c.n, c.p, m, v, err, c.m, c.v)
		}
	}
}

func TestBinomialMeanVariance_errors(t *testing.T) {
	_, _, err := BinomialMeanVariance(-1, 0.5)
	requireErr(t, "n < 0", err, ErrOutOfRange)
	_, _, err = BinomialMeanVariance(1, inf)
	requireErr(t, "Inf", err, ErrNotFinite)
	_, _, err = BinomialMeanVariance(1, 2)
	requireErr(t, "p > 1", err, ErrProbability)
}

func TestBinomialMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "BinomialMeanVariance", func() { sinkFloat, _, _ = BinomialMeanVariance(10, 0.5) })
}

func BenchmarkBinomialMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = BinomialMeanVariance(10, 0.5)
	}
}
