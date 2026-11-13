package prob

import (
	"math"
	"testing"
)

// JKB vol. 2 ch. 25 §2: α = 2, β = 3 → (0.4, 6/150 = 0.04).
func TestBetaMeanVariance_sourceExample(t *testing.T) {
	if m, v, err := BetaMeanVariance(2, 3); err != nil || !within(m, 0.4, 1e-15) || !within(v, 0.04, 1e-14) {
		t.Errorf("BetaMeanVariance(2, 3) = (%v, %v, %v)", m, v, err)
	}
}

func TestBetaMeanVariance_matchesMoments(t *testing.T) {
	f := func(x float64) float64 { d, _ := BetaPDF(x, 2, 3); return d }
	m1 := simpson(func(x float64) float64 { return x * f(x) }, 0, 1, 2000)
	m2 := simpson(func(x float64) float64 { return x * x * f(x) }, 0, 1, 2000)
	m, v, _ := BetaMeanVariance(2, 3)
	if math.Abs(m1-m) > 1e-12 || math.Abs(m2-m1*m1-v) > 1e-12 {
		t.Errorf("moments (%v, %v), want (%v, %v)", m1, m2-m1*m1, m, v)
	}
}

func TestBetaMeanVariance_cases(t *testing.T) {
	if m, v, err := BetaMeanVariance(math.MaxFloat64, math.MaxFloat64); err != nil || m != 0.5 || v != 0 {
		t.Errorf("huge shapes: (%v, %v, %v), want (0.5, 0)", m, v, err)
	}
	_, _, err := BetaMeanVariance(0, 1)
	requireErr(t, "α = 0", err, ErrOutOfRange)
	_, _, err = BetaMeanVariance(1, inf)
	requireErr(t, "β Inf", err, ErrNotFinite)
}

func BenchmarkBetaMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = BetaMeanVariance(2, 3)
	}
}
