package signal_test

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
	"github.com/CrimsonCoder42/calcbench/num/signal"
)

// TestConvolve_matchesPolynomialMultiply: mode full is the coefficient product of two polynomials
// (NUM-ALGEB-003, Cormen et al. 3rd ed. §30.1). Tolerance 1e-12·√(n·m): the two direct sums may add
// terms in different orders, and inputs of length ≥ 64 take the FFT path (error ≈ ε·log L·‖a‖₂‖v‖₂).
// This external test package is the only place num/signal meets num/algebra (test-only import).
func TestConvolve_matchesPolynomialMultiply(t *testing.T) {
	r := rand.New(rand.NewPCG(11, 12))
	for _, p := range [][2]int{{1, 1}, {2, 2}, {3, 7}, {40, 50}, {64, 64}, {200, 90}} {
		a := make([]float64, p[0])
		v := make([]float64, p[1])
		for i := range a {
			a[i] = 2*r.Float64() - 1
		}
		for i := range v {
			v[i] = 2*r.Float64() - 1
		}
		want, err := algebra.PolynomialMultiply(a, v)
		if err != nil {
			t.Fatal(err)
		}
		got, err := signal.Convolve(nil, a, v, signal.ModeFull)
		if err != nil {
			t.Fatal(err)
		}
		tol := 1e-12 * math.Sqrt(float64(p[0]*p[1]))
		for i := range want {
			if math.Abs(got[i]-want[i]) > tol {
				t.Fatalf("%dx%d i=%d: %v vs %v", p[0], p[1], i, got[i], want[i])
			}
		}
	}
}
