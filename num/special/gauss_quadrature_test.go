package special

import (
	"errors"
	"math"
	"testing"
)

type gaussRuleFunc func(int) ([]float64, []float64, error)

func checkRule(t *testing.T, name string, f gaussRuleFunc, n int, wantNodes, wantWeights []float64, tol float64) {
	t.Helper()
	x, w, err := f(n)
	if err != nil || len(x) != n || len(w) != n {
		t.Fatalf("%s(%d): len %d, %d, err %v", name, n, len(x), len(w), err)
	}
	for i := range x {
		if math.Abs(x[i]-wantNodes[i]) > tol || math.Abs(w[i]-wantWeights[i]) > tol {
			t.Errorf("%s(%d)[%d] = (%v, %v); want (%v, %v)", name, n, i, x[i], w[i], wantNodes[i], wantWeights[i])
		}
	}
}

func TestGaussRules_SmallTables(t *testing.T) {
	// A&S Table 25.4 (Legendre), 25.9 (Laguerre), 25.10 (Hermite); closed forms for n ≤ 3.
	// Tolerance 1e-14: nodes are Newton-polished to within an ulp; weights from a short sum.
	r := math.Sqrt(0.6)
	checkRule(t, "GaussLegendre", GaussLegendre, 1, []float64{0}, []float64{2}, 1e-14)
	checkRule(t, "GaussLegendre", GaussLegendre, 2, []float64{-1 / math.Sqrt(3), 1 / math.Sqrt(3)}, []float64{1, 1}, 1e-14)
	checkRule(t, "GaussLegendre", GaussLegendre, 3, []float64{-r, 0, r}, []float64{5.0 / 9, 8.0 / 9, 5.0 / 9}, 1e-14)
	h := math.Sqrt(math.Pi) / 2
	checkRule(t, "GaussHermite", GaussHermite, 2, []float64{-1 / math.Sqrt2, 1 / math.Sqrt2}, []float64{h, h}, 1e-14)
	he := math.Sqrt(2*math.Pi) / 2
	checkRule(t, "GaussHermiteE", GaussHermiteE, 2, []float64{-1, 1}, []float64{he, he}, 1e-14)
	checkRule(t, "GaussLaguerre", GaussLaguerre, 2, []float64{2 - math.Sqrt2, 2 + math.Sqrt2},
		[]float64{(2 + math.Sqrt2) / 4, (2 - math.Sqrt2) / 4}, 1e-14)
	checkRule(t, "GaussChebyshevT", GaussChebyshevT, 2, []float64{-1 / math.Sqrt2, 1 / math.Sqrt2}, []float64{math.Pi / 2, math.Pi / 2}, 1e-15)
	checkRule(t, "GaussChebyshevU", GaussChebyshevU, 1, []float64{0}, []float64{math.Pi / 2}, 1e-15)
}

