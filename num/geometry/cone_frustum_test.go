package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.16 p. 363 (frustum: ℓ = √((r₁ − r₂)² + h²), total area
// π(r₁² + r₂² + (r₁ + r₂)ℓ), volume ⅓πh(r₁² + r₂² + r₁r₂)), hand-derived for 4, 1, 4 (ℓ = 5):
// volume ⅓π·4·21 = 28π; area π(16 + 1 + 25) = 42π.
func TestConeFrustum_example(t *testing.T) {
	runValCases(t, []valCase{
		{"ConeFrustumVolume(4, 1, 4)", func() (float64, error) { return ConeFrustumVolume(4, 1, 4) }, 28 * math.Pi, tolFew},
		{"ConeFrustumSurfaceArea(4, 1, 4)", func() (float64, error) { return ConeFrustumSurfaceArea(4, 1, 4) }, 42 * math.Pi, tolFew},
	})
}

func TestConeFrustum_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"radii in either order", func() (float64, error) { return ConeFrustumVolume(1, 4, 4) }, 28 * math.Pi, tolFew},
		{"equal radii is a cylinder", func() (float64, error) { return ConeFrustumVolume(2, 2, 5) }, 20 * math.Pi, tolFew},
		{"equal radii area is the cylinder's", func() (float64, error) { return ConeFrustumSurfaceArea(2, 2, 5) }, 28 * math.Pi, tolFew},
		{"wide and thin, r² would overflow", func() (float64, error) { return ConeFrustumVolume(1e200, 1e200, 1e-300) }, math.Pi * 1e100, tolSome},
	})
}

func TestConeFrustum_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(ConeFrustumVolume(nan, 1, 1)) }, ErrNotFinite},
		{"volume zero radius", func() error { return e1(ConeFrustumVolume(1, 0, 1)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(ConeFrustumVolume(1e200, 1, 1)) }, ErrOverflow},
		{"surface Inf", func() error { return e1(ConeFrustumSurfaceArea(1, 1, inf)) }, ErrNotFinite},
		{"surface negative", func() error { return e1(ConeFrustumSurfaceArea(-1, 1, 1)) }, ErrNonPositive},
		{"surface overflow", func() error { return e1(ConeFrustumSurfaceArea(1e200, 1, 1)) }, ErrOverflow},
	})
}

// Oracles: the frustum is the large cone less the small cone above it (volume and lateral area), and the
// CRC's total-area formula evaluated directly.
func TestConeFrustum_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		big, h := 1+rng.Float64()*10, 0.1+rng.Float64()*10
		small := big * (0.05 + 0.85*rng.Float64())
		apex := h * big / (big - small) // height of the full cone
		vBig, _ := ConeVolume(big, apex)
		vSmall, _ := ConeVolume(small, apex-h)
		if got, err := ConeFrustumVolume(big, small, h); err != nil || !relNear(got, vBig-vSmall, 1e-12) {
			t.Fatalf("ConeFrustumVolume(%v, %v, %v) = (%v, %v), want %v", big, small, h, got, err, vBig-vSmall)
		}
		l := math.Sqrt((big-small)*(big-small) + h*h)
		want := math.Pi * (big*big + small*small + (big+small)*l)
		if got, err := ConeFrustumSurfaceArea(big, small, h); err != nil || !relNear(got, want, 1e-14) {
			t.Fatalf("ConeFrustumSurfaceArea(%v, %v, %v) = (%v, %v), want %v", big, small, h, got, err, want)
		}
	}
}

func TestConeFrustum_allocs(t *testing.T) {
	assertNoAllocs(t, "ConeFrustumVolume", func() { sink, _ = ConeFrustumVolume(4, 1, 4) })
	assertNoAllocs(t, "ConeFrustumSurfaceArea", func() { sink, _ = ConeFrustumSurfaceArea(4, 1, 4) })
}

func BenchmarkConeFrustumVolume(b *testing.B) {
	b.ReportAllocs()
	r1, r2, h := 4.0, 1.0, 4.0
	for b.Loop() {
		sink, _ = ConeFrustumVolume(r1, r2, h)
	}
}

func BenchmarkConeFrustumSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	r1, r2, h := 4.0, 1.0, 4.0
	for b.Loop() {
		sink, _ = ConeFrustumSurfaceArea(r1, r2, h)
	}
}
