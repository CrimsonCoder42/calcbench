package prob

import (
	"math"
	"testing"
)

// JKB vol. 2 ch. 27 eq. 27.4 with d₁ = d₂ = 2: f(x) = (1 + x)⁻² (hand-derived); x = 1 → 0.25.
func TestFDistributionPDF_sourceExample(t *testing.T) {
	for _, x := range []float64{0.25, 1, 9} {
		if got, err := FDistributionPDF(x, 2, 2); err != nil || !within(got, 1/((1+x)*(1+x)), 1e-14) {
			t.Errorf("x %v: (%v, %v)", x, got, err)
		}
	}
}

func TestFDistributionPDF_cases(t *testing.T) {
	cases := []struct{ x, d1, d2, want float64 }{{-1, 3, 4, 0}, {0, 3, 4, 0}, {0, 2, 7, 1}}
	for _, c := range cases {
		if got, err := FDistributionPDF(c.x, c.d1, c.d2); err != nil || got != c.want {
			t.Errorf("(%v, %v, %v): (%v, %v), want %v", c.x, c.d1, c.d2, got, err, c.want)
		}
	}
	f := func(x float64) float64 { d, _ := FDistributionPDF(x, 5, 10); return d }
	ca, _ := FDistributionCDF(0.2, 5, 10)
	cb, _ := FDistributionCDF(3, 5, 10)
	if got := simpson(f, 0.2, 3, 2000); math.Abs(got-(cb-ca)) > 1e-10 {
		t.Errorf("∫f = %v, CDF diff = %v", got, cb-ca)
	}
}

func TestFDistributionPDF_errors(t *testing.T) {
	_, err := FDistributionPDF(0, 1, 3)
	requireErr(t, "x = 0, d₁ < 2", err, ErrOutOfRange)
	_, err = FDistributionPDF(1, 0, 3)
	requireErr(t, "d₁ = 0", err, ErrOutOfRange)
	_, err = FDistributionPDF(nan, 1, 3)
	requireErr(t, "x NaN", err, ErrNotFinite)
}

func BenchmarkFDistributionPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = FDistributionPDF(1.5, 5, 10)
	}
}
