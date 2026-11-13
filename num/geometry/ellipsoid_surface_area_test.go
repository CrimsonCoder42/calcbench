package geometry

import (
	"math"
	"testing"
)

// Carlson 1995, §3 "Numerical checks": real-argument values for R_D and R_G, printed to 14 significant
// digits (hence 1e-13).
func TestCarlson_checkValues(t *testing.T) {
	rd := func(x, y, z float64) func() (float64, error) {
		return func() (float64, error) { v, _ := carlsonRD(x, y, z); return v, nil }
	}
	rg := func(x, y, z float64) func() (float64, error) {
		return func() (float64, error) { v, _ := carlsonRG(x, y, z); return v, nil }
	}
	runValCases(t, []valCase{
		{"R_D(0, 2, 1)", rd(0, 2, 1), 1.7972103521034, 1e-13},
		{"R_D(2, 3, 4)", rd(2, 3, 4), 0.16510527294261, 1e-13},
		{"R_G(2, 3, 4)", rg(2, 3, 4), 1.7255030280692, 1e-13},
		{"R_G(0, 16, 16) = π", rg(0, 16, 16), math.Pi, 1e-13},
		{"R_G(0, 0.0796, 4) = E(0.99)", rg(0, 0.0796, 4), 1.0284758090288, 1e-13},
		{"R_G(0, 0, z) = √z/2", rg(0, 0, 9), 1.5, 0},
	})
	if _, ok := carlsonRD(math.NaN(), 1, 1); ok {
		t.Error("carlsonRD(NaN, 1, 1) reported convergence; the iteration cap is not enforced")
	}
}

// Worked examples: DLMF 19.33.1 evaluated at 60 digits with R_G from DLMF 19.21.10 (a route independent of
// the code's 19.21.11), hand-derived: (3, 2, 1) → 48.88214630258206; (2, 1, 1) → 21.478435327883737, the
// prolate-spheroid closed form CRC 31st ed. eq. (4.18.1); (2, 2, 1) → 34.68753081338021; the sphere 4πr².
func TestEllipsoidSurfaceArea_example(t *testing.T) {
	runValCases(t, []valCase{
		{"EllipsoidSurfaceArea(3, 2, 1)", func() (float64, error) { return EllipsoidSurfaceArea(3, 2, 1) }, 48.88214630258206, tolFew},
		{"EllipsoidSurfaceArea(2, 1, 1)", func() (float64, error) { return EllipsoidSurfaceArea(2, 1, 1) }, 21.478435327883737, tolFew},
		{"EllipsoidSurfaceArea(2, 2, 1)", func() (float64, error) { return EllipsoidSurfaceArea(2, 2, 1) }, 34.68753081338021, tolFew},
		{"sphere", func() (float64, error) { return EllipsoidSurfaceArea(1, 1, 1) }, 4 * math.Pi, tolFew},
	})
}

func TestEllipsoidSurfaceArea_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"any order", func() (float64, error) { return EllipsoidSurfaceArea(1, 3, 2) }, 48.88214630258206, tolFew},
		{"thin plate", func() (float64, error) { return EllipsoidSurfaceArea(5, 3, 0.01) }, 94.25245185439326, tolFew},
		{"flat limit: both ratios underflow, 2πab", func() (float64, error) { return EllipsoidSurfaceArea(1, 1, 1e-200) }, 2 * math.Pi, tolFew},
		{"needle: a subnormal ratio, π²ab in the limit", func() (float64, error) { return EllipsoidSurfaceArea(1e160, 1, 1) }, math.Pi * math.Pi * 1e160, tolSome},
	})
}

func TestEllipsoidSurfaceArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(EllipsoidSurfaceArea(nan, 1, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(EllipsoidSurfaceArea(1, inf, 1)) }, ErrNotFinite},
		{"zero", func() error { return e1(EllipsoidSurfaceArea(1, 1, 0)) }, ErrNonPositive},
		{"overflow", func() error { return e1(EllipsoidSurfaceArea(1e200, 1e200, 1)) }, ErrOverflow},
	})
}

// spheroidArea is the CRC 31st ed. eq. (4.18.1) closed form for x²/a² + (y² + z²)/b² = 1: prolate (a > b)
// with arcsin, oblate (a < b) with the logarithm.
func spheroidArea(a, b float64) float64 {
	if a > b {
		e := math.Sqrt(a*a - b*b)
		return 2*math.Pi*b*b + 2*math.Pi*a*a*b/e*math.Asin(e/a)
	}
	e := math.Sqrt(b*b - a*a)
	return 2*math.Pi*b*b + math.Pi*a*a*b/e*math.Log((b+e)/(b-e))
}

// Oracle: the spheroid closed forms on random prolate and oblate spheroids (ratios kept 10% away from the
// sphere, where the closed forms themselves cancel).
func TestEllipsoidSurfaceArea_vsSpheroid(t *testing.T) {
	rng := newRNG()
	for range 400 {
		b := 0.1 + rng.Float64()*10
		ratio := 0.05 + 0.85*rng.Float64()
		if rng.IntN(2) == 0 {
			ratio = 1 / ratio
		}
		a := b * ratio
		want := spheroidArea(a, b)
		if got, err := EllipsoidSurfaceArea(a, b, b); err != nil || !relNear(got, want, 1e-12) {
			t.Fatalf("EllipsoidSurfaceArea(%v, %v, %v) = (%v, %v), spheroid %v", a, b, b, got, err, want)
		}
	}
}

func TestEllipsoidSurfaceArea_allocs(t *testing.T) {
	assertNoAllocs(t, "EllipsoidSurfaceArea", func() { sink, _ = EllipsoidSurfaceArea(3, 2, 1) })
}

func BenchmarkEllipsoidSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	x, y, z := 3.0, 2.0, 1.0
	for b.Loop() {
		sink, _ = EllipsoidSurfaceArea(x, y, z)
	}
}
