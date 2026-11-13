package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: StRD NumAcc1 → (n − 1)·s² = 2·1² = 2 from the certified s = 1. The e-Handbook §7.2.6.2
// resistivity data → 0.02164115 (hand-checked, 7 significant digits).
func TestSumOfSquares_worked(t *testing.T) {
	if got, err := SumOfSquares(numAcc1); err != nil || got != 2 {
		t.Errorf("SumOfSquares(NumAcc1) = %v, %v; want 2", got, err)
	}
	// Tolerance 1e-6: the hand-checked value has 7 significant digits.
	if got, err := SumOfSquares(resistivity); err != nil || !near(got, 0.02164115, 1e-6) {
		t.Errorf("SumOfSquares(resistivity) = %v, %v; want 0.02164115", got, err)
	}
}

func TestSumOfSquares_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
		tol    float64
	}{
		{"one element", []float64{7}, 0, 0},
		{"two elements", []float64{1, 3}, 2, 0},
		{"all equal", []float64{5, 5, 5, 5}, 0, 0},
		{"negative and zero", []float64{-1, 0, 1}, 2, 0},
		{"large offset", []float64{1e9 + 1, 1e9 + 2, 1e9 + 3}, 2, 0},
		{"tiny spread", []float64{1, 1 + 0x1p-52}, 0x1p-105, 1e-15},
		{"large but finite", []float64{-1e153, 1e153}, 2e306, 1e-15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Exact where tol is 0 (small integers); otherwise one or two roundings.
			got, err := SumOfSquares(tt.values)
			if err != nil || (tt.tol == 0 && got != tt.want) || !near(got, tt.want, tt.tol) {
				t.Errorf("SumOfSquares(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSumOfSquares_errors(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"NaN", []float64{math.NaN()}, ErrNotFinite},
		{"Inf", []float64{1, math.Inf(1)}, ErrNotFinite},
		{"overflow", []float64{-m, m}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SumOfSquares(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("SumOfSquares(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSumOfSquares_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(400)
		offset := math.Pow(10, float64(rng.IntN(12)))
		values := make([]float64, n)
		for i := range values {
			values[i] = offset + rng.NormFloat64()
		}
		want := bigSumOfSquares(values)
		// Tolerance 1e-10 relative: each deviation x − x̄ is computed exactly (Sterbenz) from the rounded
		// mean, and the correction term −(Σd)²/n removes that rounding, leaving O(n·u) error.
		const tol = 1e-10
		if got, err := SumOfSquares(values); err != nil || !near(got, want, tol) {
			t.Fatalf("trial %d: SumOfSquares(n=%d, offset=%g) = %v, %v; oracle %v", trial, n, offset, got, err, want)
		}
	}
}

func TestSumOfSquares_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = SumOfSquares(michelso) }); a != 0 {
		t.Errorf("SumOfSquares allocates %v times per call, want 0", a)
	}
}

func BenchmarkSumOfSquares(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = SumOfSquares(inputs[i&7])
			}
		})
	}
}

// Tiny spread. 0 and 1e-160: Σ(x − x̄)² = (1e-160)²/2 ≈ 5e-321, a subnormal; the result is the exact
// rational value rounded to float64, to within one subnormal step 2^-1074 (the scaled sum is rounded once
// more on the way down). 2^-1000·(1, 2, 3): Σd² = 2·2^-2000, below the float64 range, rounds to 0, as
// documented, with no error. All-equal tiny values: exactly 0.
func TestSumOfSquares_tinySpread(t *testing.T) {
	values := []float64{0, 1e-160}
	_, _, sxx, _, _ := exactCentred(values, nil)
	want := ratFloat(sxx)
	if got, err := SumOfSquares(values); err != nil || math.Abs(got-want) > 0x1p-1074 || got == 0 {
		t.Errorf("SumOfSquares(0, 1e-160) = %v, %v; want %v", got, err, want)
	}
	if got, err := SumOfSquares([]float64{0x1p-1000, 0x2p-1000, 0x3p-1000}); err != nil || got != 0 {
		t.Errorf("SumOfSquares(2^-1000·(1 2 3)) = %v, %v; want 0 (below the float64 range)", got, err)
	}
	if got, err := SumOfSquares([]float64{1e-300, 1e-300, 1e-300}); err != nil || got != 0 {
		t.Errorf("SumOfSquares(1e-300 ×3) = %v, %v; want 0", got, err)
	}
}
