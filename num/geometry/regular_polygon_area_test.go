package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.5.3 table, p. 324 (area of the regular k-gon of side a, to 5 decimal
// places): triangle 0.43301a², square a², pentagon 1.72048a², hexagon
// 2.59808a², octagon 4.82843a², dodecagon 11.19615a². The table's r (apothem) and R (circumradius) columns
// give the same areas through the other two functions.
func TestRegularPolygonArea_crcTable(t *testing.T) {
	// The table rounds to 5 decimal places: at most 1.2e-5 relative for the areas (0.43301), and an apothem
	// or circumradius so rounded is squared, which doubles its relative error (up to 3.5e-5 for 0.28868).
	const tol, radiusTol = 2e-5, 1e-4
	rows := []struct {
		k          int
		area, r, R float64
	}{
		{3, 0.43301, 0.28868, 0.57735},
		{4, 1, 0.5, 0.70711},
		{5, 1.72048, 0.68819, 0.85065},
		{6, 2.59808, 0.86603, 1},
		{8, 4.82843, 1.20711, 1.30656},
		{12, 11.19615, 1.86603, 1.93185},
	}
	for _, row := range rows {
		if got, err := RegularPolygonAreaFromSide(row.k, 1); err != nil || !relNear(got, row.area, tol) {
			t.Errorf("RegularPolygonAreaFromSide(%d, 1) = (%v, %v), want %v", row.k, got, err, row.area)
		}
		if got, err := RegularPolygonAreaFromApothem(row.k, row.r); err != nil || !relNear(got, row.area, radiusTol) {
			t.Errorf("RegularPolygonAreaFromApothem(%d, %v) = (%v, %v), want %v", row.k, row.r, got, err, row.area)
		}
		if got, err := RegularPolygonAreaFromCircumradius(row.k, row.R); err != nil || !relNear(got, row.area, radiusTol) {
			t.Errorf("RegularPolygonAreaFromCircumradius(%d, %v) = (%v, %v), want %v", row.k, row.R, got, err, row.area)
		}
	}
}

// Hand-derived: the hexagon of side 1 is six equilateral triangles, 3√3/2; its apothem is √3/2 and its
// circumradius 1.
func TestRegularPolygonArea_hexagon(t *testing.T) {
	want := 3 * math.Sqrt(3) / 2
	runValCases(t, []valCase{
		{"side 1", func() (float64, error) { return RegularPolygonAreaFromSide(6, 1) }, want, tolFew},
		{"apothem √3/2", func() (float64, error) { return RegularPolygonAreaFromApothem(6, math.Sqrt(3)/2) }, want, tolFew},
		{"circumradius 1", func() (float64, error) { return RegularPolygonAreaFromCircumradius(6, 1) }, want, tolFew},
		{"huge n tends to the circle", func() (float64, error) { return RegularPolygonAreaFromCircumradius(1<<40, 1) }, math.Pi, 1e-12},
	})
}

func TestRegularPolygonArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"side: too few sides", func() error { return e1(RegularPolygonAreaFromSide(2, 1)) }, ErrTooFewSides},
		{"side: NaN", func() error { return e1(RegularPolygonAreaFromSide(5, nan)) }, ErrNotFinite},
		{"side: zero", func() error { return e1(RegularPolygonAreaFromSide(5, 0)) }, ErrNonPositive},
		{"side: overflow", func() error { return e1(RegularPolygonAreaFromSide(5, 1e200)) }, ErrOverflow},
		{"apothem: too few sides", func() error { return e1(RegularPolygonAreaFromApothem(-1, 1)) }, ErrTooFewSides},
		{"apothem: Inf", func() error { return e1(RegularPolygonAreaFromApothem(5, inf)) }, ErrNotFinite},
		{"apothem: negative", func() error { return e1(RegularPolygonAreaFromApothem(5, -1)) }, ErrNonPositive},
		{"apothem: overflow", func() error { return e1(RegularPolygonAreaFromApothem(5, 1e200)) }, ErrOverflow},
		{"circumradius: too few sides", func() error { return e1(RegularPolygonAreaFromCircumradius(0, 1)) }, ErrTooFewSides},
		{"circumradius: NaN", func() error { return e1(RegularPolygonAreaFromCircumradius(5, nan)) }, ErrNotFinite},
		{"circumradius: zero", func() error { return e1(RegularPolygonAreaFromCircumradius(5, 0)) }, ErrNonPositive},
		{"circumradius: overflow", func() error { return e1(RegularPolygonAreaFromCircumradius(5, 1e200)) }, ErrOverflow},
	})
}

// Oracle: the three forms agree, with apothem s/(2 tan(π/n)) and circumradius s/(2 sin(π/n)) (CRC §4.5.3).
func TestRegularPolygonArea_formsAgree(t *testing.T) {
	rng := newRNG()
	for range 400 {
		n := 3 + rng.IntN(2000)
		s := rng.Float64()*10 + 0.01
		fromSide, err := RegularPolygonAreaFromSide(n, s)
		if err != nil {
			t.Fatal(err)
		}
		fromApothem, _ := RegularPolygonAreaFromApothem(n, s/(2*math.Tan(math.Pi/float64(n))))
		fromCircum, _ := RegularPolygonAreaFromCircumradius(n, s/(2*math.Sin(math.Pi/float64(n))))
		if !relNear(fromApothem, fromSide, 1e-13) || !relNear(fromCircum, fromSide, 1e-13) {
			t.Fatalf("n %d, s %v: side %v, apothem %v, circumradius %v", n, s, fromSide, fromApothem, fromCircum)
		}
	}
}

func TestRegularPolygonArea_allocs(t *testing.T) {
	assertNoAllocs(t, "RegularPolygonAreaFromSide", func() { sink, _ = RegularPolygonAreaFromSide(5, 2) })
	assertNoAllocs(t, "RegularPolygonAreaFromApothem", func() { sink, _ = RegularPolygonAreaFromApothem(5, 2) })
	assertNoAllocs(t, "RegularPolygonAreaFromCircumradius", func() { sink, _ = RegularPolygonAreaFromCircumradius(5, 2) })
}

func BenchmarkRegularPolygonAreaFromSide(b *testing.B) {
	b.ReportAllocs()
	n, s := 5, 2.0
	for b.Loop() {
		sink, _ = RegularPolygonAreaFromSide(n, s)
	}
}

func BenchmarkRegularPolygonAreaFromApothem(b *testing.B) {
	b.ReportAllocs()
	n, r := 5, 2.0
	for b.Loop() {
		sink, _ = RegularPolygonAreaFromApothem(n, r)
	}
}

func BenchmarkRegularPolygonAreaFromCircumradius(b *testing.B) {
	b.ReportAllocs()
	n, r := 5, 2.0
	for b.Loop() {
		sink, _ = RegularPolygonAreaFromCircumradius(n, r)
	}
}
