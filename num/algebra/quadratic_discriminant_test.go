package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// bigDiscriminant returns b² − 4ac exactly (math/big).
func bigDiscriminant(a, b, c float64) *big.Float {
	bb := new(big.Float).SetPrec(bigPrec).Mul(newBig(b), newBig(b))
	ac := new(big.Float).SetPrec(bigPrec).Mul(newBig(a), newBig(c))
	ac.Mul(ac, newBig(4))
	return bb.Sub(bb, ac)
}

// withinRelative reports whether |got − exact| ≤ tol·|exact|, evaluated exactly.
func withinRelative(got float64, exact *big.Float, tol float64) bool {
	diff := new(big.Float).SetPrec(bigPrec).Sub(newBig(got), exact)
	diff.Abs(diff)
	limit := new(big.Float).SetPrec(bigPrec).Abs(exact)
	limit.Mul(limit, newBig(tol))
	return diff.Cmp(limit) <= 0
}

// Worked example: Goldberg 1991, §1.4 "Cancellation", p. 10: b = 3.34, a = 1.22, c = 2.28. The exact
// discriminant of the decimals is 0.0292; three-digit arithmetic gives 0.1. Ours must be within 2u
// relative of the exact discriminant of the float64-rounded inputs (Jeannerod, Louvet, Muller 2013), and
// so within 1e-13 of 0.0292 (the rounding of the inputs moves D by about 11·2⁻⁵³ ≈ 1.2e-15).
func TestQuadraticDiscriminant_workedExample(t *testing.T) {
	a, b, c := 1.22, 3.34, 2.28
	got, err := QuadraticDiscriminant(a, b, c)
	if err != nil || !withinRelative(got, bigDiscriminant(a, b, c), 2*unitRoundoff) || math.Abs(got-0.0292) > 1e-13 {
		t.Errorf("QuadraticDiscriminant(1.22, 3.34, 2.28) = (%v, %v), want 0.0292 within 2 ulps of the exact value", got, err)
	}
}

func TestQuadraticDiscriminant_cases(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c float64
		want    float64
	}{
		{"positive", 1, -3, 2, 1},
		{"zero (double root)", 1, -2, 1, 0},
		{"negative", 1, 2, 5, -16},
		{"a = 0 gives b²", 0, 3, 7, 9},
		{"c = 0 gives b²", 5, -3, 0, 9},
		{"b = 0 gives −4ac", 2, 0, 3, -24},
		{"all zero", 0, 0, 0, 0},
		{"b² and 4ac overflow, D = 0", 0x1p699, 0x1p700, 0x1p699, 0},
		// b² = 2¹¹⁰⁰(1 + 2⁻⁵¹ + 2⁻¹⁰⁴), 4ac = 2¹¹⁰⁰(1 + 2⁻⁵¹): D = 2⁹⁹⁶ survives 104 bits of cancellation.
		{"b² and 4ac overflow, D finite", 0x1p549, 0x1p550 * (1 + 0x1p-52), 0x1p549 * (1 + 0x1p-51), 0x1p996},
		{"tiny inputs", 0x1p-500, 3 * 0x1p-500, 0x1p-500, 5 * 0x1p-1000},
		{"a and c of very different size", 0x1p900, 0, -0x1p-900, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QuadraticDiscriminant(tt.a, tt.b, tt.c)
			if err != nil || got != tt.want {
				t.Errorf("QuadraticDiscriminant(%v, %v, %v) = (%v, %v), want (%v, nil)", tt.a, tt.b, tt.c, got, err, tt.want)
			}
		})
	}
}

func TestQuadraticDiscriminant_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name    string
		a, b, c float64
		want    error
	}{
		{"NaN a", nan, 1, 1, ErrNotFinite},
		{"Inf b", 1, inf, 1, ErrNotFinite},
		{"−Inf c", 1, 1, -inf, ErrNotFinite},
		{"D overflows", 0x1p700, 0, -0x1p700, ErrOverflow},
		{"b² overflows with ac = 0", 0, 0x1p600, 1, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QuadraticDiscriminant(tt.a, tt.b, tt.c)
			if !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("QuadraticDiscriminant(%v, %v, %v) = (%v, %v), want (0, %v)", tt.a, tt.b, tt.c, got, err, tt.want)
			}
		})
	}
}

// Oracle: the exact discriminant (math/big), on inputs built so that b² ≈ 4ac (the hard case: the naive
// formula loses every digit) and on inputs of random size. Tolerance: relative error ≤ 2u, the bound of
// Jeannerod, Louvet, Muller 2013 for Kahan's algorithm with an FMA (no underflow or overflow occurs here).
func TestQuadraticDiscriminant_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 600 {
		scale := math.Ldexp(1, rng.IntN(1200)-600)
		a := (2*rng.Float64() - 1) * scale
		b := (2*rng.Float64() - 1) * scale
		var c float64
		if trial%2 == 0 {
			// Near a double root: c = b²/(4a)·(1 + δ), |δ| ≤ 2⁻⁴⁰.
			c = b * b / (4 * a) * (1 + (2*rng.Float64()-1)*0x1p-40)
		} else {
			c = (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(1200)-600)
		}
		if a == 0 || !reduce.AllFinite(c) {
			continue
		}
		exact := bigDiscriminant(a, b, c)
		if v := bigToFloat(exact); math.Abs(v) < 0x1p-1000 || math.Abs(v) > 0x1p1000 {
			continue // outside the normal range the 2u bound does not apply
		}
		got, err := QuadraticDiscriminant(a, b, c)
		if err != nil || !withinRelative(got, exact, 2*unitRoundoff) {
			t.Fatalf("trial %d: QuadraticDiscriminant(%v, %v, %v) = (%v, %v), exact %v", trial, a, b, c, got, err, bigToFloat(exact))
		}
	}
}

func TestQuadraticDiscriminant_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSinkFloat, _ = QuadraticDiscriminant(1.22, 3.34, 2.28) }); allocs != 0 {
		t.Errorf("QuadraticDiscriminant allocates %v times per call, want 0", allocs)
	}
}

func BenchmarkQuadraticDiscriminant(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		a := 1.22
		for b.Loop() {
			benchSinkFloat, _ = QuadraticDiscriminant(a, 3.34, 2.28)
		}
	})
}
