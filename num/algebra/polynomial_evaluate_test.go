package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Shared by the tests and benchmarks of this package.
var (
	benchSinkFloat   float64
	benchSinkComplex complex128
)

// bigPrec is enough bits for every exact computation the tests make: products of up to 20 float64 values
// over the whole exponent range.
const bigPrec = 1 << 13

func newBig(v float64) *big.Float { return new(big.Float).SetPrec(bigPrec).SetFloat64(v) }

// bigHorner returns Σ coeffs[i]·xⁱ exactly (for the degrees the tests use).
func bigHorner(coeffs []float64, x float64) *big.Float {
	s := newBig(0)
	bx := newBig(x)
	for i := len(coeffs) - 1; i >= 0; i-- {
		s.Mul(s, bx)
		s.Add(s, newBig(coeffs[i]))
	}
	return s
}

// bigToFloat rounds f to the nearest float64.
func bigToFloat(f *big.Float) float64 {
	v, _ := f.Float64()
	return v
}

// gamma returns γₖ = k·u/(1 − k·u), Higham 2002, §3.1.
func gamma(k int) float64 {
	ku := float64(k) * unitRoundoff
	return ku / (1 - ku)
}

// randomCoeffs returns n values uniform in [−1, 1).
func randomCoeffs(rng *rand.Rand, n int) []float64 {
	c := make([]float64, n)
	for i := range c {
		c[i] = 2*rng.Float64() - 1
	}
	return c
}

// Worked example, hand-checked by Horner's rule as written in Knuth, TAOCP Vol. 2, 3rd ed., §4.6.4 and
// DLMF §1.11(i): p(x) = 2x³ − 6x² + 2x − 1 at x = 3: 2 → 2·3 − 6 = 0 → 0·3 + 2 = 2 → 2·3 − 1 = 5.
func TestPolynomialEvaluate_workedExample(t *testing.T) {
	got, err := PolynomialEvaluate([]float64{-1, 2, -6, 2}, 3)
	if err != nil || got != 5 {
		t.Errorf("PolynomialEvaluate = (%v, %v), want (5, nil)", got, err)
	}
}

func TestPolynomialEvaluate_cases(t *testing.T) {
	tests := []struct {
		name   string
		coeffs []float64
		x      float64
		want   float64
	}{
		{"constant", []float64{7}, 123, 7},
		{"zero polynomial", []float64{0}, 5, 0},
		{"x = 0 gives coeffs[0]", []float64{-4, 9, 9}, 0, -4},
		{"x = 1 gives the sum", []float64{1, 2, 3, 4}, 1, 10},
		{"negative x", []float64{0, 0, 1}, -3, 9},
		{"trailing zeros allowed", []float64{1, 1, 0, 0}, 2, 3},
		{"linear", []float64{3, 2}, 4, 11},
		{"one rounding per step (fma)", []float64{-1, 1 + 0x1p-30}, 1 - 0x1p-30, -0x1p-60},
		{"largest finite value", []float64{math.MaxFloat64}, 2, math.MaxFloat64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PolynomialEvaluate(tt.coeffs, tt.x)
			if err != nil || got != tt.want {
				t.Errorf("PolynomialEvaluate(%v, %v) = (%v, %v), want (%v, nil)", tt.coeffs, tt.x, got, err, tt.want)
			}
		})
	}
}

func TestPolynomialEvaluate_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name   string
		coeffs []float64
		x      float64
		want   error
	}{
		{"nil", nil, 1, ErrEmpty},
		{"empty", []float64{}, 1, ErrEmpty},
		{"NaN x", []float64{1, 2}, nan, ErrNotFinite},
		{"Inf x", []float64{1, 2}, -inf, ErrNotFinite},
		{"NaN leading coefficient", []float64{1, nan}, 1, ErrNotFinite},
		{"Inf inner coefficient", []float64{1, inf, 1}, 1, ErrNotFinite},
		{"NaN constant term", []float64{nan, 1}, 1, ErrNotFinite},
		{"overflow", []float64{1, 0, 1e308}, 1e200, ErrOverflow},
		{"overflow part-way through", []float64{0, math.MaxFloat64, math.MaxFloat64}, 2, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PolynomialEvaluate(tt.coeffs, tt.x)
			if !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("PolynomialEvaluate(%v, %v) = (%v, %v), want (0, %v)", tt.coeffs, tt.x, got, err, tt.want)
			}
		})
	}
}

// Oracle: the exact value (math/big). Tolerance: Higham 2002 §5.1, |p̂ − p| ≤ γ₂ₙ·Σ|aᵢ||x|ⁱ; the bound is
// itself computed in float64, so it is enlarged by 1%.
func TestPolynomialEvaluate_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		coeffs := randomCoeffs(rng, 1+rng.IntN(20))
		x := 4*rng.Float64() - 2
		got, err := PolynomialEvaluate(coeffs, x)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		abs := make([]float64, len(coeffs))
		for i, c := range coeffs {
			abs[i] = math.Abs(c)
		}
		bound := 1.01 * gamma(2*(len(coeffs)-1)) * bigToFloat(bigHorner(abs, math.Abs(x)))
		diff := new(big.Float).Sub(newBig(got), bigHorner(coeffs, x))
		if math.Abs(bigToFloat(diff)) > bound {
			t.Fatalf("trial %d: |error| = %g > bound %g", trial, bigToFloat(diff), bound)
		}
	}
}

func TestPolynomialEvaluate_noAllocation(t *testing.T) {
	coeffs := []float64{-1, 2, -6, 2}
	if allocs := testing.AllocsPerRun(100, func() { benchSinkFloat, _ = PolynomialEvaluate(coeffs, 3) }); allocs != 0 {
		t.Errorf("PolynomialEvaluate allocates %v times per call, want 0", allocs)
	}
}

// polynomialBenchDegrees are the degrees used as sizes by the O(n) polynomial benchmarks.
var polynomialBenchDegrees = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

func BenchmarkPolynomialEvaluate(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range polynomialBenchDegrees {
		coeffs := randomCoeffs(rng, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSinkFloat, _ = PolynomialEvaluate(coeffs, 0.999)
			}
		})
	}
}
