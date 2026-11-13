package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.15 (V = Ah, lateral 2πrh) with the annulus of §4.6.6, hand-derived for
// R 3, r 2, h 5: volume π(9 − 4)·5 = 25π; area 2π(3 + 2)·5 + 2π(9 − 4) = 60π.
func TestHollowCylinder_example(t *testing.T) {
	runValCases(t, []valCase{
		{"HollowCylinderVolume(3, 2, 5)", func() (float64, error) { return HollowCylinderVolume(3, 2, 5) }, 25 * math.Pi, tolFew},
		{"HollowCylinderSurfaceArea(3, 2, 5)", func() (float64, error) { return HollowCylinderSurfaceArea(3, 2, 5) }, 60 * math.Pi, tolFew},
	})
}

func TestHollowCylinder_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"thin wall keeps its digits", func() (float64, error) { return HollowCylinderVolume(1+0x1p-40, 1, 1) }, math.Pi * 0x1p-40 * (2 + 0x1p-40), tolFew},
	})
}

func TestHollowCylinder_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(HollowCylinderVolume(nan, 1, 1)) }, ErrNotFinite},
		{"volume Inf height", func() error { return e1(HollowCylinderVolume(2, 1, inf)) }, ErrNotFinite},
		{"volume zero inner", func() error { return e1(HollowCylinderVolume(2, 0, 1)) }, ErrNonPositive},
		{"volume zero height", func() error { return e1(HollowCylinderVolume(2, 1, 0)) }, ErrNonPositive},
		{"volume inner ≥ outer", func() error { return e1(HollowCylinderVolume(2, 2, 1)) }, ErrOutOfRange},
		{"volume overflow", func() error { return e1(HollowCylinderVolume(1e200, 1, 1)) }, ErrOverflow},
		{"surface NaN", func() error { return e1(HollowCylinderSurfaceArea(2, nan, 1)) }, ErrNotFinite},
		{"surface negative height", func() error { return e1(HollowCylinderSurfaceArea(2, 1, -1)) }, ErrNonPositive},
		{"surface inner > outer", func() error { return e1(HollowCylinderSurfaceArea(1, 2, 1)) }, ErrOutOfRange},
		{"surface overflow", func() error { return e1(HollowCylinderSurfaceArea(1e200, 1, 1e200)) }, ErrOverflow},
	})
}

// Oracle: the outer cylinder less the inner one, on walls that are not thin; and the area as the two
// lateral walls plus two annular ends.
func TestHollowCylinder_vsCylinders(t *testing.T) {
	rng := newRNG()
	for range 400 {
		outer, h := 1+rng.Float64()*10, 0.1+rng.Float64()*10
		inner := outer * (0.01 + 0.8*rng.Float64())
		vOut, _ := CylinderVolume(outer, h)
		vIn, _ := CylinderVolume(inner, h)
		if got, err := HollowCylinderVolume(outer, inner, h); err != nil || !relNear(got, vOut-vIn, 1e-14) {
			t.Fatalf("HollowCylinderVolume(%v, %v, %v) = (%v, %v), want %v", outer, inner, h, got, err, vOut-vIn)
		}
		lOut, _ := CylinderLateralArea(outer, h)
		lIn, _ := CylinderLateralArea(inner, h)
		ring, _ := AnnulusArea(outer, inner)
		if got, err := HollowCylinderSurfaceArea(outer, inner, h); err != nil || !relNear(got, lOut+lIn+2*ring, 1e-14) {
			t.Fatalf("HollowCylinderSurfaceArea(%v, %v, %v) = (%v, %v), want %v", outer, inner, h, got, err, lOut+lIn+2*ring)
		}
	}
}

func TestHollowCylinder_allocs(t *testing.T) {
	assertNoAllocs(t, "HollowCylinderVolume", func() { sink, _ = HollowCylinderVolume(3, 2, 5) })
	assertNoAllocs(t, "HollowCylinderSurfaceArea", func() { sink, _ = HollowCylinderSurfaceArea(3, 2, 5) })
}

func BenchmarkHollowCylinderVolume(b *testing.B) {
	b.ReportAllocs()
	outer, inner, h := 3.0, 2.0, 5.0
	for b.Loop() {
		sink, _ = HollowCylinderVolume(outer, inner, h)
	}
}

func BenchmarkHollowCylinderSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	outer, inner, h := 3.0, 2.0, 5.0
	for b.Loop() {
		sink, _ = HollowCylinderSurfaceArea(outer, inner, h)
	}
}
