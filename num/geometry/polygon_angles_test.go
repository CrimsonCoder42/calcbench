package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.5 p. 318 ("2(k − 2) right angles") and §4.5.3 p. 324
// (θ = ((k − 2)/k)·180°), hand-derived: pentagon 3π; hexagon 2π/3; octagon exterior π/4.
func TestPolygonAngles_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PolygonInteriorAngleSumRad(5)", func() (float64, error) { return PolygonInteriorAngleSumRad(5) }, 3 * math.Pi, tolFew},
		{"PolygonInteriorAngleSumRad(3)", func() (float64, error) { return PolygonInteriorAngleSumRad(3) }, math.Pi, 0},
		{"RegularPolygonInteriorAngleRad(6)", func() (float64, error) { return RegularPolygonInteriorAngleRad(6) }, 2 * math.Pi / 3, tolFew},
		{"RegularPolygonInteriorAngleRad(4)", func() (float64, error) { return RegularPolygonInteriorAngleRad(4) }, math.Pi / 2, tolFew},
		{"RegularPolygonExteriorAngleRad(8)", func() (float64, error) { return RegularPolygonExteriorAngleRad(8) }, math.Pi / 4, 0},
		{"RegularPolygonExteriorAngleRad(3)", func() (float64, error) { return RegularPolygonExteriorAngleRad(3) }, 2 * math.Pi / 3, tolFew},
	})
}

func TestPolygonAngles_largeN(t *testing.T) {
	runValCases(t, []valCase{
		{"interior sum, n = MaxInt", func() (float64, error) { return PolygonInteriorAngleSumRad(math.MaxInt) }, float64(math.MaxInt-2) * math.Pi, tolFew},
		{"interior angle tends to π", func() (float64, error) { return RegularPolygonInteriorAngleRad(math.MaxInt) }, math.Pi, tolFew},
		{"exterior angle tends to 0", func() (float64, error) { return RegularPolygonExteriorAngleRad(math.MaxInt) }, 2 * math.Pi / float64(math.MaxInt), tolFew},
	})
}

func TestPolygonAngles_errors(t *testing.T) {
	var cases []errCase
	for _, n := range []int{2, 0, -5, math.MinInt} {
		cases = append(cases,
			errCase{"sum", func() error { return e1(PolygonInteriorAngleSumRad(n)) }, ErrTooFewSides},
			errCase{"interior", func() error { return e1(RegularPolygonInteriorAngleRad(n)) }, ErrTooFewSides},
			errCase{"exterior", func() error { return e1(RegularPolygonExteriorAngleRad(n)) }, ErrTooFewSides},
		)
	}
	runErrCases(t, cases)
}

// Property: interior + exterior = π, and n interior angles make the interior sum.
func TestPolygonAngles_consistent(t *testing.T) {
	for n := 3; n <= 1000; n++ {
		in, _ := RegularPolygonInteriorAngleRad(n)
		ex, _ := RegularPolygonExteriorAngleRad(n)
		sum, _ := PolygonInteriorAngleSumRad(n)
		if !relNear(in+ex, math.Pi, tolSome) || !relNear(float64(n)*in, sum, tolSome) {
			t.Fatalf("n = %d: interior %v, exterior %v, sum %v", n, in, ex, sum)
		}
	}
}

func TestPolygonAngles_allocs(t *testing.T) {
	assertNoAllocs(t, "PolygonInteriorAngleSumRad", func() { sink, _ = PolygonInteriorAngleSumRad(7) })
	assertNoAllocs(t, "RegularPolygonInteriorAngleRad", func() { sink, _ = RegularPolygonInteriorAngleRad(7) })
	assertNoAllocs(t, "RegularPolygonExteriorAngleRad", func() { sink, _ = RegularPolygonExteriorAngleRad(7) })
}

func BenchmarkPolygonInteriorAngleSumRad(b *testing.B) {
	b.ReportAllocs()
	n := 7
	for b.Loop() {
		sink, _ = PolygonInteriorAngleSumRad(n)
	}
}

func BenchmarkRegularPolygonInteriorAngleRad(b *testing.B) {
	b.ReportAllocs()
	n := 7
	for b.Loop() {
		sink, _ = RegularPolygonInteriorAngleRad(n)
	}
}

func BenchmarkRegularPolygonExteriorAngleRad(b *testing.B) {
	b.ReportAllocs()
	n := 7
	for b.Loop() {
		sink, _ = RegularPolygonExteriorAngleRad(n)
	}
}
