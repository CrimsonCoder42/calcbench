package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 21 eq. 21.3: k = 2, λ = 1, x = 1 → 2e⁻¹; k = 1 is exponential with rate 1/λ.
func TestWeibullPDF_sourceExample(t *testing.T) {
	cases := []struct{ x, k, l, want float64 }{
		{1, 2, 1, 2 * math.Exp(-1)},
		{3, 1, 2, math.Exp(-1.5) / 2},
		{0, 1, 4, 0.25},
		{0, 3, 1, 0},
		{-1, 3, 1, 0},
	}
	for _, c := range cases {
		if got, err := WeibullPDF(c.x, c.k, c.l); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("WeibullPDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.k, c.l, got, err, c.want)
		}
	}
	f := func(x float64) float64 { d, _ := WeibullPDF(x, 1.7, 2); return d }
	ca, _ := WeibullCDF(0.5, 1.7, 2)
	cb, _ := WeibullCDF(4, 1.7, 2)
	if got := simpson(f, 0.5, 4, 2000); math.Abs(got-(cb-ca)) > 1e-11 {
		t.Errorf("∫f = %v, CDF diff = %v", got, cb-ca)
	}
}

func TestWeibullPDF_errors(t *testing.T) {
	_, err := WeibullPDF(0, 0.5, 1)
	requireErr(t, "x = 0, k < 1", err, ErrOutOfRange)
	_, err = WeibullPDF(1, 0, 1)
	requireErr(t, "k = 0", err, ErrOutOfRange)
	_, err = WeibullPDF(nan, 1, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = WeibullPDF(5e-324, 0.01, 1)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkWeibullPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = WeibullPDF(1.5, 2, 1)
	}
}
