package geometry

import (
	"math"
	"testing"
)

// Worked example: OpenStax Prealgebra 2e §9.3 (supplementary angles add to 180°), hand-derived:
// 60° → 120°, i.e. π/3 → 2π/3.
func TestSupplementaryAngleRad_example(t *testing.T) {
	runValCases(t, []valCase{
		{"SupplementaryAngleRad(π/3)", func() (float64, error) { return SupplementaryAngleRad(math.Pi / 3) }, 2 * math.Pi / 3, tolFew},
	})
}

func TestSupplementaryAngleRad_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"zero", func() (float64, error) { return SupplementaryAngleRad(0) }, math.Pi, 0},
		{"straight angle", func() (float64, error) { return SupplementaryAngleRad(math.Pi) }, 0, 0},
		{"right angle", func() (float64, error) { return SupplementaryAngleRad(math.Pi / 2) }, math.Pi / 2, 0},
	})
}

func TestSupplementaryAngleRad_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(SupplementaryAngleRad(nan)) }, ErrNotFinite},
		{"-Inf", func() error { return e1(SupplementaryAngleRad(-inf)) }, ErrNotFinite},
		{"negative", func() error { return e1(SupplementaryAngleRad(-0.1)) }, ErrOutOfRange},
		{"reflex", func() error { return e1(SupplementaryAngleRad(4)) }, ErrOutOfRange},
	})
}

// Property: θ + (π − θ) = π to within an ulp of π, over the whole range.
func TestSupplementaryAngleRad_sum(t *testing.T) {
	rng := newRNG()
	for range 400 {
		theta := rng.Float64() * math.Pi
		got, err := SupplementaryAngleRad(theta)
		if err != nil || math.Abs(got+theta-math.Pi) > 8e-16 {
			t.Fatalf("SupplementaryAngleRad(%v) = (%v, %v)", theta, got, err)
		}
	}
}

func TestSupplementaryAngleRad_allocs(t *testing.T) {
	assertNoAllocs(t, "SupplementaryAngleRad", func() { sink, _ = SupplementaryAngleRad(0.5) })
}

func BenchmarkSupplementaryAngleRad(b *testing.B) {
	b.ReportAllocs()
	theta := 0.5
	for b.Loop() {
		sink, _ = SupplementaryAngleRad(theta)
	}
}
