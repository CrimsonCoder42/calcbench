package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 21 §2: k = 1 → exponential (λ, λ²); k = 2, λ = 1 → (√π/2, 1 − π/4).
func TestWeibullMeanVariance_sourceExample(t *testing.T) {
	m, v, err := WeibullMeanVariance(1, 2)
	if err != nil || !within(m, 2, 1e-15) || !within(v, 4, 1e-14) {
		t.Errorf("(1, 2): (%v, %v, %v)", m, v, err)
	}
	m, v, _ = WeibullMeanVariance(2, 1)
	if !within(m, math.Sqrt(math.Pi)/2, 1e-15) || !within(v, 1-math.Pi/4, 1e-14) {
		t.Errorf("(2, 1): (%v, %v)", m, v)
	}
	_, _, err = WeibullMeanVariance(0.005, 1)
	requireErr(t, "overflow", err, ErrOverflow)
	_, _, err = WeibullMeanVariance(-1, 1)
	requireErr(t, "k < 0", err, ErrOutOfRange)
	_, _, err = WeibullMeanVariance(1, inf)
	requireErr(t, "λ Inf", err, ErrNotFinite)
}

func BenchmarkWeibullMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = WeibullMeanVariance(2, 1)
	}
}
