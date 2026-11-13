package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.15 p. 361 (lateral 2πrh, total 2πr(r + h), volume πr²h),
// hand-derived for r 2, h 5: volume 20π, lateral 20π, total 28π.
func TestCylinder_example(t *testing.T) {
	runValCases(t, []valCase{
		{"CylinderVolume(2, 5)", func() (float64, error) { return CylinderVolume(2, 5) }, 20 * math.Pi, tolFew},
		{"CylinderLateralArea(2, 5)", func() (float64, error) { return CylinderLateralArea(2, 5) }, 20 * math.Pi, tolFew},
		{"CylinderSurfaceArea(2, 5)", func() (float64, error) { return CylinderSurfaceArea(2, 5) }, 28 * math.Pi, tolFew},
	})
}

func TestCylinder_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"volume, wide and thin, no spurious overflow", func() (float64, error) { return CylinderVolume(1e200, 1e-300) }, math.Pi * 1e100, tolSome},
		{"lateral, no spurious overflow", func() (float64, error) { return CylinderLateralArea(1e300, 1e-300) }, 2 * math.Pi, tolSome},
	})
}

func TestCylinder_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(CylinderVolume(nan, 1)) }, ErrNotFinite},
		{"volume zero height", func() error { return e1(CylinderVolume(1, 0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(CylinderVolume(1e200, 1)) }, ErrOverflow},
		{"lateral Inf", func() error { return e1(CylinderLateralArea(inf, 1)) }, ErrNotFinite},
		{"lateral negative", func() error { return e1(CylinderLateralArea(-1, 1)) }, ErrNonPositive},
		{"lateral overflow", func() error { return e1(CylinderLateralArea(1e200, 1e200)) }, ErrOverflow},
		{"surface NaN", func() error { return e1(CylinderSurfaceArea(1, nan)) }, ErrNotFinite},
		{"surface zero", func() error { return e1(CylinderSurfaceArea(0, 1)) }, ErrNonPositive},
		{"surface overflow", func() error { return e1(CylinderSurfaceArea(1e200, 1)) }, ErrOverflow},
	})
}

// Oracle: the total area is the lateral area plus two discs.
func TestCylinder_surfaceIsLateralPlusEnds(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r, h := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		lateral, _ := CylinderLateralArea(r, h)
		disc, _ := CircleArea(r)
		if got, err := CylinderSurfaceArea(r, h); err != nil || !relNear(got, lateral+2*disc, tolSome) {
			t.Fatalf("CylinderSurfaceArea(%v, %v) = (%v, %v), want %v", r, h, got, err, lateral+2*disc)
		}
	}
}

func TestCylinder_allocs(t *testing.T) {
	assertNoAllocs(t, "CylinderVolume", func() { sink, _ = CylinderVolume(2, 5) })
	assertNoAllocs(t, "CylinderLateralArea", func() { sink, _ = CylinderLateralArea(2, 5) })
	assertNoAllocs(t, "CylinderSurfaceArea", func() { sink, _ = CylinderSurfaceArea(2, 5) })
}

func BenchmarkCylinderVolume(b *testing.B) {
	b.ReportAllocs()
	r, h := 2.0, 5.0
	for b.Loop() {
		sink, _ = CylinderVolume(r, h)
	}
}

func BenchmarkCylinderLateralArea(b *testing.B) {
	b.ReportAllocs()
	r, h := 2.0, 5.0
	for b.Loop() {
		sink, _ = CylinderLateralArea(r, h)
	}
}

func BenchmarkCylinderSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	r, h := 2.0, 5.0
	for b.Loop() {
		sink, _ = CylinderSurfaceArea(r, h)
	}
}
