package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.6.3 p. 330 (x²/a² + y²/b² = 1), hand-derived: a 5, b 3 at offset 1.8
// from the centre (depth 1.2 from the vertex): x = 5·√(1 − 0.36) = 4, so the chord is 8.
func TestEllipseChord_example(t *testing.T) {
	runValCases(t, []valCase{
		{"EllipseChordFromCenter(5, 3, 1.8)", func() (float64, error) { return EllipseChordFromCenter(5, 3, 1.8) }, 8, tolFew},
		{"EllipseChordFromEdge(5, 3, 1.2)", func() (float64, error) { return EllipseChordFromEdge(5, 3, 1.2) }, 8, tolFew},
	})
}

func TestEllipseChord_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"centre: offset 0 is the axis", func() (float64, error) { return EllipseChordFromCenter(5, 3, 0) }, 10, tolFew},
		{"centre: offset b is the vertex", func() (float64, error) { return EllipseChordFromCenter(5, 3, 3) }, 0, 0},
		{"centre: axes swapped (a vertical chord)", func() (float64, error) { return EllipseChordFromCenter(3, 5, 3) }, 4.8, tolFew},
		{"edge: depth 0 is the vertex", func() (float64, error) { return EllipseChordFromEdge(5, 3, 0) }, 0, 0},
		{"edge: depth b is the axis", func() (float64, error) { return EllipseChordFromEdge(5, 3, 3) }, 10, tolFew},
		{"edge: depth 2b is the far vertex", func() (float64, error) { return EllipseChordFromEdge(5, 3, 6) }, 0, 0},
		{"edge: near the vertex keeps its digits", func() (float64, error) { return EllipseChordFromEdge(1, 1, 1e-10) }, 2 * math.Sqrt(1e-10*(2-1e-10)), tolFew},
		{"edge: 2b would overflow", func() (float64, error) { return EllipseChordFromEdge(1, 1e308, 1e308) }, 2, tolFew},
	})
}

func TestEllipseChord_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"centre NaN", func() error { return e1(EllipseChordFromCenter(nan, 1, 0)) }, ErrNotFinite},
		{"centre Inf offset", func() error { return e1(EllipseChordFromCenter(1, 1, inf)) }, ErrNotFinite},
		{"centre zero semi-axis", func() error { return e1(EllipseChordFromCenter(1, 0, 0)) }, ErrNonPositive},
		{"centre negative offset", func() error { return e1(EllipseChordFromCenter(1, 1, -0.5)) }, ErrNegative},
		{"centre offset beyond b", func() error { return e1(EllipseChordFromCenter(5, 3, 3.5)) }, ErrOutOfRange},
		{"centre overflow", func() error { return e1(EllipseChordFromCenter(1e308, 1, 0)) }, ErrOverflow},
		{"edge NaN", func() error { return e1(EllipseChordFromEdge(1, nan, 0)) }, ErrNotFinite},
		{"edge negative semi-axis", func() error { return e1(EllipseChordFromEdge(-1, 1, 0)) }, ErrNonPositive},
		{"edge negative depth", func() error { return e1(EllipseChordFromEdge(1, 1, -1)) }, ErrNegative},
		{"edge depth beyond 2b", func() error { return e1(EllipseChordFromEdge(5, 3, 6.5)) }, ErrOutOfRange},
		{"edge overflow", func() error { return e1(EllipseChordFromEdge(1e308, 1, 1)) }, ErrOverflow},
	})
}

// Oracle: the chord's end points (±L/2, y) lie on x²/a² + y²/b² = 1, and the two functions agree at
// depth h = b − y.
func TestEllipseChord_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		a, b := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		y := b * rng.Float64() * 0.9
		fromCenter, err := EllipseChordFromCenter(a, b, y)
		if err != nil {
			t.Fatal(err)
		}
		x := fromCenter / 2
		if lhs := x*x/(a*a) + y*y/(b*b); math.Abs(lhs-1) > 1e-14 {
			t.Fatalf("(%v, %v) on the ellipse a %v, b %v: x²/a² + y²/b² = %v", x, y, a, b, lhs)
		}
		if fromEdge, err := EllipseChordFromEdge(a, b, b-y); err != nil || !relNear(fromEdge, fromCenter, 1e-13) {
			t.Fatalf("a %v, b %v, y %v: from centre %v, from edge (%v, %v)", a, b, y, fromCenter, fromEdge, err)
		}
	}
}

func TestEllipseChord_allocs(t *testing.T) {
	assertNoAllocs(t, "EllipseChordFromCenter", func() { sink, _ = EllipseChordFromCenter(5, 3, 1.8) })
	assertNoAllocs(t, "EllipseChordFromEdge", func() { sink, _ = EllipseChordFromEdge(5, 3, 1.2) })
}

func BenchmarkEllipseChordFromCenter(b *testing.B) {
	b.ReportAllocs()
	x, y, off := 5.0, 3.0, 1.8
	for b.Loop() {
		sink, _ = EllipseChordFromCenter(x, y, off)
	}
}

func BenchmarkEllipseChordFromEdge(b *testing.B) {
	b.ReportAllocs()
	x, y, depth := 5.0, 3.0, 1.2
	for b.Loop() {
		sink, _ = EllipseChordFromEdge(x, y, depth)
	}
}
