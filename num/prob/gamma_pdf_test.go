package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 17 eq. 17.1: f(x) = x^(k−1)e^(−x/θ)/(Γ(k)θ^k). Hand-computed: k = 2, θ = 1, x = 1 → e⁻¹;
// k = 3, θ = 2, x = 2 → 4e⁻¹/(2·8) = e⁻¹/4.
func TestGammaPDF_sourceExample(t *testing.T) {
	cases := []struct{ x, k, th, want float64 }{
		{1, 2, 1, math.Exp(-1)},
		{2, 3, 2, math.Exp(-1) / 4},
		{3, 1, 2, math.Exp(-1.5) / 2}, // the exponential case
	}
	for _, c := range cases {
		if got, err := GammaPDF(c.x, c.k, c.th); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("GammaPDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.k, c.th, got, err, c.want)
		}
	}
}

func TestGammaPDF_cases(t *testing.T) {
	cases := []struct {
		name         string
		x, k, th, wt float64
	}{
		{"negative x", -1, 2, 1, 0},
		{"x = 0, k > 1", 0, 2, 1, 0},
		{"x = 0, k = 1", 0, 1, 4, 0.25},
		{"far tail underflows to 0", 1e6, 2, 1, 0},
	}
	for _, c := range cases {
		if got, err := GammaPDF(c.x, c.k, c.th); err != nil || got != c.wt {
			t.Errorf("%s: (%v, %v), want %v", c.name, got, err, c.wt)
		}
	}
}

// The density integrates to the CDF difference (Simpson, 2,000 panels; tolerance 1e-9 covers the rule's
// h⁴ error on smooth densities).
func TestGammaPDF_integratesToCDF(t *testing.T) {
	for _, p := range [][2]float64{{2, 1}, {3.5, 0.7}, {10, 2}} {
		f := func(x float64) float64 { d, _ := GammaPDF(x, p[0], p[1]); return d }
		a, b := 0.5*p[0]*p[1], 2*p[0]*p[1]
		ca, _ := GammaCDF(a, p[0], p[1])
		cb, _ := GammaCDF(b, p[0], p[1])
		if got := simpson(f, a, b, 2000); math.Abs(got-(cb-ca)) > 1e-9 {
			t.Errorf("shape %v scale %v: ∫f = %v, CDF diff = %v", p[0], p[1], got, cb-ca)
		}
	}
}

func TestGammaPDF_errors(t *testing.T) {
	_, err := GammaPDF(0, 0.5, 1)
	requireErr(t, "x = 0, k < 1", err, ErrOutOfRange)
	_, err = GammaPDF(1, 0, 1)
	requireErr(t, "shape 0", err, ErrOutOfRange)
	_, err = GammaPDF(1, 1, -1)
	requireErr(t, "scale < 0", err, ErrOutOfRange)
	_, err = GammaPDF(nan, 1, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = GammaPDF(1, inf, 1)
	requireErr(t, "shape Inf", err, ErrNotFinite)
	_, err = GammaPDF(5e-324, 0.01, 1)
	requireErr(t, "density overflow", err, ErrOverflow)
}

func TestGammaPDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "GammaPDF", func() { sinkFloat, _ = GammaPDF(2, 3, 1) })
}

func BenchmarkGammaPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GammaPDF(2, 3, 1)
	}
}
