package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.6.6 p. 334 (area πr² of each circle), hand-derived: 25π − 9π = 16π.
func TestAnnulusArea_example(t *testing.T) {
	runValCases(t, []valCase{
		{"AnnulusArea(5, 3)", func() (float64, error) { return AnnulusArea(5, 3) }, 16 * math.Pi, tolFew},
	})
}

func TestAnnulusArea_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"thin ring keeps its digits", func() (float64, error) { return AnnulusArea(1+0x1p-40, 1) }, math.Pi * 0x1p-40 * (2 + 0x1p-40), tolFew},
		{"large radii, no overflow of R²", func() (float64, error) { return AnnulusArea(5e153, 4e153) }, math.Pi * 9e306, tolSome},
	})
}

func TestAnnulusArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(AnnulusArea(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(AnnulusArea(2, inf)) }, ErrNotFinite},
		{"zero inner", func() error { return e1(AnnulusArea(2, 0)) }, ErrNonPositive},
		{"inner equals outer", func() error { return e1(AnnulusArea(2, 2)) }, ErrOutOfRange},
		{"inner larger", func() error { return e1(AnnulusArea(2, 3)) }, ErrOutOfRange},
		{"overflow", func() error { return e1(AnnulusArea(1e200, 1)) }, ErrOverflow},
	})
}

// Oracle: the difference of two circle areas, on rings that are not thin (no cancellation there).
func TestAnnulusArea_vsCircles(t *testing.T) {
	rng := newRNG()
	for range 400 {
		outer := 1 + rng.Float64()*100
		inner := outer * (0.01 + 0.8*rng.Float64())
		a, _ := CircleArea(outer)
		b, _ := CircleArea(inner)
		if got, err := AnnulusArea(outer, inner); err != nil || !relNear(got, a-b, 1e-14) {
			t.Fatalf("AnnulusArea(%v, %v) = (%v, %v), want %v", outer, inner, got, err, a-b)
		}
	}
}

func TestAnnulusArea_allocs(t *testing.T) {
	assertNoAllocs(t, "AnnulusArea", func() { sink, _ = AnnulusArea(5, 3) })
}

func BenchmarkAnnulusArea(b *testing.B) {
	b.ReportAllocs()
	outer, inner := 5.0, 3.0
	for b.Loop() {
		sink, _ = AnnulusArea(outer, inner)
	}
}
