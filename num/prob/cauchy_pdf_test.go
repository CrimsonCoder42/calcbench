package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 16 eq. 16.1: standard Cauchy f(0) = 1/π, f(1) = 1/(2π); equals Student t with ν = 1.
func TestCauchyPDF_sourceExample(t *testing.T) {
	cases := []struct{ x, x0, g, want float64 }{
		{0, 0, 1, 1 / math.Pi}, {1, 0, 1, 1 / (2 * math.Pi)}, {5, 3, 2, 1 / (4 * math.Pi)},
	}
	for _, c := range cases {
		if got, err := CauchyPDF(c.x, c.x0, c.g); err != nil || !within(got, c.want, 1e-15) {
			t.Errorf("CauchyPDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.x0, c.g, got, err, c.want)
		}
	}
	a, _ := CauchyPDF(2.3, 0, 1)
	b, _ := StudentTPDF(2.3, 1)
	if !within(a, b, 1e-14) {
		t.Errorf("Cauchy %v vs t₁ %v", a, b)
	}
	_, err := CauchyPDF(1, 0, 0)
	requireErr(t, "γ = 0", err, ErrOutOfRange)
	_, err = CauchyPDF(1, nan, 1)
	requireErr(t, "x₀ NaN", err, ErrNotFinite)
	_, err = CauchyPDF(0, 0, 1e-310)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkCauchyPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = CauchyPDF(1.5, 0, 1)
	}
}
