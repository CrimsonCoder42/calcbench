package special

import (
	"math"
	"testing"
)

// The direct form exp(a·ln x − x − ln Γ(a)) is accurate to a few ulps of a·ln x for moderate a, so it
// is a reference for the saddle-point branch (a > 15) there. Tolerance 1e-11: at a = 501, a·ln x ≈ 3·10³,
// whose ulp (4.5e-13) is the direct form's own relative error.
func TestGammaPrefactor_SaddlePointMatchesDirectForm(t *testing.T) {
	for _, a := range []float64{15.5, 16, 20.25, 37, 81.5, 120, 501} {
		for _, r := range []float64{0.5, 0.9, 1, 1.1, 2} {
			x := r * a
			lg, _ := math.Lgamma(a)
			want := math.Exp(a*math.Log(x) - x - lg)
			if got := gammaPrefactor(a, x); math.Abs(got-want) > 1e-11*want {
				t.Errorf("gammaPrefactor(%v, %v) = %v; want %v", a, x, got, want)
			}
		}
	}
}

// DLMF 8.8.5: P(a+1, x) = P(a, x) − t(a), t(a) = x^a e^(−x)/Γ(a+1), and t(a+1) = t(a)·x/(a+1). So with
// d₀ = P(a, x) − P(a+1, x) and d₁ = P(a+1, x) − P(a+2, x), d₁ = d₀·x/(a+1): an identity among P values
// alone, with no prefactor computed in the test. The same holds for Q with the signs reversed. A prefactor
// formed from ln Γ(a) carries a relative error near 1e-9 at a = 5e5 (ulp of ln Γ ≈ 6·10⁶), which moves
// d₁ by about 1e-6 of itself; tolerance 1e-10 of d₁ leaves room for the few-ulp rounding of the sums
// amplified by P/t ≤ about √a.
func TestRegularizedGamma_LargeShapeRecurrence(t *testing.T) {
	for _, a := range []float64{1000.5, 1e4, 100000.25, 5e5} {
		s := math.Sqrt(a)
		// Lower side (x < a+1: the series branch for a, a+1, a+2) uses P; upper side (x ≥ a+3: the
		// continued fraction for all three) uses Q.
		for _, x := range []float64{a - 3*s, a - s, a - 0.5, a + 3, a + s, a + 3*s} {
			f := RegularizedGammaP
			sign := 1.0
			if x >= a+1 {
				f, sign = RegularizedGammaQ, -1
			}
			v0, e0 := f(a, x)
			v1, e1 := f(a+1, x)
			v2, e2 := f(a+2, x)
			if e0 != nil || e1 != nil || e2 != nil {
				t.Fatalf("a = %v, x = %v: errors %v, %v, %v", a, x, e0, e1, e2)
			}
			d0, d1 := sign*(v0-v1), sign*(v1-v2)
			if d1 <= 0 || math.Abs(d1-d0*x/(a+1)) > 1e-10*d1 {
				t.Errorf("a = %v, x = %v: d1 = %v, d0·x/(a+1) = %v", a, x, d1, d0*x/(a+1))
			}
		}
	}
}

// DLMF 8.2.5: P + Q = 1, at large a on both branches.
func TestRegularizedGamma_LargeShapeComplement(t *testing.T) {
	for _, a := range []float64{1e3, 1e5, 5e5} {
		for _, x := range []float64{0.99 * a, a, 1.01 * a} {
			p, _ := RegularizedGammaP(a, x)
			q, _ := RegularizedGammaQ(a, x)
			if math.Abs(p+q-1) > 4e-16 {
				t.Errorf("P + Q at a = %v, x = %v = %v; want 1", a, x, p+q)
			}
		}
	}
}

func BenchmarkRegularizedGammaP_LargeShape(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = RegularizedGammaP(5e5, 5e5-700)
	}
}
