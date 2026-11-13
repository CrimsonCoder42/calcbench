package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.19 (mean λ, standard deviation √λ): λ = 2 → (2, 2). Hand-checked.
func TestPoissonMeanVariance_sourceExample(t *testing.T) {
	if m, v, err := PoissonMeanVariance(2); err != nil || m != 2 || v != 2 {
		t.Errorf("PoissonMeanVariance(2) = (%v, %v, %v), want (2, 2, nil)", m, v, err)
	}
}

func TestPoissonMeanVariance_cases(t *testing.T) {
	for _, l := range []float64{0, 1e-300, math.MaxFloat64} {
		if m, v, err := PoissonMeanVariance(l); err != nil || m != l || v != l {
			t.Errorf("PoissonMeanVariance(%v) = (%v, %v, %v)", l, m, v, err)
		}
	}
}

func TestPoissonMeanVariance_errors(t *testing.T) {
	_, _, err := PoissonMeanVariance(-1)
	requireErr(t, "lambda < 0", err, ErrOutOfRange)
	_, _, err = PoissonMeanVariance(inf)
	requireErr(t, "Inf", err, ErrNotFinite)
}

func TestPoissonMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "PoissonMeanVariance", func() { sinkFloat, _, _ = PoissonMeanVariance(2) })
}

func BenchmarkPoissonMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = PoissonMeanVariance(2)
	}
}
