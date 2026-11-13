package prob

import "testing"

// NIST §1.3.6.6.7 (mean β, standard deviation β, β = 1/λ): λ = 2 → (0.5, 0.25). Hand-checked.
func TestExponentialMeanVariance_sourceExample(t *testing.T) {
	if m, v, err := ExponentialMeanVariance(2); err != nil || m != 0.5 || v != 0.25 {
		t.Errorf("ExponentialMeanVariance(2) = (%v, %v, %v), want (0.5, 0.25, nil)", m, v, err)
	}
}

func TestExponentialMeanVariance_cases(t *testing.T) {
	for _, c := range []struct{ rate, m, v float64 }{{1, 1, 1}, {0.5, 2, 4}, {1e100, 1e-100, 1e-200}} {
		if m, v, err := ExponentialMeanVariance(c.rate); err != nil || !within(m, c.m, 5e-16) || !within(v, c.v, 5e-16) {
			t.Errorf("ExponentialMeanVariance(%v) = (%v, %v, %v), want (%v, %v)", c.rate, m, v, err, c.m, c.v)
		}
	}
}

func TestExponentialMeanVariance_errors(t *testing.T) {
	_, _, err := ExponentialMeanVariance(0)
	requireErr(t, "rate = 0", err, ErrOutOfRange)
	_, _, err = ExponentialMeanVariance(1e-200)
	requireErr(t, "variance overflows", err, ErrOverflow)
	_, _, err = ExponentialMeanVariance(nan)
	requireErr(t, "NaN", err, ErrNotFinite)
}

func TestExponentialMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ExponentialMeanVariance", func() { sinkFloat, _, _ = ExponentialMeanVariance(2) })
}

func BenchmarkExponentialMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = ExponentialMeanVariance(2)
	}
}
