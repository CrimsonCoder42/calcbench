package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.18 p. 364 ((4/3)πabc), hand-derived: 3, 2, 1 → 8π.
func TestEllipsoidVolume_example(t *testing.T) {
	runValCases(t, []valCase{
		{"EllipsoidVolume(3, 2, 1)", func() (float64, error) { return EllipsoidVolume(3, 2, 1) }, 8 * math.Pi, tolFew},
	})
}

func TestEllipsoidVolume_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"any order", func() (float64, error) { return EllipsoidVolume(1, 3, 2) }, 8 * math.Pi, tolFew},
		{"no spurious overflow", func() (float64, error) { return EllipsoidVolume(1e200, 1e200, 1e-300) }, 4 * math.Pi / 3 * 1e100, tolSome},
	})
}

func TestEllipsoidVolume_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(EllipsoidVolume(nan, 1, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(EllipsoidVolume(1, 1, inf)) }, ErrNotFinite},
		{"zero", func() error { return e1(EllipsoidVolume(1, 0, 1)) }, ErrNonPositive},
		{"overflow", func() error { return e1(EllipsoidVolume(1e200, 1e200, 1)) }, ErrOverflow},
	})
}

// Oracle: the ellipsoid is the unit sphere stretched by a, b and c, so its volume is abc times the unit
// sphere's.
func TestEllipsoidVolume_vsSphere(t *testing.T) {
	unit, _ := SphereVolume(1)
	rng := newRNG()
	for range 400 {
		a, b, c := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		if got, err := EllipsoidVolume(a, b, c); err != nil || !relNear(got, unit*a*b*c, tolSome) {
			t.Fatalf("EllipsoidVolume(%v, %v, %v) = (%v, %v), want %v", a, b, c, got, err, unit*a*b*c)
		}
	}
}

func TestEllipsoidVolume_allocs(t *testing.T) {
	assertNoAllocs(t, "EllipsoidVolume", func() { sink, _ = EllipsoidVolume(3, 2, 1) })
}

func BenchmarkEllipsoidVolume(b *testing.B) {
	b.ReportAllocs()
	x, y, z := 3.0, 2.0, 1.0
	for b.Loop() {
		sink, _ = EllipsoidVolume(x, y, z)
	}
}
