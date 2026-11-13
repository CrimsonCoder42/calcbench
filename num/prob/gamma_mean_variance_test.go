package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 17 §3: mean kθ, variance kθ². k = 3, θ = 2 → (6, 12).
func TestGammaMeanVariance_sourceExample(t *testing.T) {
	if m, v, err := GammaMeanVariance(3, 2); err != nil || m != 6 || v != 12 {
		t.Errorf("GammaMeanVariance(3, 2) = (%v, %v, %v), want (6, 12)", m, v, err)
	}
}

// The moments of the density agree (Simpson on [0, 120], 24,000 panels; tail beyond is below 1e-20).
func TestGammaMeanVariance_matchesMoments(t *testing.T) {
	f := func(x float64) float64 { d, _ := GammaPDF(x, 3, 2); return d }
	m1 := simpson(func(x float64) float64 { return x * f(x) }, 0, 120, 24000)
	m2 := simpson(func(x float64) float64 { return x * x * f(x) }, 0, 120, 24000)
	m, v, _ := GammaMeanVariance(3, 2)
	if math.Abs(m1-m) > 1e-8 || math.Abs(m2-m1*m1-v) > 1e-7 {
		t.Errorf("moments (%v, %v), want (%v, %v)", m1, m2-m1*m1, m, v)
	}
}

func TestGammaMeanVariance_errors(t *testing.T) {
	_, _, err := GammaMeanVariance(0, 1)
	requireErr(t, "shape 0", err, ErrOutOfRange)
	_, _, err = GammaMeanVariance(1, nan)
	requireErr(t, "scale NaN", err, ErrNotFinite)
	_, _, err = GammaMeanVariance(1e300, 1e10)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkGammaMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = GammaMeanVariance(3, 2)
	}
}
