package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.6.3 eq. (4.6.11) (area of C = πab), hand-derived: 5, 3 → 15π.
func TestEllipseArea_example(t *testing.T) {
	runValCases(t, []valCase{
		{"EllipseArea(5, 3)", func() (float64, error) { return EllipseArea(5, 3) }, 15 * math.Pi, tolFew},
	})
}

func TestEllipseArea_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"circle", func() (float64, error) { return EllipseArea(2, 2) }, 4 * math.Pi, tolFew},
		{"either order", func() (float64, error) { return EllipseArea(3, 5) }, 15 * math.Pi, tolFew},
		{"no spurious overflow", func() (float64, error) { return EllipseArea(1e300, 1e-300) }, math.Pi, tolFew},
	})
}

func TestEllipseArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(EllipseArea(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(EllipseArea(1, inf)) }, ErrNotFinite},
		{"zero", func() error { return e1(EllipseArea(0, 1)) }, ErrNonPositive},
		{"overflow", func() error { return e1(EllipseArea(1e200, 1e200)) }, ErrOverflow},
	})
}

// Oracle: the inscribed 4096-gon through (a cos t, b sin t) has area ½·N·ab·sin(2π/N), within
// (2π/N)²/6 ≈ 4e-7 of πab.
func TestEllipseArea_vsPolygon(t *testing.T) {
	const n = 4096
	vs := make([]Point, n)
	rng := newRNG()
	for range 20 {
		a, b := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		for i := range vs {
			s, c := math.Sincos(2 * math.Pi * float64(i) / n)
			vs[i] = Point{a * c, b * s}
		}
		poly, err := PolygonArea(vs)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := EllipseArea(a, b); err != nil || !relNear(got, poly, 1e-6) {
			t.Fatalf("EllipseArea(%v, %v) = (%v, %v), inscribed polygon %v", a, b, got, err, poly)
		}
	}
}

func TestEllipseArea_allocs(t *testing.T) {
	assertNoAllocs(t, "EllipseArea", func() { sink, _ = EllipseArea(5, 3) })
}

func BenchmarkEllipseArea(b *testing.B) {
	b.ReportAllocs()
	x, y := 5.0, 3.0
	for b.Loop() {
		sink, _ = EllipseArea(x, y)
	}
}
