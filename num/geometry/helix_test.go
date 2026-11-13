package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.20.1.3 p. 375: for (a cos t, a sin t, bt) the arc length is
// s = t√(a² + b²). With a 2, b 0.5 and one turn (t = 2π, pitch 2πb = π): s = 2π√4.25. Hand-derived:
// r 1, pitch 2, 3 turns → 3√(4π² + 4) ≈ 19.781449856854, and that length back gives height 6.
func TestHelix_example(t *testing.T) {
	runValCases(t, []valCase{
		{"HelixArcLength(2, π, 1), the CRC parametrisation", func() (float64, error) { return HelixArcLength(2, math.Pi, 1) }, 2 * math.Pi * math.Sqrt(4.25), tolFew},
		{"HelixArcLength(1, 2, 3)", func() (float64, error) { return HelixArcLength(1, 2, 3) }, 3 * math.Sqrt(4*math.Pi*math.Pi+4), tolFew},
		{"HelixHeight(1, 3√(4π² + 4), 3)", func() (float64, error) { return HelixHeight(1, 3*math.Sqrt(4*math.Pi*math.Pi+4), 3) }, 6, 1e-13}, // tolerance: (L/h)² ≈ 11 times a few ulps of L
	})
}

func TestHelix_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"fractional turns", func() (float64, error) { return HelixArcLength(1, 2, 0.5) }, 0.5 * math.Sqrt(4*math.Pi*math.Pi+4), tolFew},
		{"tiny pitch tends to the circle", func() (float64, error) { return HelixArcLength(1, 1e-12, 1) }, 2 * math.Pi, tolFew},
		{"2πr would overflow, length finite", func() (float64, error) { return HelixArcLength(1e308, 1, 1e-300) }, 2 * math.Pi * 1e8, tolSome},
		{"flat coil: length exactly 2πrn", func() (float64, error) { return HelixHeight(1, 2*math.Pi, 1) }, 0, 0},
	})
}

func TestHelix_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"length NaN radius", func() error { return e1(HelixArcLength(nan, 1, 1)) }, ErrNotFinite},
		{"length Inf pitch", func() error { return e1(HelixArcLength(1, inf, 1)) }, ErrNotFinite},
		{"length zero pitch", func() error { return e1(HelixArcLength(1, 0, 1)) }, ErrNonPositive},
		{"length zero turns", func() error { return e1(HelixArcLength(1, 1, 0)) }, ErrNonPositive},
		{"length overflow", func() error { return e1(HelixArcLength(1e308, 1, 1)) }, ErrOverflow},
		{"height NaN", func() error { return e1(HelixHeight(1, nan, 1)) }, ErrNotFinite},
		{"height negative radius", func() error { return e1(HelixHeight(-1, 10, 1)) }, ErrNonPositive},
		{"height zero turns", func() error { return e1(HelixHeight(1, 10, 0)) }, ErrNonPositive},
		{"height too short for the turns", func() error { return e1(HelixHeight(1, 6, 1)) }, ErrOutOfRange},
		{"height: 2πrn beyond float64", func() error { return e1(HelixHeight(1e308, 1e308, 10)) }, ErrOutOfRange},
	})
}

// Oracle: the length and the height are inverse (pitch-to-circumference ratios in [0.1, 10], where the
// rounded length costs only a few ulps of the height).
func TestHelix_roundTrip(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r := 0.1 + rng.Float64()*10
		pitch := 2 * math.Pi * r * math.Pow(10, 2*rng.Float64()-1)
		turns := 0.1 + rng.Float64()*20
		length, err := HelixArcLength(r, pitch, turns)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := HelixHeight(r, length, turns); err != nil || !relNear(got, pitch*turns, 1e-12) {
			t.Fatalf("HelixHeight(%v, %v, %v) = (%v, %v), want %v", r, length, turns, got, err, pitch*turns)
		}
	}
}

func TestHelix_allocs(t *testing.T) {
	assertNoAllocs(t, "HelixArcLength", func() { sink, _ = HelixArcLength(1, 2, 3) })
	assertNoAllocs(t, "HelixHeight", func() { sink, _ = HelixHeight(1, 20, 3) })
}

func BenchmarkHelixArcLength(b *testing.B) {
	b.ReportAllocs()
	r, p, n := 1.0, 2.0, 3.0
	for b.Loop() {
		sink, _ = HelixArcLength(r, p, n)
	}
}

func BenchmarkHelixHeight(b *testing.B) {
	b.ReportAllocs()
	r, l, n := 1.0, 20.0, 3.0
	for b.Loop() {
		sink, _ = HelixHeight(r, l, n)
	}
}
