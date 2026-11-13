package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example, hand-checked with the coefficient-form product of Cormen et al., 3rd ed., §30.1,
// eqs. (30.1)–(30.2): (2x + 3)(x − 4) = 2x² − 8x + 3x − 12 = 2x² − 5x − 12.
func TestPolynomialMultiply_workedExample(t *testing.T) {
	got, err := PolynomialMultiply([]float64{3, 2}, []float64{-4, 1})
	if err != nil || !slices.Equal(got, []float64{-12, -5, 2}) {
		t.Errorf("PolynomialMultiply = (%v, %v), want ([-12 -5 2], nil)", got, err)
	}
}

func TestPolynomialMultiply_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []float64
		want []float64
	}{
		{"constants", []float64{3}, []float64{-2}, []float64{-6}},
		{"constant times polynomial", []float64{2}, []float64{1, 2, 3}, []float64{2, 4, 6}},
		{"zero polynomial", []float64{0}, []float64{1, 2}, []float64{0, 0}},
		{"difference of squares", []float64{1, 1}, []float64{-1, 1}, []float64{-1, 0, 1}},
		{"square of a binomial", []float64{1, 1}, []float64{1, 1}, []float64{1, 2, 1}},
		{"trailing zeros kept", []float64{1, 0}, []float64{1, 0}, []float64{1, 0, 0}},
		{"different lengths", []float64{1, 2, 3}, []float64{4, 5}, []float64{4, 13, 22, 15}},
		{"negative coefficients", []float64{-1, -1}, []float64{-1, 1}, []float64{1, 0, -1}},
		{"largest finite value", []float64{math.MaxFloat64}, []float64{1}, []float64{math.MaxFloat64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PolynomialMultiply(tt.a, tt.b)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("PolynomialMultiply(%v, %v) = (%v, %v), want (%v, nil)", tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

func TestPolynomialMultiply_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name string
		a, b []float64
		want error
	}{
		{"a empty", nil, []float64{1}, ErrEmpty},
		{"b empty", []float64{1}, []float64{}, ErrEmpty},
		{"NaN in a", []float64{1, nan}, []float64{1}, ErrNotFinite},
		{"Inf in b", []float64{1}, []float64{inf, 1}, ErrNotFinite},
		{"product overflows", []float64{1e200}, []float64{1, 1e200}, ErrOverflow},
		{"sum of products overflows", []float64{math.MaxFloat64, math.MaxFloat64}, []float64{1, 1}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PolynomialMultiply(tt.a, tt.b)
			if !errors.Is(err, tt.want) || got != nil {
				t.Errorf("PolynomialMultiply(%v, %v) = (%v, %v), want (nil, %v)", tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

func TestPolynomialMultiply_doesNotModifyInput(t *testing.T) {
	a, b := []float64{1, 2, 3}, []float64{4, 5}
	if _, err := PolynomialMultiply(a, b); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(a, []float64{1, 2, 3}) || !slices.Equal(b, []float64{4, 5}) {
		t.Errorf("inputs modified: a = %v, b = %v", a, b)
	}
}

// Oracle: the exact convolution (math/big). Tolerance: each c[k] is a dot product of at most
// m = min(len(a), len(b)) terms accumulated with one rounding per term, so
// |ĉ[k] − c[k]| ≤ γₘ·Σ|a[i]·b[k−i]| (Higham 2002, §3.1); 1% added because the bound is computed in float64.
func TestPolynomialMultiply_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		a := randomCoeffs(rng, 1+rng.IntN(25))
		b := randomCoeffs(rng, 1+rng.IntN(25))
		got, err := PolynomialMultiply(a, b)
		if err != nil || len(got) != len(a)+len(b)-1 {
			t.Fatalf("trial %d: len %d, err %v", trial, len(got), err)
		}
		g := gamma(min(len(a), len(b)))
		for k := range got {
			exact := newBig(0)
			absSum := 0.0
			for i := max(0, k-len(b)+1); i <= min(k, len(a)-1); i++ {
				exact.Add(exact, new(big.Float).SetPrec(bigPrec).Mul(newBig(a[i]), newBig(b[k-i])))
				absSum += math.Abs(a[i] * b[k-i])
			}
			diff := bigToFloat(new(big.Float).Sub(newBig(got[k]), exact))
			if bound := 1.01 * g * absSum; math.Abs(diff) > bound {
				t.Fatalf("trial %d, c[%d]: |error| = %g > bound %g", trial, k, diff, bound)
			}
		}
	}
}

func TestPolynomialMultiply_oneAllocation(t *testing.T) {
	a, b := []float64{1, 2, 3}, []float64{4, 5}
	if allocs := testing.AllocsPerRun(100, func() { _, _ = PolynomialMultiply(a, b) }); allocs != 1 {
		t.Errorf("PolynomialMultiply allocates %v times per call, want 1 (the result)", allocs)
	}
}

func BenchmarkPolynomialMultiply(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range []struct {
		name string
		n    int
	}{{"1e1", 10}, {"1e2", 100}, {"1e3", 1_000}} {
		p, q := randomCoeffs(rng, size.n+1), randomCoeffs(rng, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				c, _ := PolynomialMultiply(p, q)
				benchSinkFloat = c[0]
			}
		})
	}
}
