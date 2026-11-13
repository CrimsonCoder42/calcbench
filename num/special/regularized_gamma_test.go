package special

import (
	"errors"
	"math"
	"testing"
)

// within reports |got − want| ≤ tol·max(1, |want|): relative for large values, absolute near zero.
func within(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(1, math.Abs(want))
}

var sinkF float64

// Tolerance 1e-14: both branches stop at ε = 2⁻⁵² per term; a few ulps of rounding accumulate.
func TestRegularizedGammaP_ExponentialCase(t *testing.T) {
	// DLMF 8.4.4: γ(1, x) = 1 − e^(−x), so P(1, x) = 1 − e^(−x) and Q(1, x) = e^(−x).
	for _, x := range []float64{0.1, 0.5, 1, 2, 5, 10, 30} {
		p, err := RegularizedGammaP(1, x)
		if err != nil || !within(p, -math.Expm1(-x), 1e-14) {
			t.Errorf("P(1, %v) = %v, %v; want %v", x, p, err, -math.Expm1(-x))
		}
		q, err := RegularizedGammaQ(1, x)
		if err != nil || math.Abs(q-math.Exp(-x)) > 1e-14*math.Exp(-x) {
			t.Errorf("Q(1, %v) = %v, %v; want %v", x, q, err, math.Exp(-x))
		}
	}
}

func TestRegularizedGammaP_HalfIsErf(t *testing.T) {
	// DLMF 8.4.1: γ(½, x) = √π erf(√x), so P(½, x) = erf(√x) and Q(½, x) = erfc(√x).
	for _, x := range []float64{0.01, 0.25, 1, 1.5, 4, 9, 16} {
		p, _ := RegularizedGammaP(0.5, x)
		if !within(p, math.Erf(math.Sqrt(x)), 1e-14) {
			t.Errorf("P(0.5, %v) = %v; want %v", x, p, math.Erf(math.Sqrt(x)))
		}
		q, _ := RegularizedGammaQ(0.5, x)
		want := math.Erfc(math.Sqrt(x))
		if math.Abs(q-want) > 1e-13*want {
			t.Errorf("Q(0.5, %v) = %v; want %v", x, q, want)
		}
	}
}

func TestRegularizedGammaQ_IntegerShapeIsPoissonSum(t *testing.T) {
	// DLMF 8.4.10: Q(n, x) = e^(−x) Σ_(k<n) x^k/k! for integer n (both branches: x < n+1 and x ≥ n+1).
	for _, n := range []int{2, 5, 12} {
		for _, x := range []float64{0.5, 3, 6, 13, 25} {
			sum, term := 0.0, 1.0
			for k := 0; k < n; k++ {
				if k > 0 {
					term *= x / float64(k)
				}
				sum += term
			}
			want := math.Exp(-x) * sum
			q, err := RegularizedGammaQ(float64(n), x)
			if err != nil || math.Abs(q-want) > 1e-13*math.Max(want, 1e-300)+1e-16 {
				t.Errorf("Q(%d, %v) = %v, %v; want %v", n, x, q, err, want)
			}
			p, _ := RegularizedGammaP(float64(n), x)
			if !within(p+q, 1, 1e-15) {
				t.Errorf("P+Q(%d, %v) = %v; want 1", n, x, p+q)
			}
		}
	}
}

func TestRegularizedGammaP_Edges(t *testing.T) {
	if p, err := RegularizedGammaP(3, 0); p != 0 || err != nil {
		t.Errorf("P(3, 0) = %v, %v; want 0", p, err)
	}
	if q, err := RegularizedGammaQ(3, 0); q != 1 || err != nil {
		t.Errorf("Q(3, 0) = %v, %v; want 1", q, err)
	}
	// Large shape, x = a: P is near ½ (median of the gamma distribution ≈ a − 1/3, DLMF 8.11 regime).
	p, err := RegularizedGammaP(1e4, 1e4)
	if err != nil || p < 0.5 || p > 0.51 {
		t.Errorf("P(1e4, 1e4) = %v, %v; want in (0.5, 0.51)", p, err)
	}
	// Far tails: P → 1, Q → 0 without NaN.
	q, err := RegularizedGammaQ(2, 800)
	if err != nil || q != 0 && q > 1e-300 {
		t.Errorf("Q(2, 800) = %v, %v; want underflow to ~0", q, err)
	}
}

func TestRegularizedGammaP_Errors(t *testing.T) {
	cases := []struct {
		name string
		a, x float64
		want error
	}{
		{"a zero", 0, 1, ErrOutOfRange},
		{"a negative", -1, 1, ErrOutOfRange},
		{"x negative", 1, -0.5, ErrOutOfRange},
		{"a NaN", math.NaN(), 1, ErrNotFinite},
		{"x NaN", 1, math.NaN(), ErrNotFinite},
		{"a Inf", math.Inf(1), 1, ErrNotFinite},
		{"x Inf", 1, math.Inf(1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := RegularizedGammaP(c.a, c.x); !errors.Is(err, c.want) {
			t.Errorf("P %s: err = %v; want %v", c.name, err, c.want)
		}
		if _, err := RegularizedGammaQ(c.a, c.x); !errors.Is(err, c.want) {
			t.Errorf("Q %s: err = %v; want %v", c.name, err, c.want)
		}
	}
}

func TestRegularizedGammaP_NoAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkF, _ = RegularizedGammaP(4.5, 7) }); n != 0 {
		t.Errorf("allocs = %v; want 0", n)
	}
}

func BenchmarkRegularizedGammaP(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = RegularizedGammaP(4.5, 3)
	}
}

func BenchmarkRegularizedGammaQ(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = RegularizedGammaQ(4.5, 9)
	}
}
