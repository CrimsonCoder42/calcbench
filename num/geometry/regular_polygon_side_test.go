package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.5.3 p. 324 (a = 2R sin(180°/k) = 2r tan(180°/k)); hand-derived:
// square in the unit circle √2, hexagon in the unit circle 1 (the table's R = a), square about the unit
// circle 2. The table's hexagon apothem 0.86603a (5 decimals, 5.3e-6 relative) gives side 1.
func TestRegularPolygonSide_example(t *testing.T) {
	runValCases(t, []valCase{
		{"inscribed square", func() (float64, error) { return RegularPolygonSideInscribed(4, 1) }, math.Sqrt2, tolFew},
		{"inscribed hexagon", func() (float64, error) { return RegularPolygonSideInscribed(6, 1) }, 1, tolFew},
		{"circumscribed square", func() (float64, error) { return RegularPolygonSideCircumscribed(4, 1) }, 2, tolFew},
		{"circumscribed hexagon, CRC table", func() (float64, error) { return RegularPolygonSideCircumscribed(6, 0.86603) }, 1, 1e-5},
		{"inscribed triangle", func() (float64, error) { return RegularPolygonSideInscribed(3, 1) }, math.Sqrt(3), tolFew},
	})
}

func TestRegularPolygonSide_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"inscribed: too few sides", func() error { return e1(RegularPolygonSideInscribed(2, 1)) }, ErrTooFewSides},
		{"inscribed: NaN", func() error { return e1(RegularPolygonSideInscribed(5, nan)) }, ErrNotFinite},
		{"inscribed: zero", func() error { return e1(RegularPolygonSideInscribed(5, 0)) }, ErrNonPositive},
		{"inscribed: overflow", func() error { return e1(RegularPolygonSideInscribed(3, math.MaxFloat64)) }, ErrOverflow},
		{"circumscribed: too few sides", func() error { return e1(RegularPolygonSideCircumscribed(1, 1)) }, ErrTooFewSides},
		{"circumscribed: Inf", func() error { return e1(RegularPolygonSideCircumscribed(5, inf)) }, ErrNotFinite},
		{"circumscribed: negative", func() error { return e1(RegularPolygonSideCircumscribed(5, -2)) }, ErrNonPositive},
		{"circumscribed: overflow", func() error { return e1(RegularPolygonSideCircumscribed(3, math.MaxFloat64)) }, ErrOverflow},
	})
}

// Oracle: n inscribed sides make a perimeter just under 2πR and n circumscribed sides one just over it,
// both within the next term of the sine and tangent series, π³/n² relative.
func TestRegularPolygonSide_perimeterBracketsCircle(t *testing.T) {
	for _, n := range []int{3, 10, 100, 1000, 100000} {
		in, _ := RegularPolygonSideInscribed(n, 1)
		out, _ := RegularPolygonSideCircumscribed(n, 1)
		pIn, pOut := float64(n)*in, float64(n)*out
		bound := math.Pow(math.Pi, 3) / float64(n*n)
		if !(pIn < 2*math.Pi && 2*math.Pi < pOut) || (2*math.Pi-pIn) > bound || (pOut-2*math.Pi) > 2*bound {
			t.Errorf("n %d: inscribed perimeter %v, circumscribed %v, circle %v", n, pIn, pOut, 2*math.Pi)
		}
	}
}

func TestRegularPolygonSide_allocs(t *testing.T) {
	assertNoAllocs(t, "RegularPolygonSideInscribed", func() { sink, _ = RegularPolygonSideInscribed(5, 2) })
	assertNoAllocs(t, "RegularPolygonSideCircumscribed", func() { sink, _ = RegularPolygonSideCircumscribed(5, 2) })
}

func BenchmarkRegularPolygonSideInscribed(b *testing.B) {
	b.ReportAllocs()
	n, r := 5, 2.0
	for b.Loop() {
		sink, _ = RegularPolygonSideInscribed(n, r)
	}
}

func BenchmarkRegularPolygonSideCircumscribed(b *testing.B) {
	b.ReportAllocs()
	n, r := 5, 2.0
	for b.Loop() {
		sink, _ = RegularPolygonSideCircumscribed(n, r)
	}
}
