package prob

import (
	"math"
	"testing"
)

// d₁ = d₂ = 2: F(x) = x/(1 + x). F(1, ν) is t² (JKB ch. 27 §1): P(F ≤ t²) = 2·P(T ≤ t) − 1.
func TestFDistributionCDF_sourceExample(t *testing.T) {
	if got, err := FDistributionCDF(1, 2, 2); err != nil || !within(got, 0.5, 1e-14) {
		t.Errorf("FDistributionCDF(1, 2, 2) = (%v, %v)", got, err)
	}
	for _, x := range []float64{0.5, 1.7, 4} {
		f, _ := FDistributionCDF(x*x, 1, 6)
		c, _ := StudentTCDF(x, 6)
		if !within(f, 2*c-1, 1e-12) {
			t.Errorf("t² relation x %v: %v vs %v", x, f, 2*c-1)
		}
	}
}

func TestFDistributionCDF_cases(t *testing.T) {
	if got, _ := FDistributionCDF(-3, 2, 2); got != 0 {
		t.Errorf("x < 0: %v", got)
	}
	if got, _ := FDistributionCDF(math.MaxFloat64, 2, 1e-300); got != 1 {
		t.Errorf("r overflow: %v", got)
	}
	_, err := FDistributionCDF(1, 2, 0)
	requireErr(t, "d₂ = 0", err, ErrOutOfRange)
	_, err = FDistributionCDF(inf, 2, 2)
	requireErr(t, "x Inf", err, ErrNotFinite)
}

func BenchmarkFDistributionCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = FDistributionCDF(2.5, 5, 10)
	}
}
