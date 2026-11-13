package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 18 §1: f(x) = x^(k/2−1)e^(−x/2)/(2^(k/2)Γ(k/2)); df = 2 is exponential with mean 2.
func TestChiSquarePDF_sourceExample(t *testing.T) {
	cases := []struct{ x, df, want float64 }{
		{2, 2, 0.5 * math.Exp(-1)},
		{0, 2, 0.5},
		{0, 3, 0},
		{-1, 3, 0},
		{1, 1, math.Exp(-0.5) / math.Sqrt(2*math.Pi)},
		{3, 4, 3 * math.Exp(-1.5) / 4},
	}
	for _, c := range cases {
		if got, err := ChiSquarePDF(c.x, c.df); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("ChiSquarePDF(%v, %v) = (%v, %v), want %v", c.x, c.df, got, err, c.want)
		}
	}
	f := func(x float64) float64 { d, _ := ChiSquarePDF(x, 5); return d }
	ca, _ := ChiSquareCDF(0.5, 5)
	cb, _ := ChiSquareCDF(9, 5)
	if got := simpson(f, 0.5, 9, 2000); math.Abs(got-(cb-ca)) > 1e-11 {
		t.Errorf("∫f = %v, CDF diff = %v", got, cb-ca)
	}
}

func TestChiSquarePDF_errors(t *testing.T) {
	_, err := ChiSquarePDF(0, 1)
	requireErr(t, "x = 0, df < 2", err, ErrOutOfRange)
	_, err = ChiSquarePDF(1, 0)
	requireErr(t, "df = 0", err, ErrOutOfRange)
	_, err = ChiSquarePDF(nan, 2)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = ChiSquarePDF(1, inf)
	requireErr(t, "df Inf", err, ErrNotFinite)
	_, err = ChiSquarePDF(5e-324, 1e-3)
	requireErr(t, "overflow", err, ErrOverflow)
}

func TestChiSquareMeanVariance(t *testing.T) {
	if m, v, err := ChiSquareMeanVariance(3); err != nil || m != 3 || v != 6 {
		t.Errorf("df 3: (%v, %v, %v)", m, v, err)
	}
	// Moments against the density: ∫x f = k, ∫x² f = k² + 2k.
	f1 := func(x float64) float64 { d, _ := ChiSquarePDF(x, 4); return x * d }
	f2 := func(x float64) float64 { d, _ := ChiSquarePDF(x, 4); return x * x * d }
	if m := simpson(f1, 0, 200, 20000); math.Abs(m-4) > 1e-9 {
		t.Errorf("∫x f = %v", m)
	}
	if m2 := simpson(f2, 0, 200, 20000); math.Abs(m2-24) > 1e-8 {
		t.Errorf("∫x² f = %v", m2)
	}
	_, _, err := ChiSquareMeanVariance(0)
	requireErr(t, "df 0", err, ErrOutOfRange)
	_, _, err = ChiSquareMeanVariance(nan)
	requireErr(t, "df NaN", err, ErrNotFinite)
	_, _, err = ChiSquareMeanVariance(math.MaxFloat64)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkChiSquarePDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ChiSquarePDF(3, 4)
	}
}

func BenchmarkChiSquareMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = ChiSquareMeanVariance(4)
	}
}
