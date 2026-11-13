package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.6.6 p. 334 (c = √(4h(2R − h))), hand-derived: chord 8, sagitta 2 →
// R 5 (half-chord 4 and R − h = 3 make a 3-4-5 triangle).
func TestArcRadiusFromChordSagitta_example(t *testing.T) {
	runValCases(t, []valCase{
		{"ArcRadiusFromChordSagitta(8, 2)", func() (float64, error) { return ArcRadiusFromChordSagitta(8, 2) }, 5, 0},
	})
}

func TestArcRadiusFromChordSagitta_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"semicircle: sagitta = chord/2", func() (float64, error) { return ArcRadiusFromChordSagitta(10, 5) }, 5, 0},
		{"major arc: sagitta > chord/2", func() (float64, error) { return ArcRadiusFromChordSagitta(8, 8) }, 5, 0},
		{"flat arch, c² would overflow", func() (float64, error) { return ArcRadiusFromChordSagitta(1e200, 1e195) }, 1.25e204 + 5e194, tolSome},
		{"tiny sagitta, 1/s would overflow", func() (float64, error) {
			return ArcRadiusFromChordSagitta(1e-15, math.SmallestNonzeroFloat64)
		}, 1e-30 / (8 * math.SmallestNonzeroFloat64), tolSome},
	})
}

func TestArcRadiusFromChordSagitta_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(ArcRadiusFromChordSagitta(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(ArcRadiusFromChordSagitta(1, inf)) }, ErrNotFinite},
		{"zero chord", func() error { return e1(ArcRadiusFromChordSagitta(0, 1)) }, ErrNonPositive},
		{"negative sagitta", func() error { return e1(ArcRadiusFromChordSagitta(1, -1)) }, ErrNonPositive},
		{"overflow", func() error { return e1(ArcRadiusFromChordSagitta(1e200, 1e-200)) }, ErrOverflow},
	})
}

// Oracle: build the chord of a circle from R and a sagitta h by the CRC relation c = 2√(h(2R − h)), then
// recover R.
func TestArcRadiusFromChordSagitta_roundTrip(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r := 0.1 + rng.Float64()*100
		h := r * (0.001 + 1.998*rng.Float64())
		c := 2 * math.Sqrt(h*(2*r-h))
		if got, err := ArcRadiusFromChordSagitta(c, h); err != nil || !relNear(got, r, 1e-14) {
			t.Fatalf("ArcRadiusFromChordSagitta(%v, %v) = (%v, %v), want %v", c, h, got, err, r)
		}
	}
}

func TestArcRadiusFromChordSagitta_allocs(t *testing.T) {
	assertNoAllocs(t, "ArcRadiusFromChordSagitta", func() { sink, _ = ArcRadiusFromChordSagitta(8, 2) })
}

func BenchmarkArcRadiusFromChordSagitta(b *testing.B) {
	b.ReportAllocs()
	c, h := 8.0, 2.0
	for b.Loop() {
		sink, _ = ArcRadiusFromChordSagitta(c, h)
	}
}
