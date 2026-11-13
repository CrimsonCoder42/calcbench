package geometry

import (
	"math"
	"testing"
)

// Worked example: Euclid, Elements III.36 (tangent² = whole secant × external part), with the secant
// through the centre: t² = (d + r)(d − r). Hand-derived for a point 2 above a sphere of radius 6 371 000:
// t = √(12742002 × 2) ≈ 5048.168380706. Also the 3-4-5 case: r 3, d 5 → 4.
func TestCircleTangentLength_example(t *testing.T) {
	runValCases(t, []valCase{
		{"CircleTangentLength(6371000, 6371002)", func() (float64, error) { return CircleTangentLength(6371000, 6371002) }, math.Sqrt(25484004), tolFew},
		{"CircleTangentLength(3, 5)", func() (float64, error) { return CircleTangentLength(3, 5) }, 4, tolFew},
	})
}

func TestCircleTangentLength_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"point on the circle", func() (float64, error) { return CircleTangentLength(2, 2) }, 0, 0},
		{"just outside: d − r one ulp", func() (float64, error) { return CircleTangentLength(1, 1+0x1p-52) }, math.Sqrt(0x1p-52 * (2 + 0x1p-52)), tolFew},
		{"d + r would overflow", func() (float64, error) { return CircleTangentLength(0.9e308, 1.5e308) }, 1.2e308, tolSome},
	})
}

func TestCircleTangentLength_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(CircleTangentLength(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(CircleTangentLength(1, inf)) }, ErrNotFinite},
		{"zero radius", func() error { return e1(CircleTangentLength(0, 1)) }, ErrNonPositive},
		{"inside the circle", func() error { return e1(CircleTangentLength(2, 1)) }, ErrOutOfRange},
	})
}

// Oracle: √(d² − r²) directly, for points well outside the circle (no cancellation).
func TestCircleTangentLength_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r := 0.1 + rng.Float64()*100
		d := r * (1.5 + rng.Float64()*10)
		want := math.Sqrt(d*d - r*r)
		if got, err := CircleTangentLength(r, d); err != nil || !relNear(got, want, 1e-14) {
			t.Fatalf("CircleTangentLength(%v, %v) = (%v, %v), want %v", r, d, got, err, want)
		}
	}
}

func TestCircleTangentLength_allocs(t *testing.T) {
	assertNoAllocs(t, "CircleTangentLength", func() { sink, _ = CircleTangentLength(3, 5) })
}

func BenchmarkCircleTangentLength(b *testing.B) {
	b.ReportAllocs()
	r, d := 3.0, 5.0
	for b.Loop() {
		sink, _ = CircleTangentLength(r, d)
	}
}
