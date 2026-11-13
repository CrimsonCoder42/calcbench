package geometry

import (
	"math"
	"testing"
)

// Worked example: OpenStax Prealgebra 2e §9.3 (complementary angles add to 90°), hand-derived:
// 30° → 60°, i.e. π/6 → π/3.
func TestComplementaryAngleRad_example(t *testing.T) {
	runValCases(t, []valCase{
		{"ComplementaryAngleRad(π/6)", func() (float64, error) { return ComplementaryAngleRad(math.Pi / 6) }, math.Pi / 3, tolFew},
	})
}

func TestComplementaryAngleRad_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"zero", func() (float64, error) { return ComplementaryAngleRad(0) }, math.Pi / 2, 0},
		{"right angle", func() (float64, error) { return ComplementaryAngleRad(math.Pi / 2) }, 0, 0},
		{"half a right angle", func() (float64, error) { return ComplementaryAngleRad(math.Pi / 4) }, math.Pi / 4, 0},
	})
}

func TestComplementaryAngleRad_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(ComplementaryAngleRad(nan)) }, ErrNotFinite},
		{"Inf", func() error { return e1(ComplementaryAngleRad(inf)) }, ErrNotFinite},
		{"negative", func() error { return e1(ComplementaryAngleRad(-1e-300)) }, ErrOutOfRange},
		{"obtuse", func() error { return e1(ComplementaryAngleRad(math.Nextafter(math.Pi/2, 4))) }, ErrOutOfRange},
	})
}

// Property: θ + (π/2 − θ) = π/2 to within an ulp of π/2, over the whole range.
func TestComplementaryAngleRad_sum(t *testing.T) {
	rng := newRNG()
	for range 400 {
		theta := rng.Float64() * math.Pi / 2
		got, err := ComplementaryAngleRad(theta)
		if err != nil || math.Abs(got+theta-math.Pi/2) > 4e-16 {
			t.Fatalf("ComplementaryAngleRad(%v) = (%v, %v)", theta, got, err)
		}
	}
}

func TestComplementaryAngleRad_allocs(t *testing.T) {
	assertNoAllocs(t, "ComplementaryAngleRad", func() { sink, _ = ComplementaryAngleRad(0.5) })
}

func BenchmarkComplementaryAngleRad(b *testing.B) {
	b.ReportAllocs()
	theta := 0.5
	for b.Loop() {
		sink, _ = ComplementaryAngleRad(theta)
	}
}
