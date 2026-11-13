package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.6.6 p. 334 (circumference 2πr, area πr²), hand-derived: r 1 → 2π;
// r 2 → 4π.
func TestCircle_example(t *testing.T) {
	runValCases(t, []valCase{
		{"CircleCircumference(1)", func() (float64, error) { return CircleCircumference(1) }, 2 * math.Pi, tolFew},
		{"CircleArea(2)", func() (float64, error) { return CircleArea(2) }, 4 * math.Pi, tolFew},
	})
}

func TestCircle_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"circumference, tiny radius", func() (float64, error) { return CircleCircumference(1e-300) }, 2 * math.Pi * 1e-300, tolFew},
		{"area, large radius below overflow", func() (float64, error) { return CircleArea(1e150) }, math.Pi * 1e300, tolFew},
	})
}

func TestCircle_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"circumference NaN", func() error { return e1(CircleCircumference(nan)) }, ErrNotFinite},
		{"circumference zero", func() error { return e1(CircleCircumference(0)) }, ErrNonPositive},
		{"circumference overflow", func() error { return e1(CircleCircumference(math.MaxFloat64)) }, ErrOverflow},
		{"area Inf", func() error { return e1(CircleArea(inf)) }, ErrNotFinite},
		{"area negative", func() error { return e1(CircleArea(-1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(CircleArea(1e160)) }, ErrOverflow},
	})
}

// Oracle: a regular n-gon inscribed in the circle approaches it (CRC §4.5.3): perimeter and area within
// the next series term, about (2π)²/(6n²) relative.
func TestCircle_vsRegularPolygon(t *testing.T) {
	const n = 1_000_000
	r := 2.5
	c, _ := CircleCircumference(r)
	a, _ := CircleArea(r)
	side, _ := RegularPolygonSideInscribed(n, r)
	area, _ := RegularPolygonAreaFromCircumradius(n, r)
	if !relNear(float64(n)*side, c, 1e-11) || !relNear(area, a, 1e-11) {
		t.Errorf("r %v: circumference %v vs %v, area %v vs %v", r, c, float64(n)*side, a, area)
	}
}

func TestCircle_allocs(t *testing.T) {
	assertNoAllocs(t, "CircleCircumference", func() { sink, _ = CircleCircumference(2) })
	assertNoAllocs(t, "CircleArea", func() { sink, _ = CircleArea(2) })
}

func BenchmarkCircleCircumference(b *testing.B) {
	b.ReportAllocs()
	r := 2.0
	for b.Loop() {
		sink, _ = CircleCircumference(r)
	}
}

func BenchmarkCircleArea(b *testing.B) {
	b.ReportAllocs()
	r := 2.0
	for b.Loop() {
		sink, _ = CircleArea(r)
	}
}
