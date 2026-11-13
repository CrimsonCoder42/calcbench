package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked example, hand-checked on the slope-intercept form f(x) = mx + b of Abramson, College Algebra 2e,
// OpenStax, 2021, §4.1: m = 2, b = 3, x = 4 → 2·4 + 3 = 11.
func TestLinearEvaluate_workedExample(t *testing.T) {
	if got, err := LinearEvaluate(2, 3, 4); err != nil || got != 11 {
		t.Errorf("LinearEvaluate(2, 3, 4) = (%v, %v), want (11, nil)", got, err)
	}
}

func TestLinearEvaluate_cases(t *testing.T) {
	tests := []struct {
		name                      string
		slope, intercept, x, want float64
	}{
		{"x = 0 gives the intercept", 5, -7, 0, -7},
		{"zero slope is a constant", 0, 4, 1e300, 4},
		{"negative slope", -3, 1, 2, -5},
		{"negative x", 2, 0, -8, -16},
		// slope·x = 1 − 2⁻⁶⁰ exactly; two roundings would give 1 − 1 = 0, one rounding gives −2⁻⁶⁰.
		{"one rounding (IEEE 754 fusedMultiplyAdd)", 1 + 0x1p-30, -1, 1 - 0x1p-30, -0x1p-60},
		{"largest finite value", 1, 0, math.MaxFloat64, math.MaxFloat64},
		{"smallest subnormal", 0, math.SmallestNonzeroFloat64, 1, math.SmallestNonzeroFloat64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LinearEvaluate(tt.slope, tt.intercept, tt.x)
			if err != nil || got != tt.want {
				t.Errorf("LinearEvaluate(%v, %v, %v) = (%v, %v), want (%v, nil)", tt.slope, tt.intercept, tt.x, got, err, tt.want)
			}
		})
	}
}

func TestLinearEvaluate_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name                string
		slope, intercept, x float64
		want                error
	}{
		{"NaN slope", nan, 0, 1, ErrNotFinite},
		{"Inf intercept", 1, -inf, 1, ErrNotFinite},
		{"Inf x", 1, 0, inf, ErrNotFinite},
		{"overflow", 1e200, 0, 1e200, ErrOverflow},
		{"overflow in the sum", 1, math.MaxFloat64, math.MaxFloat64, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LinearEvaluate(tt.slope, tt.intercept, tt.x)
			if !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("LinearEvaluate(%v, %v, %v) = (%v, %v), want (0, %v)", tt.slope, tt.intercept, tt.x, got, err, tt.want)
			}
		})
	}
}

// Oracle: slope·x + intercept computed exactly (math/big) and rounded once to float64 must equal the result
// bit for bit, since IEEE 754 fusedMultiplyAdd rounds once (round to nearest even in both).
func TestLinearEvaluate_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		slope := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(40)-20)
		intercept := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(40)-20)
		x := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(40)-20)
		got, err := LinearEvaluate(slope, intercept, x)
		exact := new(big.Float).SetPrec(bigPrec).Mul(newBig(slope), newBig(x))
		exact.Add(exact, newBig(intercept))
		if want := bigToFloat(exact); err != nil || got != want {
			t.Fatalf("trial %d: LinearEvaluate(%v, %v, %v) = (%v, %v), want %v", trial, slope, intercept, x, got, err, want)
		}
	}
}

func TestLinearEvaluate_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSinkFloat, _ = LinearEvaluate(2, 3, 4) }); allocs != 0 {
		t.Errorf("LinearEvaluate allocates %v times per call, want 0", allocs)
	}
}

func BenchmarkLinearEvaluate(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		x := 4.0
		for b.Loop() {
			benchSinkFloat, _ = LinearEvaluate(2, 3, x)
		}
	})
}
