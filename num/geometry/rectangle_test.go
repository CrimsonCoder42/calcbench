package geometry

import "testing"

// Worked examples: CRC 31st ed. §4.5.2 p. 322 (rectangle: area = ab), hand-derived: 3 × 4 → 12, and
// perimeter 2(3 + 4) → 14.
func TestRectangle_example(t *testing.T) {
	runValCases(t, []valCase{
		{"RectangleArea(3, 4)", func() (float64, error) { return RectangleArea(3, 4) }, 12, 0},
		{"RectanglePerimeter(3, 4)", func() (float64, error) { return RectanglePerimeter(3, 4) }, 14, 0},
	})
}

func TestRectangle_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"square area", func() (float64, error) { return RectangleArea(5, 5) }, 25, 0},
		{"square perimeter", func() (float64, error) { return RectanglePerimeter(5, 5) }, 20, 0},
		{"area, no spurious overflow", func() (float64, error) { return RectangleArea(1e300, 1e-300) }, 1, tolFew},
		{"perimeter, tiny sides", func() (float64, error) { return RectanglePerimeter(5e-324, 5e-324) }, 2e-323, 0},
	})
}

func TestRectangle_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"area NaN", func() error { return e1(RectangleArea(nan, 1)) }, ErrNotFinite},
		{"area Inf", func() error { return e1(RectangleArea(1, inf)) }, ErrNotFinite},
		{"area zero", func() error { return e1(RectangleArea(0, 1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(RectangleArea(1e200, 1e200)) }, ErrOverflow},
		{"perimeter NaN", func() error { return e1(RectanglePerimeter(1, nan)) }, ErrNotFinite},
		{"perimeter negative", func() error { return e1(RectanglePerimeter(-1, 1)) }, ErrNonPositive},
		{"perimeter overflow", func() error { return e1(RectanglePerimeter(1e308, 1e308)) }, ErrOverflow},
	})
}

func TestRectangle_allocs(t *testing.T) {
	assertNoAllocs(t, "RectangleArea", func() { sink, _ = RectangleArea(3, 4) })
	assertNoAllocs(t, "RectanglePerimeter", func() { sink, _ = RectanglePerimeter(3, 4) })
}

func BenchmarkRectangleArea(b *testing.B) {
	b.ReportAllocs()
	w, h := 3.0, 4.0
	for b.Loop() {
		sink, _ = RectangleArea(w, h)
	}
}

func BenchmarkRectanglePerimeter(b *testing.B) {
	b.ReportAllocs()
	w, h := 3.0, 4.0
	for b.Loop() {
		sink, _ = RectanglePerimeter(w, h)
	}
}
