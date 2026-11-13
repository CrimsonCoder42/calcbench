package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked example, hand-checked with the extended Horner scheme of DLMF §1.11(i): p(x) = 2x³ − 6x² + 2x − 1,
// p′(x) = 6x² − 12x + 2; at x = 3, p = 5 and p′ = 54 − 36 + 2 = 20.
func TestPolynomialValueAndDerivative_workedExample(t *testing.T) {
	v, d, err := PolynomialValueAndDerivative([]float64{-1, 2, -6, 2}, 3)
	if err != nil || v != 5 || d != 20 {
		t.Errorf("PolynomialValueAndDerivative = (%v, %v, %v), want (5, 20, nil)", v, d, err)
	}
}

func TestPolynomialValueAndDerivative_cases(t *testing.T) {
	tests := []struct {
		name         string
		coeffs       []float64
		x            float64
		value, deriv float64
	}{
		{"constant has derivative 0", []float64{7}, 3, 7, 0},
		{"linear", []float64{3, 2}, 10, 23, 2},
		{"x = 0 gives coeffs[0] and coeffs[1]", []float64{4, -5, 6, 7}, 0, 4, -5},
		{"negative x", []float64{0, 0, 1}, -3, 9, -6},
		{"stationary point", []float64{1, -2, 1}, 1, 0, 0},
		{"trailing zeros allowed", []float64{1, 1, 0, 0}, 2, 3, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, d, err := PolynomialValueAndDerivative(tt.coeffs, tt.x)
			if err != nil || v != tt.value || d != tt.deriv {
				t.Errorf("PolynomialValueAndDerivative(%v, %v) = (%v, %v, %v), want (%v, %v, nil)",
					tt.coeffs, tt.x, v, d, err, tt.value, tt.deriv)
			}
		})
	}
}

func TestPolynomialValueAndDerivative_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name   string
		coeffs []float64
		x      float64
		want   error
	}{
		{"nil", nil, 1, ErrEmpty},
		{"NaN x", []float64{1, 2}, nan, ErrNotFinite},
		{"Inf x", []float64{1, 2}, inf, ErrNotFinite},
		{"NaN leading coefficient", []float64{1, nan}, 1, ErrNotFinite},
		{"Inf constant term", []float64{-inf, 1}, 1, ErrNotFinite},
		{"value overflows", []float64{1, 0, 1e308}, 1e200, ErrOverflow},
		// p(x) = 2¹⁰²³·x³ at x = 1: p = 2¹⁰²³ is finite, p′ = 3·2¹⁰²³ is not.
		{"derivative overflows, value finite", []float64{0, 0, 0, 0x1p1023}, 1, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, d, err := PolynomialValueAndDerivative(tt.coeffs, tt.x)
			if !errors.Is(err, tt.want) || v != 0 || d != 0 {
				t.Errorf("PolynomialValueAndDerivative(%v, %v) = (%v, %v, %v), want (0, 0, %v)", tt.coeffs, tt.x, v, d, err, tt.want)
			}
		})
	}
}

// The value is computed by exactly the same recurrence as PolynomialEvaluate, so the two agree bit for bit.
func TestPolynomialValueAndDerivative_matchesEvaluate(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		coeffs := randomCoeffs(rng, 1+rng.IntN(30))
		x := 4*rng.Float64() - 2
		v, _, err := PolynomialValueAndDerivative(coeffs, x)
		want, err2 := PolynomialEvaluate(coeffs, x)
		if err != nil || err2 != nil || v != want {
			t.Fatalf("trial %d: value %v (%v), PolynomialEvaluate %v (%v)", trial, v, err, want, err2)
		}
	}
}

// Oracle: the exact derivative Σ i·aᵢ·xⁱ⁻¹ (math/big). Tolerance: 4n·u·Σ i·|aᵢ|·|x|ⁱ⁻¹, a safe multiple of
// the first-order bound for the derivative recurrence (Higham 2002, §5.2), which runs n−1 multiply-adds
// on values each carrying the error of the value recurrence.
func TestPolynomialValueAndDerivative_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		coeffs := randomCoeffs(rng, 2+rng.IntN(20))
		n := len(coeffs) - 1
		x := 4*rng.Float64() - 2
		_, got, err := PolynomialValueAndDerivative(coeffs, x)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		absDeriv := make([]float64, n)
		for i := 1; i <= n; i++ {
			absDeriv[i-1] = float64(i) * math.Abs(coeffs[i]) // for the bound only
		}
		exact := newBig(0)
		bx := newBig(x)
		for i := n; i >= 1; i-- {
			exact.Mul(exact, bx)
			exact.Add(exact, new(big.Float).SetPrec(bigPrec).Mul(newBig(float64(i)), newBig(coeffs[i])))
		}
		bound := 4 * float64(n) * unitRoundoff * bigToFloat(bigHorner(absDeriv, math.Abs(x)))
		diff := bigToFloat(new(big.Float).Sub(newBig(got), exact))
		if math.Abs(diff) > bound {
			t.Fatalf("trial %d: |error| = %g > bound %g", trial, diff, bound)
		}
	}
}

func TestPolynomialValueAndDerivative_noAllocation(t *testing.T) {
	coeffs := []float64{-1, 2, -6, 2}
	if allocs := testing.AllocsPerRun(100, func() { benchSinkFloat, _, _ = PolynomialValueAndDerivative(coeffs, 3) }); allocs != 0 {
		t.Errorf("PolynomialValueAndDerivative allocates %v times per call, want 0", allocs)
	}
}

func BenchmarkPolynomialValueAndDerivative(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range polynomialBenchDegrees {
		coeffs := randomCoeffs(rng, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSinkFloat, _, _ = PolynomialValueAndDerivative(coeffs, 0.999)
			}
		})
	}
}
