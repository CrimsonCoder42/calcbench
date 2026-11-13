package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.16 p. 362 (ℓ = √(r² + h²), lateral πrℓ, total πr(ℓ + r), volume
// ⅓πr²h), hand-derived for r 3, h 4 (ℓ = 5): volume 12π, lateral 15π, total 24π.
func TestCone_example(t *testing.T) {
	runValCases(t, []valCase{
		{"ConeVolume(3, 4)", func() (float64, error) { return ConeVolume(3, 4) }, 12 * math.Pi, tolFew},
		{"ConeLateralArea(3, 4)", func() (float64, error) { return ConeLateralArea(3, 4) }, 15 * math.Pi, tolFew},
		{"ConeSurfaceArea(3, 4)", func() (float64, error) { return ConeSurfaceArea(3, 4) }, 24 * math.Pi, tolFew},
	})
}

func TestCone_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"lateral, slant height near the float64 limit", func() (float64, error) { return ConeLateralArea(1e-300, 1.5e308) }, math.Pi * 1.5e8, tolSome},
		{"flat cone: lateral area tends to the disc", func() (float64, error) { return ConeLateralArea(1, 1e-10) }, math.Pi, tolFew},
	})
}

func TestCone_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(ConeVolume(nan, 1)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(ConeVolume(1, 0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(ConeVolume(1e200, 1)) }, ErrOverflow},
		{"lateral Inf", func() error { return e1(ConeLateralArea(1, inf)) }, ErrNotFinite},
		{"lateral negative", func() error { return e1(ConeLateralArea(-1, 1)) }, ErrNonPositive},
		{"lateral overflow", func() error { return e1(ConeLateralArea(1e200, 1e200)) }, ErrOverflow},
		{"surface NaN", func() error { return e1(ConeSurfaceArea(1, nan)) }, ErrNotFinite},
		{"surface zero", func() error { return e1(ConeSurfaceArea(0, 1)) }, ErrNonPositive},
		{"surface overflow", func() error { return e1(ConeSurfaceArea(1e200, 1)) }, ErrOverflow},
	})
}

// Oracles: the cone is the pyramid over a disc; its total area is the lateral area plus the base; the
// lateral area matches πr√(r² + h²) computed directly.
func TestCone_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r, h := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		disc, _ := CircleArea(r)
		pyramid, _ := PyramidVolume(disc, h)
		if got, err := ConeVolume(r, h); err != nil || !relNear(got, pyramid, tolSome) {
			t.Fatalf("ConeVolume(%v, %v) = (%v, %v), pyramid %v", r, h, got, err, pyramid)
		}
		lateral, _ := ConeLateralArea(r, h)
		if want := math.Pi * r * math.Sqrt(r*r+h*h); !relNear(lateral, want, tolSome) {
			t.Fatalf("ConeLateralArea(%v, %v) = %v, want %v", r, h, lateral, want)
		}
		if got, err := ConeSurfaceArea(r, h); err != nil || !relNear(got, lateral+disc, tolSome) {
			t.Fatalf("ConeSurfaceArea(%v, %v) = (%v, %v), want %v", r, h, got, err, lateral+disc)
		}
	}
}

func TestCone_allocs(t *testing.T) {
	assertNoAllocs(t, "ConeVolume", func() { sink, _ = ConeVolume(3, 4) })
	assertNoAllocs(t, "ConeLateralArea", func() { sink, _ = ConeLateralArea(3, 4) })
	assertNoAllocs(t, "ConeSurfaceArea", func() { sink, _ = ConeSurfaceArea(3, 4) })
}

func BenchmarkConeVolume(b *testing.B) {
	b.ReportAllocs()
	r, h := 3.0, 4.0
	for b.Loop() {
		sink, _ = ConeVolume(r, h)
	}
}

func BenchmarkConeLateralArea(b *testing.B) {
	b.ReportAllocs()
	r, h := 3.0, 4.0
	for b.Loop() {
		sink, _ = ConeLateralArea(r, h)
	}
}

func BenchmarkConeSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	r, h := 3.0, 4.0
	for b.Loop() {
		sink, _ = ConeSurfaceArea(r, h)
	}
}
