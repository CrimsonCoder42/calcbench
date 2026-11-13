package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.5.2 p. 322 (h = a sin A, area = bh), hand-derived: 3, 4, π/6 → 6.
func TestParallelogramAreaSAS_example(t *testing.T) {
	runValCases(t, []valCase{
		{"ParallelogramAreaSAS(3, 4, π/6)", func() (float64, error) { return ParallelogramAreaSAS(3, 4, math.Pi/6) }, 6, tolFew},
	})
}

func TestParallelogramAreaSAS_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"rectangle (right angle)", func() (float64, error) { return ParallelogramAreaSAS(3, 4, math.Pi/2) }, 12, 0},
		{"rhombus", func() (float64, error) { return ParallelogramAreaSAS(2, 2, math.Pi/6) }, 2, tolFew},
		{"obtuse angle equals its supplement", func() (float64, error) { return ParallelogramAreaSAS(3, 4, 5*math.Pi/6) }, 6, tolFew},
	})
}

func TestParallelogramAreaSAS_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(ParallelogramAreaSAS(nan, 1, 1)) }, ErrNotFinite},
		{"Inf angle", func() error { return e1(ParallelogramAreaSAS(1, 1, -inf)) }, ErrNotFinite},
		{"zero side", func() error { return e1(ParallelogramAreaSAS(1, 0, 1)) }, ErrNonPositive},
		{"angle 0", func() error { return e1(ParallelogramAreaSAS(1, 1, 0)) }, ErrOutOfRange},
		{"angle π", func() error { return e1(ParallelogramAreaSAS(1, 1, math.Pi)) }, ErrOutOfRange},
		{"overflow", func() error { return e1(ParallelogramAreaSAS(1e200, 1e200, 1)) }, ErrOverflow},
	})
}

// Oracle: a parallelogram is two congruent triangles with the same two sides and angle.
func TestParallelogramAreaSAS_vsTriangle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		a, b, angle := rng.Float64()*100+0.1, rng.Float64()*100+0.1, 0.01+rng.Float64()*3.1
		tri, err := TriangleAreaSAS(a, b, angle)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := ParallelogramAreaSAS(a, b, angle); err != nil || !relNear(got, 2*tri, tolSome) {
			t.Fatalf("ParallelogramAreaSAS(%v, %v, %v) = (%v, %v), want %v", a, b, angle, got, err, 2*tri)
		}
	}
}

func TestParallelogramAreaSAS_allocs(t *testing.T) {
	assertNoAllocs(t, "ParallelogramAreaSAS", func() { sink, _ = ParallelogramAreaSAS(3, 4, 0.5) })
}

func BenchmarkParallelogramAreaSAS(b *testing.B) {
	b.ReportAllocs()
	x, y, angle := 3.0, 4.0, 0.5
	for b.Loop() {
		sink, _ = ParallelogramAreaSAS(x, y, angle)
	}
}
