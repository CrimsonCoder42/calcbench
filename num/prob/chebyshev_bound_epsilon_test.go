package prob

import "testing"

// Ross §8.2 (the weekly-production illustration of Chebyshev's inequality, widely reproduced from Ross;
// the example number is not confirmed in print): variance 25, ε = 10 → 25/100 = 0.25, so
// P(40 < X < 60) ≥ 0.75 for a mean-50 variable. Hand-checked.
func TestChebyshevBoundEpsilon_sourceExample(t *testing.T) {
	if got, err := ChebyshevBoundEpsilon(25, 10); err != nil || got != 0.25 {
		t.Errorf("ChebyshevBoundEpsilon(25, 10) = (%v, %v), want (0.25, nil)", got, err)
	}
}

func TestChebyshevBoundEpsilon_cases(t *testing.T) {
	for _, c := range []struct{ v, e, want float64 }{{0, 1, 0}, {4, 1, 1}, {1, 2, 0.25}, {1e-300, 1e-160, 1}, {1e300, 1e200, 1e-100}} {
		if got, err := ChebyshevBoundEpsilon(c.v, c.e); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("ChebyshevBoundEpsilon(%v, %v) = (%v, %v), want %v", c.v, c.e, got, err, c.want)
		}
	}
}

func TestChebyshevBoundEpsilon_errors(t *testing.T) {
	_, err := ChebyshevBoundEpsilon(-1, 1)
	requireErr(t, "variance < 0", err, ErrOutOfRange)
	_, err = ChebyshevBoundEpsilon(1, 0)
	requireErr(t, "epsilon = 0", err, ErrOutOfRange)
	_, err = ChebyshevBoundEpsilon(inf, 1)
	requireErr(t, "Inf", err, ErrNotFinite)
	_, err = ChebyshevBoundEpsilon(1, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
}

func TestChebyshevBoundEpsilon_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ChebyshevBoundEpsilon", func() { sinkFloat, _ = ChebyshevBoundEpsilon(25, 10) })
}

func BenchmarkChebyshevBoundEpsilon(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ChebyshevBoundEpsilon(25, 10)
	}
}
