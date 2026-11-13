package special

import (
	"errors"
	"math"
	"testing"
)

// Tolerance 1e-14: the continued fraction stops at ε per step; the prefactor adds a few ulps.
func TestRegularizedBeta_ClosedForms(t *testing.T) {
	for _, x := range []float64{0.01, 0.2, 0.5, 0.7, 0.99} {
		// I_x(1, 1) = x; I_x(a, 1) = x^a; I_x(1, b) = 1 − (1−x)^b (DLMF 8.17.1–8.17.2 integrated).
		if got, _ := RegularizedBeta(x, 1, 1); !within(got, x, 1e-14) {
			t.Errorf("I_%v(1,1) = %v; want %v", x, got, x)
		}
		if got, _ := RegularizedBeta(x, 3.5, 1); !within(got, math.Pow(x, 3.5), 1e-14) {
			t.Errorf("I_%v(3.5,1) = %v; want %v", x, got, math.Pow(x, 3.5))
		}
		if got, _ := RegularizedBeta(x, 1, 4); !within(got, 1-math.Pow(1-x, 4), 1e-14) {
			t.Errorf("I_%v(1,4) = %v; want %v", x, got, 1-math.Pow(1-x, 4))
		}
	}
}

func TestRegularizedBeta_BinomialSum(t *testing.T) {
	// DLMF 8.17.5: I_x(m, n−m+1) = Σ_(j=m..n) C(n,j) x^j (1−x)^(n−j). m = 2, n = 4, x = 0.3 → 0.3483.
	got, err := RegularizedBeta(0.3, 2, 3)
	if err != nil || !within(got, 0.3483, 1e-14) {
		t.Errorf("I_0.3(2,3) = %v, %v; want 0.3483", got, err)
	}
}

func TestRegularizedBeta_SymmetricMidpoint(t *testing.T) {
	// DLMF 8.17.4 with a = b and x = ½: I_½(a, a) = ½.
	// Tolerance 1e-13, and 1e-11 at a = 400: ln B(400, 400) cancels terms near 4.5e3 (about 5e-13 absolute).
	for _, a := range []float64{0.3, 2.5, 10, 400} {
		tol := 1e-13
		if a > 100 {
			tol = 1e-11
		}
		got, err := RegularizedBeta(0.5, a, a)
		if err != nil || math.Abs(got-0.5) > tol {
			t.Errorf("I_0.5(%v,%v) = %v, %v; want 0.5", a, a, got, err)
		}
	}
}

func TestRegularizedBeta_Edges(t *testing.T) {
	if got, err := RegularizedBeta(0, 2, 3); got != 0 || err != nil {
		t.Errorf("I_0 = %v, %v; want 0", got, err)
	}
	if got, err := RegularizedBeta(1, 2, 3); got != 1 || err != nil {
		t.Errorf("I_1 = %v, %v; want 1", got, err)
	}
}

func TestRegularizedBeta_Errors(t *testing.T) {
	cases := []struct {
		name    string
		x, a, b float64
		want    error
	}{
		{"x below 0", -0.1, 1, 1, ErrOutOfRange},
		{"x above 1", 1.1, 1, 1, ErrOutOfRange},
		{"a zero", 0.5, 0, 1, ErrOutOfRange},
		{"b negative", 0.5, 1, -1, ErrOutOfRange},
		{"x NaN", math.NaN(), 1, 1, ErrNotFinite},
		{"a Inf", 0.5, math.Inf(1), 1, ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := RegularizedBeta(c.x, c.a, c.b); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v; want %v", c.name, err, c.want)
		}
	}
}

func BenchmarkRegularizedBeta(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = RegularizedBeta(0.3, 2.5, 7.5)
	}
}
