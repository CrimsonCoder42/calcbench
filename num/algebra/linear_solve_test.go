package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked example, hand-checked on the slope-intercept form of Abramson, College Algebra 2e, OpenStax, 2021,
// §4.1: 2x + 3 = 11 gives x = (11 − 3)/2 = 4, the inverse of LinearEvaluate's example.
func TestLinearSolve_workedExample(t *testing.T) {
	if got, err := LinearSolve(2, 3, 11); err != nil || got != 4 {
		t.Errorf("LinearSolve(2, 3, 11) = (%v, %v), want (4, nil)", got, err)
	}
}

func TestLinearSolve_cases(t *testing.T) {
	tests := []struct {
		name                      string
		slope, intercept, y, want float64
	}{
		{"y = intercept gives x = 0", 5, -7, -7, 0},
		{"negative slope", -2, 1, 7, -3},
		{"zero intercept", 4, 0, 2, 0.5},
		{"tiny slope, finite result", 0x1p-1000, 0, 0x1p-30, 0x1p970},
		{"largest finite value", 1, 0, math.MaxFloat64, math.MaxFloat64},
		{"result underflows to a subnormal", 0x1p60, 0, 0x1p-1030, 0x1p-1090},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LinearSolve(tt.slope, tt.intercept, tt.y)
			if err != nil || got != tt.want {
				t.Errorf("LinearSolve(%v, %v, %v) = (%v, %v), want (%v, nil)", tt.slope, tt.intercept, tt.y, got, err, tt.want)
			}
		})
	}
}

func TestLinearSolve_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name                string
		slope, intercept, y float64
		want                error
	}{
		{"zero slope", 0, 1, 2, ErrDegenerate},
		{"negative zero slope", math.Copysign(0, -1), 1, 1, ErrDegenerate},
		{"NaN slope", nan, 0, 1, ErrNotFinite},
		{"Inf intercept", 1, inf, 1, ErrNotFinite},
		{"NaN y", 1, 0, nan, ErrNotFinite},
		{"y − intercept overflows", 1, -math.MaxFloat64, math.MaxFloat64, ErrOverflow},
		{"quotient overflows", 1e-300, 0, 1e300, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LinearSolve(tt.slope, tt.intercept, tt.y)
			if !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("LinearSolve(%v, %v, %v) = (%v, %v), want (0, %v)", tt.slope, tt.intercept, tt.y, got, err, tt.want)
			}
		})
	}
}

// Oracle: the exact (y − intercept)/slope (math/big). Tolerance: one rounding in the subtraction and one in
// the division put the result within (1 + u)² − 1 < 2.0001u of the exact quotient, and want (the exact
// quotient rounded) is within u/2 of it, so |got − want| ≤ 2.5001u·|want|.
func TestLinearSolve_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		slope := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(40)-20)
		if slope == 0 {
			continue
		}
		intercept := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(40)-20)
		y := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(40)-20)
		got, err := LinearSolve(slope, intercept, y)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		exact := new(big.Float).SetPrec(bigPrec).Sub(newBig(y), newBig(intercept))
		exact.Quo(exact, newBig(slope))
		want := bigToFloat(exact)
		if math.Abs(got-want) > 2.5001*unitRoundoff*math.Abs(want) {
			t.Fatalf("trial %d: LinearSolve(%v, %v, %v) = %v, want %v", trial, slope, intercept, y, got, want)
		}
	}
}

func TestLinearSolve_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSinkFloat, _ = LinearSolve(2, 3, 11) }); allocs != 0 {
		t.Errorf("LinearSolve allocates %v times per call, want 0", allocs)
	}
}

func BenchmarkLinearSolve(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		y := 11.0
		for b.Loop() {
			benchSinkFloat, _ = LinearSolve(2, 3, y)
		}
	})
}
