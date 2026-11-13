package prob

import (
	"math"
	"testing"
)

// Mardia & Jupp §3.5.4: κ = 0 is uniform; the density integrates to the CDF (Hill 1977's series); F(μ) = ½.
func TestVonMises(t *testing.T) {
	if got, err := VonMisesPDF(1, 0, 0); err != nil || !within(got, 1/(2*math.Pi), 1e-15) {
		t.Errorf("uniform PDF = (%v, %v)", got, err)
	}
	if got, _ := VonMisesCDF(0, 0, 0); !within(got, 0.5, 1e-15) {
		t.Errorf("uniform CDF(0) = %v", got)
	}
	// f(μ) = e^κ/(2π I₀(κ)); I₀(1) = 1.2660658777520082 (A&S Table 9.8).
	if got, _ := VonMisesPDF(0, 0, 1); !within(got, math.E/(2*math.Pi*1.2660658777520082), 1e-14) {
		t.Errorf("PDF(μ), κ = 1 = %v", got)
	}
	for _, kappa := range []float64{0.3, 1, 10.5, 40, 650, 750, 1e4} {
		if c, _ := VonMisesCDF(2, 2, kappa); !within(c, 0.5, 1e-14) {
			t.Errorf("κ = %v: F(μ) = %v", kappa, c)
		}
		if c, _ := VonMisesCDF(math.Pi, 0, kappa); !within(c, 1, 1e-14) {
			t.Errorf("κ = %v: F(μ + π) = %v", kappa, c)
		}
		if c, _ := VonMisesCDF(-math.Pi, 0, kappa); math.Abs(c) > 1e-14 {
			t.Errorf("κ = %v: F(μ − π) = %v", kappa, c)
		}
		f := func(x float64) float64 { v, _ := VonMisesPDF(x, 0, kappa); return v }
		w := math.Min(math.Pi, 12/math.Sqrt(kappa))
		for _, frac := range []float64{-0.6, -0.1, 0.25} {
			x := frac * w
			want := simpson(f, -math.Min(math.Pi, 20/math.Sqrt(kappa)), x, 20000)
			if c, _ := VonMisesCDF(x, 0, kappa); math.Abs(c-want) > 1e-12 {
				t.Errorf("κ = %v, x = %v: CDF %v, ∫f %v", kappa, x, c, want)
			}
		}
	}
	// The scaled I₀ branches agree at κ = 700.
	a, _ := scaledBesselI0(math.Nextafter(700, 0))
	b, _ := scaledBesselI0(700)
	if !within(a, b, 1e-12) {
		t.Errorf("I₀ switch: %v vs %v", a, b)
	}
	_, err := VonMisesPDF(nan, 0, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = VonMisesPDF(0, inf, 1)
	requireErr(t, "location Inf", err, ErrNotFinite)
	_, err = VonMisesCDF(0, 0, nan)
	requireErr(t, "kappa NaN", err, ErrNotFinite)
	_, err = VonMisesCDF(0, 0, -1)
	requireErr(t, "kappa < 0", err, ErrOutOfRange)
	_, err = VonMisesCDF(0, 0, 2e6)
	requireErr(t, "kappa > 1e6", err, ErrOutOfRange)
	_, err = VonMisesPDF(4, 0, 1)
	requireErr(t, "x beyond π", err, ErrOutOfRange)
	requireNoAlloc(t, "VonMisesCDF", func() { sinkFloat, _ = VonMisesCDF(0.5, 0, 4) })
}

func BenchmarkVonMisesPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = VonMisesPDF(0.5, 0, 4)
	}
}

func BenchmarkVonMisesCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = VonMisesCDF(0.5, 0, 4)
	}
}
