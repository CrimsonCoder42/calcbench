package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 17 with DLMF 8.4.4 (integer a): k = 1 → 1 − e^(−x/θ); k = 2, θ = 1, x = 2 → 1 − 3e⁻².
func TestGammaCDF_sourceExample(t *testing.T) {
	cases := []struct{ x, k, th, want float64 }{
		{1, 1, 1, 1 - math.Exp(-1)},
		{2, 2, 1, 1 - 3*math.Exp(-2)},
		{4, 1, 2, 1 - math.Exp(-2)},
	}
	for _, c := range cases {
		if got, err := GammaCDF(c.x, c.k, c.th); err != nil || !within(got, c.want, 1e-13) {
			t.Errorf("GammaCDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.k, c.th, got, err, c.want)
		}
	}
}

// Gamma(ν/2, 2) is chi-square with ν degrees of freedom (JKB ch. 18 §1): agree with NUM-PROB-063.
func TestGammaCDF_matchesChiSquare(t *testing.T) {
	for _, df := range []float64{1, 3, 10} {
		for _, x := range []float64{0.1, 1, 5, 20} {
			g, _ := GammaCDF(x, df/2, 2)
			c, _ := ChiSquareCDF(x, df)
			if !within(g, c, 1e-12) {
				t.Errorf("df %v x %v: gamma %v, chi-square %v", df, x, g, c)
			}
		}
	}
}

func TestGammaCDF_cases(t *testing.T) {
	for _, x := range []float64{-1, 0} {
		if got, err := GammaCDF(x, 2, 1); err != nil || got != 0 {
			t.Errorf("x = %v: (%v, %v), want 0", x, got, err)
		}
	}
	if got, _ := GammaCDF(1e6, 2, 1); got != 1 {
		t.Errorf("far upper: %v, want 1", got)
	}
}

func TestGammaCDF_errors(t *testing.T) {
	_, err := GammaCDF(1, -1, 1)
	requireErr(t, "shape < 0", err, ErrOutOfRange)
	_, err = GammaCDF(nan, 1, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = GammaCDF(1, 1, inf)
	requireErr(t, "scale Inf", err, ErrNotFinite)
}

func BenchmarkGammaCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GammaCDF(2.5, 3, 1)
	}
}
