package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.4.3 (angle of slope 1 is tan⁻¹ 1), hand-derived: π/4.
func TestSlopeAngleRad_example(t *testing.T) {
	runValCases(t, []valCase{
		{"SlopeAngleRad(1, 1)", func() (float64, error) { return SlopeAngleRad(1, 1) }, math.Pi / 4, tolFew},
	})
}

func TestSlopeAngleRad_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"vertical up", func() (float64, error) { return SlopeAngleRad(1, 0) }, math.Pi / 2, tolFew},
		{"vertical down", func() (float64, error) { return SlopeAngleRad(-1, 0) }, -math.Pi / 2, tolFew},
		{"horizontal", func() (float64, error) { return SlopeAngleRad(0, 5) }, 0, 0},
		{"backwards", func() (float64, error) { return SlopeAngleRad(0, -5) }, math.Pi, 0},
		{"backwards with -0 rise is +π", func() (float64, error) { return SlopeAngleRad(math.Copysign(0, -1), -5) }, math.Pi, 0},
		{"30° (rise 1, run √3)", func() (float64, error) { return SlopeAngleRad(1, math.Sqrt(3)) }, math.Pi / 6, tolFew},
		{"third quadrant", func() (float64, error) { return SlopeAngleRad(-1, -1) }, -3 * math.Pi / 4, tolFew},
	})
}

func TestSlopeAngleRad_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN rise", func() error { return e1(SlopeAngleRad(nan, 1)) }, ErrNotFinite},
		{"Inf run", func() error { return e1(SlopeAngleRad(1, -inf)) }, ErrNotFinite},
		{"zero vector", func() error { return e1(SlopeAngleRad(0, 0)) }, ErrDegenerate},
	})
}

// Oracle: (r cos θ, r sin θ) = (run, rise) with r = √(rise² + run²), and θ lies in (−π, π]. The
// comparison is relative to r because sin and cos are well conditioned there (tan is not, near π/2).
func TestSlopeAngleRad_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		rise, run := rng.NormFloat64(), rng.NormFloat64()
		got, err := SlopeAngleRad(rise, run)
		r := math.Hypot(rise, run)
		sin, cos := math.Sincos(got)
		if err != nil || got <= -math.Pi || got > math.Pi ||
			math.Abs(r*cos-run) > tolSome*r || math.Abs(r*sin-rise) > tolSome*r {
			t.Fatalf("SlopeAngleRad(%v, %v) = (%v, %v)", rise, run, got, err)
		}
	}
}

func TestSlopeAngleRad_allocs(t *testing.T) {
	assertNoAllocs(t, "SlopeAngleRad", func() { sink, _ = SlopeAngleRad(1, 2) })
}

func BenchmarkSlopeAngleRad(b *testing.B) {
	b.ReportAllocs()
	rise, run := 1.0, 2.0
	for b.Loop() {
		sink, _ = SlopeAngleRad(rise, run)
	}
}
