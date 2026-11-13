package prob

import (
	"math"
	"testing"
)

// Chhikara & Folks 1989 ch. 2: f(μ) = √(λ/(2πμ³)) and F(x) = Φ(a) + e^(2λ/μ)Φ(−b); at μ = λ = x = 1,
// f = 1/√(2π) and F = Φ(0) + e²Φ(−2).
func TestWald(t *testing.T) {
	if got, err := WaldPDF(1, 1, 1); err != nil || !within(got, 1/math.Sqrt(2*math.Pi), 1e-15) {
		t.Errorf("PDF(1) = (%v, %v)", got, err)
	}
	want := 0.5 + math.Exp(2)*0.5*math.Erfc(math.Sqrt2)
	if got, err := WaldCDF(1, 1, 1); err != nil || !within(got, want, 1e-14) {
		t.Errorf("CDF(1) = (%v, %v), want %v", got, err, want)
	}
	if q, err := WaldQuantile(want, 1, 1); err != nil || !within(q, 1, 1e-13) {
		t.Errorf("Q(F(1)) = (%v, %v)", q, err)
	}
	if got, _ := WaldPDF(-1, 1, 1); got != 0 {
		t.Errorf("PDF(−1) = %v", got)
	}
	if got, _ := WaldCDF(0, 1, 1); got != 0 {
		t.Errorf("CDF(0) = %v", got)
	}
	// Large λ/μ: e^(2λ/μ) alone overflows; the CDF stays in [0, 1] and near ½ at the mean (normal limit).
	if got, err := WaldCDF(1, 1, 1e6); err != nil || math.Abs(got-0.5) > 1e-3 {
		t.Errorf("CDF large shape = (%v, %v)", got, err)
	}
	// The scaled tail's two branches agree at the switch point.
	if a, b := scaledNormalTail(math.Nextafter(26, 0)), scaledNormalTail(26); !within(a, b, 1e-12) {
		t.Errorf("tail switch: %v vs %v", a, b)
	}
	for _, sh := range []float64{3, 200} {
		checkClosedForm(t, closedForm{"wald",
			func(x float64) (float64, error) { return WaldPDF(x, 2, sh) },
			func(x float64) (float64, error) { return WaldCDF(x, 2, sh) },
			func(u float64) (float64, error) { return WaldQuantile(u, 2, sh) },
			func() (float64, float64, error) { return WaldMeanVariance(2, sh) }, 1e-12, 400})
	}
	if m, v, _ := WaldMeanVariance(2, 4); m != 2 || v != 2 {
		t.Errorf("MV = (%v, %v)", m, v)
	}
	_, err := WaldPDF(nan, 1, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = WaldCDF(1, 0, 1)
	requireErr(t, "mean 0", err, ErrOutOfRange)
	_, err = WaldCDF(1, 1, -1)
	requireErr(t, "shape < 0", err, ErrOutOfRange)
	_, err = WaldCDF(inf, 1, 1)
	requireErr(t, "x Inf", err, ErrNotFinite)
	_, err = WaldPDF(1e-300, 1e-300, 1e300)
	requireErr(t, "pdf overflow", err, ErrOverflow)
	_, err = WaldQuantile(0, 1, 1)
	requireErr(t, "u = 0", err, ErrOutOfRange)
	_, err = WaldQuantile(1.5, 1, 1)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = WaldQuantile(0.5, nan, 1)
	requireErr(t, "quantile mean NaN", err, ErrNotFinite)
	_, _, err = WaldMeanVariance(1e200, 1e-200)
	requireErr(t, "MV overflow", err, ErrOverflow)
	_, _, err = WaldMeanVariance(1, 0)
	requireErr(t, "MV shape 0", err, ErrOutOfRange)
	requireNoAlloc(t, "WaldQuantile", func() { sinkFloat, _ = WaldQuantile(0.7, 1, 1) })
}

func BenchmarkWaldPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = WaldPDF(1.5, 1, 2)
	}
}

func BenchmarkWaldCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = WaldCDF(1.5, 1, 2)
	}
}

func BenchmarkWaldQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = WaldQuantile(0.7, 1, 2)
	}
}

func BenchmarkWaldMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = WaldMeanVariance(1, 2)
	}
}
