package prob

import (
	"math"
	"testing"
)

// Ross §8.2 / Casella–Berger §3.6 (Chebyshev's inequality): k = 2 → 1/4; k = 0.5 → the trivial bound 1.
// Hand-checked.
func TestChebyshevBound_sourceExample(t *testing.T) {
	for _, c := range []struct{ k, want float64 }{{2, 0.25}, {0.5, 1}} {
		if got, err := ChebyshevBound(c.k); err != nil || got != c.want {
			t.Errorf("ChebyshevBound(%v) = (%v, %v), want (%v, nil)", c.k, got, err, c.want)
		}
	}
}

func TestChebyshevBound_cases(t *testing.T) {
	for _, c := range []struct{ k, want float64 }{{1, 1}, {3, 1.0 / 9}, {1e-300, 1}, {math.MaxFloat64, 0}, {10, 0.01}} {
		if got, err := ChebyshevBound(c.k); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("ChebyshevBound(%v) = (%v, %v), want %v", c.k, got, err, c.want)
		}
	}
}

func TestChebyshevBound_errors(t *testing.T) {
	for _, k := range []float64{0, -1} {
		_, err := ChebyshevBound(k)
		requireErr(t, "k <= 0", err, ErrOutOfRange)
	}
	_, err := ChebyshevBound(nan)
	requireErr(t, "NaN", err, ErrNotFinite)
}

func TestChebyshevBound_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ChebyshevBound", func() { sinkFloat, _ = ChebyshevBound(2) })
}

func BenchmarkChebyshevBound(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ChebyshevBound(2)
	}
}
