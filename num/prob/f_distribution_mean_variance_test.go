package prob

import (
	"math"
	"testing"
)

// A&S 26.6.1 / JKB ch. 27 §2: d₁ = 5, d₂ = 10 → mean 10/8 = 1.25, variance 2·100·13/(5·64·6) = 2600/1920.
func TestFDistributionMeanVariance_sourceExample(t *testing.T) {
	m, v, err := FDistributionMeanVariance(5, 10)
	if err != nil || !within(m, 1.25, 1e-15) || !within(v, 2600.0/1920, 1e-14) {
		t.Errorf("(5, 10): (%v, %v, %v)", m, v, err)
	}
}

func TestFDistributionMeanVariance_matchesMoments(t *testing.T) {
	// d₂ = 30: the tail beyond 200 holds under 1e-25 of the second moment.
	f := func(x float64) float64 { d, _ := FDistributionPDF(x, 6, 30); return d }
	m1 := simpson(func(x float64) float64 { return x * f(x) }, 0, 200, 200000)
	m2 := simpson(func(x float64) float64 { return x * x * f(x) }, 0, 200, 200000)
	m, v, _ := FDistributionMeanVariance(6, 30)
	if math.Abs(m1-m) > 1e-8 || math.Abs(m2-m1*m1-v) > 1e-7 {
		t.Errorf("moments (%v, %v), want (%v, %v)", m1, m2-m1*m1, m, v)
	}
}

func TestFDistributionMeanVariance_errors(t *testing.T) {
	_, _, err := FDistributionMeanVariance(5, 4)
	requireErr(t, "d₂ ≤ 4", err, ErrOutOfRange)
	_, _, err = FDistributionMeanVariance(0, 10)
	requireErr(t, "d₁ = 0", err, ErrOutOfRange)
	_, _, err = FDistributionMeanVariance(inf, 10)
	requireErr(t, "d₁ Inf", err, ErrNotFinite)
	_, _, err = FDistributionMeanVariance(1e-310, 4.0000000001)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkFDistributionMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = FDistributionMeanVariance(5, 10)
	}
}