// Each n-point rule integrates x^k exactly for k ≤ 2n−1 (DLMF 3.5.15). Moments: Legendre 2/(k+1) (k
// even); Hermite Γ((k+1)/2); HermiteE √(2π)(k−1)!!; Laguerre k!; Chebyshev T π·C(k, k/2)/2^k;
// Chebyshev U π·C(k, k/2)/(2^k (k+2)) — all zero for odd k except Laguerre.
func TestGaussRules_Exactness(t *testing.T) {
	gammaHalf := func(k int) float64 { return math.Gamma(float64(k+1) / 2) }
	binomMid := func(k int) float64 {
		v := 1.0
		for j := 1; j <= k/2; j++ {
			v = v * float64(k/2+j) / float64(j)
		}
		return v / math.Pow(2, float64(k))
	}
	rules := []struct {
		name   string
		f      gaussRuleFunc
		moment func(k int) float64
		odd0   bool
	}{
		{"GaussLegendre", GaussLegendre, func(k int) float64 { return 2 / float64(k+1) }, true},
		{"GaussHermite", GaussHermite, gammaHalf, true},
		{"GaussHermiteE", GaussHermiteE, func(k int) float64 {
			v := math.Sqrt(2 * math.Pi)
			for j := k - 1; j > 1; j -= 2 {
				v *= float64(j)
			}
			return v
		}, true},
		{"GaussLaguerre", GaussLaguerre, func(k int) float64 { return math.Gamma(float64(k + 1)) }, false},
		{"GaussChebyshevT", GaussChebyshevT, func(k int) float64 { return math.Pi * binomMid(k) }, true},
		{"GaussChebyshevU", GaussChebyshevU, func(k int) float64 { return math.Pi * binomMid(k) / float64(k+2) }, true},
	}
	for _, r := range rules {
		for _, n := range []int{1, 4, 9, 16} {
			x, w, err := r.f(n)
			if err != nil {
				t.Fatalf("%s(%d): %v", r.name, n, err)
			}
			for i := 1; i < n; i++ {
				if x[i] <= x[i-1] {
					t.Errorf("%s(%d): nodes not ascending at %d", r.name, n, i)
				}
			}
			for k := 0; k <= 2*n-1; k++ {
				sum, scale := 0.0, 0.0
				for i := range x {
					term := w[i] * math.Pow(x[i], float64(k))
					sum += term
					scale += math.Abs(term)
				}
				want := 0.0
				if !r.odd0 || k%2 == 0 {
					want = r.moment(k)
				}
				// Tolerance 1e-12 of the absolute sum: cancellation in odd moments is bounded by it.
				if math.Abs(sum-want) > 1e-12*math.Max(scale, 1) {
					t.Errorf("%s(%d) moment %d = %v; want %v", r.name, n, k, sum, want)
				}
			}
		}
	}
}

// The nodes are the roots of the degree-n polynomial (checked with the evaluators of this package).
func TestGaussRules_NodesAreRoots(t *testing.T) {
	pairs := []struct {
		name string
		f    gaussRuleFunc
		p    func(int, float64) (float64, error)
	}{
		{"Legendre", GaussLegendre, LegendreP},
		{"Hermite", GaussHermite, HermiteH},
		{"HermiteE", GaussHermiteE, HermiteHe},
		{"Laguerre", GaussLaguerre, LaguerreL},
		{"ChebyshevT", GaussChebyshevT, ChebyshevT},
		{"ChebyshevU", GaussChebyshevU, ChebyshevU},
	}
	const n = 12
	for _, pr := range pairs {
		x, _, _ := pr.f(n)
		for _, xi := range x {
			v, _ := pr.p(n, xi)
			// scale: |p_n| at a point one node-gap away is O(p_n'·gap); 1e-9 relative to max |p_(n−1)| bounds it.
			ref, _ := pr.p(n-1, xi)
			if math.Abs(v) > 1e-9*math.Max(1, math.Abs(ref))*math.Max(1, math.Abs(xi)) {
				t.Errorf("%s: p_%d(%v) = %v; want ≈ 0", pr.name, n, xi, v)
			}
		}
	}
}

func TestGaussRules_Errors(t *testing.T) {
	for name, f := range map[string]gaussRuleFunc{
		"GaussLegendre": GaussLegendre, "GaussHermite": GaussHermite, "GaussHermiteE": GaussHermiteE,
		"GaussLaguerre": GaussLaguerre, "GaussChebyshevT": GaussChebyshevT, "GaussChebyshevU": GaussChebyshevU,
	} {
		for _, n := range []int{0, -3} {
			if _, _, err := f(n); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("%s(%d): err = %v; want ErrOutOfRange", name, n, err)
			}
		}
	}
}

func BenchmarkGaussLegendre(b *testing.B) {
	for _, n := range []int{10, 100} {
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = GaussLegendre(n)
			}
		})
	}
}

func BenchmarkGaussHermite(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = GaussHermite(20)
	}
}

func BenchmarkGaussHermiteE(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = GaussHermiteE(20)
	}
}

func BenchmarkGaussLaguerre(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = GaussLaguerre(20)
	}
}

func BenchmarkGaussChebyshevT(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = GaussChebyshevT(100)
	}
}

func BenchmarkGaussChebyshevU(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = GaussChebyshevU(100)
	}
}
