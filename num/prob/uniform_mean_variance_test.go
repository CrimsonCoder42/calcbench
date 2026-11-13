package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.2 (mean (A + B)/2, standard deviation √((B − A)²/12)): a = 0, b = 10 → (5, 100/12).
// Hand-checked; tolerance 1 ulp.
func TestUniformMeanVariance_sourceExample(t *testing.T) {
	m, v, err := UniformMeanVariance(0, 10)
	if err != nil || m != 5 || !within(v, 100.0/12, 2.5e-16) {
		t.Errorf("UniformMeanVariance(0, 10) = (%v, %v, %v), want (5, 100/12, nil)", m, v, err)
	}
}

func TestUniformMeanVariance_cases(t *testing.T) {
	for _, c := range []struct{ a, b, m, v float64 }{{-1, 1, 0, 1.0 / 3}, {0, 12, 6, 12}, {math.MaxFloat64 / 2, math.MaxFloat64, 0.75 * math.MaxFloat64, inf}} {
		m, v, err := UniformMeanVariance(c.a, c.b)
		if math.IsInf(c.v, 1) {
			requireErr(t, "variance overflows", err, ErrOverflow)
			continue
		}
		if err != nil || !within(m, c.m, 5e-16) || !within(v, c.v, 5e-16) {
			t.Errorf("UniformMeanVariance(%v, %v) = (%v, %v, %v), want (%v, %v)", c.a, c.b, m, v, err, c.m, c.v)
		}
	}
}

func TestUniformMeanVariance_errors(t *testing.T) {
	_, _, err := UniformMeanVariance(1, 1)
	requireErr(t, "a = b", err, ErrOutOfRange)
	_, _, err = UniformMeanVariance(nan, 1)
	requireErr(t, "NaN", err, ErrNotFinite)
	// b − a overflows, and so does the variance (2e308)²/12: not representable, so ErrOverflow stays.
	_, _, err = UniformMeanVariance(-1e308, 1e308)
	requireErr(t, "variance of [−1e308, 1e308] overflows", err, ErrOverflow)
	_, _, err = UniformMeanVariance(-math.MaxFloat64, math.MaxFloat64)
	requireErr(t, "b − a overflows", err, ErrOverflow)
}

func TestUniformMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "UniformMeanVariance", func() { sinkFloat, _, _ = UniformMeanVariance(0, 10) })
}

func BenchmarkUniformMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = UniformMeanVariance(0, 10)
	}
}
