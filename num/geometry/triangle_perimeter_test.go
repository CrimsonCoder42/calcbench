package geometry

import "testing"

// Worked example: CRC 31st ed. §4.5.1 (s = ½(a + b + c)), hand-derived: 3, 4, 5 → 12.
func TestTrianglePerimeter_example(t *testing.T) {
	runValCases(t, []valCase{
		{"TrianglePerimeter(3, 4, 5)", func() (float64, error) { return TrianglePerimeter(3, 4, 5) }, 12, 0},
	})
}

func TestTrianglePerimeter_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"equilateral", func() (float64, error) { return TrianglePerimeter(2, 2, 2) }, 6, 0},
		{"any order", func() (float64, error) { return TrianglePerimeter(5, 3, 4) }, 12, 0},
		{"degenerate is accepted", func() (float64, error) { return TrianglePerimeter(1, 2, 3) }, 6, 0},
		{"needle (Kahan 2014 Table 1, row 3)", func() (float64, error) { return TrianglePerimeter(100000, 99999.99979, 0.00029) }, 200000.00008, tolFew},
	})
}

func TestTrianglePerimeter_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(TrianglePerimeter(nan, 1, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(TrianglePerimeter(1, inf, 1)) }, ErrNotFinite},
		{"zero side", func() error { return e1(TrianglePerimeter(1, 1, 0)) }, ErrNonPositive},
		{"negative side (Kahan 2014 Table 1, row 2)", func() error { return e1(TrianglePerimeter(-3, 4, 2)) }, ErrNonPositive},
		{"inequality fails", func() error { return e1(TrianglePerimeter(1, 2, 3.5)) }, ErrNotTriangle},
		{"overflow", func() error { return e1(TrianglePerimeter(1e308, 1e308, 1e308)) }, ErrOverflow},
	})
}

func TestTrianglePerimeter_allocs(t *testing.T) {
	assertNoAllocs(t, "TrianglePerimeter", func() { sink, _ = TrianglePerimeter(3, 4, 5) })
}

func BenchmarkTrianglePerimeter(b *testing.B) {
	b.ReportAllocs()
	x, y, z := 3.0, 4.0, 5.0
	for b.Loop() {
		sink, _ = TrianglePerimeter(x, y, z)
	}
}
