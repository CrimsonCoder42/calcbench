package geometry

import "testing"

// Worked example: CRC 31st ed. §4.5.2 p. 323 (area = ½(AB + CD)h), hand-derived: 3, 5, 2 → 8.
func TestTrapezoidArea_example(t *testing.T) {
	runValCases(t, []valCase{
		{"TrapezoidArea(3, 5, 2)", func() (float64, error) { return TrapezoidArea(3, 5, 2) }, 8, 0},
	})
}

func TestTrapezoidArea_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"equal bases is a parallelogram", func() (float64, error) { return TrapezoidArea(4, 4, 3) }, 12, 0},
		{"bases near the float64 limit, no overflow of the sum", func() (float64, error) { return TrapezoidArea(1.5e308, 1.5e308, 1e-10) }, 1.5e298, tolFew},
	})
}

func TestTrapezoidArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(TrapezoidArea(nan, 1, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(TrapezoidArea(1, 1, inf)) }, ErrNotFinite},
		{"zero base", func() error { return e1(TrapezoidArea(1, 0, 1)) }, ErrNonPositive},
		{"negative height", func() error { return e1(TrapezoidArea(1, 1, -1)) }, ErrNonPositive},
		{"overflow", func() error { return e1(TrapezoidArea(1e200, 1e200, 1e200)) }, ErrOverflow},
	})
}

// Oracle: a trapezoid splits along a diagonal into two triangles of the same height.
func TestTrapezoidArea_vsTriangles(t *testing.T) {
	rng := newRNG()
	for range 400 {
		a, b, h := rng.Float64()*100+0.1, rng.Float64()*100+0.1, rng.Float64()*100+0.1
		t1, _ := TriangleAreaBaseHeight(a, h)
		t2, _ := TriangleAreaBaseHeight(b, h)
		if got, err := TrapezoidArea(a, b, h); err != nil || !relNear(got, t1+t2, tolSome) {
			t.Fatalf("TrapezoidArea(%v, %v, %v) = (%v, %v), want %v", a, b, h, got, err, t1+t2)
		}
	}
}

func TestTrapezoidArea_allocs(t *testing.T) {
	assertNoAllocs(t, "TrapezoidArea", func() { sink, _ = TrapezoidArea(3, 5, 2) })
}

func BenchmarkTrapezoidArea(b *testing.B) {
	b.ReportAllocs()
	x, y, h := 3.0, 5.0, 2.0
	for b.Loop() {
		sink, _ = TrapezoidArea(x, y, h)
	}
}
