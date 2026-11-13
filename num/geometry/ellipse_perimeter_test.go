package geometry

import (
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Worked examples: DLMF 19.9.9 (L = 4aE(k)) evaluated at 60 digits by Carlson's AGM (eqs. (42)–(46)),
// hand-derived and agreeing with the catalog's quadrature: (2, 1) → 9.688448220547676,
// (5, 3) → 25.526998863398128, (3, 1) → 13.364893220555258. The circle gives 2πa.
func TestEllipsePerimeter_example(t *testing.T) {
	runValCases(t, []valCase{
		{"EllipsePerimeter(2, 1)", func() (float64, error) { return EllipsePerimeter(2, 1) }, 9.688448220547676, tolFew},
		{"EllipsePerimeter(5, 3)", func() (float64, error) { return EllipsePerimeter(5, 3) }, 25.526998863398128, tolFew},
		{"EllipsePerimeter(3, 1)", func() (float64, error) { return EllipsePerimeter(3, 1) }, 13.364893220555258, tolFew},
		{"circle", func() (float64, error) { return EllipsePerimeter(1, 1) }, 2 * math.Pi, tolFew},
	})
}

func TestEllipsePerimeter_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"either order", func() (float64, error) { return EllipsePerimeter(1, 2) }, 9.688448220547676, tolFew},
		{"flat, AGM branch (b/a = 1e-9): 4 + 4.3e-17", func() (float64, error) { return EllipsePerimeter(1, 1e-9) }, 4, tolFew},
		{"flat, limit branch (b/a = 1e-300)", func() (float64, error) { return EllipsePerimeter(1, 1e-300) }, 4, 0},
		{"smallest subnormal ratio", func() (float64, error) { return EllipsePerimeter(1, math.SmallestNonzeroFloat64) }, 4, 0},
		{"huge axes, no overflow", func() (float64, error) { return EllipsePerimeter(2e300, 1e300) }, 9.688448220547676e300, tolFew},
	})
}

func TestEllipsePerimeter_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(EllipsePerimeter(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(EllipsePerimeter(1, inf)) }, ErrNotFinite},
		{"zero", func() error { return e1(EllipsePerimeter(1, 0)) }, ErrNonPositive},
		{"overflow", func() error { return e1(EllipsePerimeter(1e308, 1)) }, ErrOverflow},
	})
	// ErrNoConvergence cannot be reached from valid input; the cap itself is exercised here: a NaN never
	// meets the stopping rule, so the loop must end at agmMaxSteps and report failure.
	if _, ok := completeEllipticE(math.NaN()); ok {
		t.Error("completeEllipticE(NaN) reported convergence; the iteration cap is not enforced")
	}
}

// ellipsePerimeterQuadrature is the trapezoid rule on L = ∫₀^2π √(a² sin²t + b² cos²t) dt. For a periodic
// analytic integrand it converges geometrically, like e^(−N·artanh(b/a)), so N = 4096 is exact to rounding
// for b/a ≥ 0.05. The terms are summed with compensation so rounding stays near one ulp.
func ellipsePerimeterQuadrature(a, b float64) float64 {
	const n = 4096
	var sum, comp float64
	for i := range n {
		s, c := math.Sincos(2 * math.Pi * float64(i) / n)
		sum, comp = reduce.NeumaierAdd(sum, comp, math.Hypot(a*s, b*c))
	}
	return (sum + comp) * (2 * math.Pi / n)
}

// Oracle: numerical integration of the arc length (an independent route to the same number).
func TestEllipsePerimeter_oracle(t *testing.T) {
	rng := newRNG()
	for range 200 {
		a := 0.1 + rng.Float64()*10
		b := a * (0.05 + 0.95*rng.Float64())
		want := ellipsePerimeterQuadrature(a, b)
		if got, err := EllipsePerimeter(a, b); err != nil || !relNear(got, want, 1e-13) {
			t.Fatalf("EllipsePerimeter(%v, %v) = (%v, %v), quadrature %v", a, b, got, err, want)
		}
	}
}

func TestEllipsePerimeter_allocs(t *testing.T) {
	assertNoAllocs(t, "EllipsePerimeter", func() { sink, _ = EllipsePerimeter(5, 3) })
}

func BenchmarkEllipsePerimeter(b *testing.B) {
	b.ReportAllocs()
	x, y := 5.0, 3.0
	for b.Loop() {
		sink, _ = EllipsePerimeter(x, y)
	}
}
