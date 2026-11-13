package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.5.1 p. 320 (c² = a² + b², Pythagoras), hand-derived: 3, 4 → 5 and
// 5, 3 → 4.
func TestRightTriangle_example(t *testing.T) {
	runValCases(t, []valCase{
		{"RightTriangleHypotenuse(3, 4)", func() (float64, error) { return RightTriangleHypotenuse(3, 4) }, 5, 0},
		{"RightTriangleLeg(5, 3)", func() (float64, error) { return RightTriangleLeg(5, 3) }, 4, tolFew},
	})
}

func TestRightTriangle_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"hypotenuse, isosceles", func() (float64, error) { return RightTriangleHypotenuse(1, 1) }, math.Sqrt2, tolFew},
		{"hypotenuse, no overflow of the squares", func() (float64, error) { return RightTriangleHypotenuse(3e200, 4e200) }, 5e200, tolFew},
		{"leg, near the hypotenuse (no cancellation)", func() (float64, error) { return RightTriangleLeg(1, 1-0x1p-52) }, math.Sqrt(0x1p-51 - 0x1p-104), tolFew},
		{"leg, c + a overflows", func() (float64, error) { return RightTriangleLeg(1.5e308, 0.9e308) }, 1.2e308, tolSome},
	})
}

func TestRightTriangle_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"hypotenuse NaN", func() error { return e1(RightTriangleHypotenuse(nan, 1)) }, ErrNotFinite},
		{"hypotenuse Inf", func() error { return e1(RightTriangleHypotenuse(1, inf)) }, ErrNotFinite},
		{"hypotenuse zero leg", func() error { return e1(RightTriangleHypotenuse(0, 1)) }, ErrNonPositive},
		{"hypotenuse overflow", func() error { return e1(RightTriangleHypotenuse(1.5e308, 1.5e308)) }, ErrOverflow},
		{"leg NaN", func() error { return e1(RightTriangleLeg(nan, 1)) }, ErrNotFinite},
		{"leg Inf", func() error { return e1(RightTriangleLeg(inf, 1)) }, ErrNotFinite},
		{"leg negative", func() error { return e1(RightTriangleLeg(5, -3)) }, ErrNonPositive},
		{"leg equals hypotenuse", func() error { return e1(RightTriangleLeg(5, 5)) }, ErrNotTriangle},
		{"leg longer than hypotenuse", func() error { return e1(RightTriangleLeg(5, 6)) }, ErrNotTriangle},
	})
}

// Oracle: the leg recovered from a hypotenuse built from two legs (leg ratios kept in [0.1, 10], where
// rounding the hypotenuse costs at most a few ulps of the recovered leg).
func TestRightTriangle_roundTrip(t *testing.T) {
	rng := newRNG()
	for range 400 {
		a := 1 + rng.Float64()*99
		b := a * math.Pow(10, rng.Float64()*2-1)
		c, err := RightTriangleHypotenuse(a, b)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := RightTriangleLeg(c, a); err != nil || !relNear(got, b, 1e-13) {
			t.Fatalf("RightTriangleLeg(%v, %v) = (%v, %v), want %v", c, a, got, err, b)
		}
	}
}

func TestRightTriangle_allocs(t *testing.T) {
	assertNoAllocs(t, "RightTriangleHypotenuse", func() { sink, _ = RightTriangleHypotenuse(3, 4) })
	assertNoAllocs(t, "RightTriangleLeg", func() { sink, _ = RightTriangleLeg(5, 3) })
}

func BenchmarkRightTriangleHypotenuse(b *testing.B) {
	b.ReportAllocs()
	x, y := 3.0, 4.0
	for b.Loop() {
		sink, _ = RightTriangleHypotenuse(x, y)
	}
}

func BenchmarkRightTriangleLeg(b *testing.B) {
	b.ReportAllocs()
	c, a := 5.0, 3.0
	for b.Loop() {
		sink, _ = RightTriangleLeg(c, a)
	}
}
