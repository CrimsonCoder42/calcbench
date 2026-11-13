package prob

import (
	"math"
	"testing"
)

// Casella–Berger §3.2 (discrete uniform on 1…N: mean (N+1)/2, variance (N+1)(N−1)/12): a fair die →
// (3.5, 35/12 ≈ 2.9166666666666665). Hand-checked; tolerance 1 ulp.
func TestDiscreteUniformMeanVariance_sourceExample(t *testing.T) {
	m, v, err := DiscreteUniformMeanVariance(1, 6)
	if err != nil || m != 3.5 || !within(v, 35.0/12, 2.5e-16) {
		t.Errorf("DiscreteUniformMeanVariance(1, 6) = (%v, %v, %v), want (3.5, 35/12, nil)", m, v, err)
	}
}

func TestDiscreteUniformMeanVariance_cases(t *testing.T) {
	for _, c := range []struct {
		a, b int
		m, v float64
	}{{5, 5, 5, 0}, {0, 1, 0.5, 0.25}, {-3, 3, 0, 4}, {math.MaxInt - 1, math.MaxInt, float64(math.MaxInt), 0.25}} {
		if m, v, err := DiscreteUniformMeanVariance(c.a, c.b); err != nil || !within(m, c.m, 1e-15) || !within(v, c.v, 1e-15) {
			t.Errorf("DiscreteUniformMeanVariance(%d, %d) = (%v, %v, %v), want (%v, %v)", c.a, c.b, m, v, err, c.m, c.v)
		}
	}
}

func TestDiscreteUniformMeanVariance_errors(t *testing.T) {
	_, _, err := DiscreteUniformMeanVariance(2, 1)
	requireErr(t, "a > b", err, ErrOutOfRange)
}

// The full int range: mean (MinInt + MaxInt)/2 = −½; variance (m² − 1)/12 with m = 2⁶⁴ → 2¹²⁸/12 in float64
// (the −1 is below half an ulp). Hand-checked. Before the fix this was ErrOverflow.
func TestDiscreteUniformMeanVariance_fullIntRange(t *testing.T) {
	m, v, err := DiscreteUniformMeanVariance(math.MinInt, math.MaxInt)
	if err != nil || m != -0.5 || !within(v, 0x1p128/12, 2.3e-16) {
		t.Errorf("DiscreteUniformMeanVariance(MinInt, MaxInt) = (%v, %v, %v), want (−0.5, 2¹²⁸/12, nil)", m, v, err)
	}
}

func TestDiscreteUniformMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "DiscreteUniformMeanVariance", func() { sinkFloat, _, _ = DiscreteUniformMeanVariance(1, 6) })
}

func BenchmarkDiscreteUniformMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = DiscreteUniformMeanVariance(1, 6)
	}
}
