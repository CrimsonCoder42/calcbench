package geometry

import (
	"math"
	"testing"
)

// Worked example: Ramanujan 1914 §15 (second approximation), evaluated at 60 digits, hand-derived:
// (2, 1) → 9.688448216130084 against the exact 9.688448220547676 (relative error −4.6e-10); (5, 3) →
// 25.52699886278876.
func TestEllipsePerimeterRamanujan_example(t *testing.T) {
	runValCases(t, []valCase{
		{"EllipsePerimeterRamanujan(2, 1)", func() (float64, error) { return EllipsePerimeterRamanujan(2, 1) }, 9.688448216130084, tolFew},
		{"EllipsePerimeterRamanujan(5, 3)", func() (float64, error) { return EllipsePerimeterRamanujan(5, 3) }, 25.52699886278876, tolFew},
	})
}

func TestEllipsePerimeterRamanujan_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"circle is exact", func() (float64, error) { return EllipsePerimeterRamanujan(1, 1) }, 2 * math.Pi, tolFew},
		{"either order", func() (float64, error) { return EllipsePerimeterRamanujan(1, 2) }, 9.688448216130084, tolFew},
		{"flat limit: π(1 + 3/(10 + 1)) = 14π/11", func() (float64, error) { return EllipsePerimeterRamanujan(1, 1e-300) }, 14 * math.Pi / 11, tolFew},
	})
}

func TestEllipsePerimeterRamanujan_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(EllipsePerimeterRamanujan(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(EllipsePerimeterRamanujan(1, -inf)) }, ErrNotFinite},
		{"negative", func() error { return e1(EllipsePerimeterRamanujan(-1, 1)) }, ErrNonPositive},
		{"overflow", func() error { return e1(EllipsePerimeterRamanujan(1e308, 1e308)) }, ErrOverflow},
	})
}

// Oracle against the exact perimeter: Ramanujan underestimates, by at most 4.1e-4 relative (as b/a → 0),
// and for λ = (a−b)/(a+b) ≤ 1/3 the shortfall over π(a+b) is the DLMF 19.9.10 term 3λ¹⁰/2¹⁷ times a factor
// in [1, 1.21] (computed at 60 digits). The factor is checked only where the term exceeds 1e-11, so float
// rounding (about 1e-15) cannot disturb it.
func TestEllipsePerimeterRamanujan_vsExact(t *testing.T) {
	rng := newRNG()
	for range 400 {
		a := 0.1 + rng.Float64()*10
		b := a * math.Pow(10, -9*rng.Float64())
		exact, _ := EllipsePerimeter(a, b)
		approx, err := EllipsePerimeterRamanujan(a, b)
		if err != nil {
			t.Fatal(err)
		}
		rel := (exact - approx) / exact
		if rel < -1e-15 || rel > 4.1e-4 {
			t.Fatalf("a %v, b %v: relative shortfall %v", a, b, rel)
		}
		lambda := (a - b) / (a + b)
		term := 3 * math.Pow(lambda, 10) / (1 << 17)
		if lambda <= 1.0/3 && term > 1e-11 {
			factor := (exact - approx) / (math.Pi * (a + b)) / term
			if factor < 0.99 || factor > 1.25 {
				t.Fatalf("a %v, b %v: shortfall / (3λ¹⁰/2¹⁷) = %v, want in [1, 1.21]", a, b, factor)
			}
		}
	}
}

func TestEllipsePerimeterRamanujan_allocs(t *testing.T) {
	assertNoAllocs(t, "EllipsePerimeterRamanujan", func() { sink, _ = EllipsePerimeterRamanujan(5, 3) })
}

func BenchmarkEllipsePerimeterRamanujan(b *testing.B) {
	b.ReportAllocs()
	x, y := 5.0, 3.0
	for b.Loop() {
		sink, _ = EllipsePerimeterRamanujan(x, y)
	}
}
