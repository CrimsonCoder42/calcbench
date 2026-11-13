package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Hand-checked from the e-Handbook §1.3.5.15 formula: 10, 20, 30, 40 against 25 each gives
// (15² + 5² + 5² + 15²)/25 = 500/25 = 20. Tolerance 1e-13 relative.
func TestChiSquareGoodnessOfFit_example(t *testing.T) {
	got, err := ChiSquareGoodnessOfFit([]int{10, 20, 30, 40}, []float64{25, 25, 25, 25})
	if err != nil || !near(got, 20, 1e-13) {
		t.Errorf("ChiSquareGoodnessOfFit = (%v, %v), want 20", got, err)
	}
}

func TestChiSquareGoodnessOfFit_cases(t *testing.T) {
	tests := []struct {
		name     string
		observed []int
		expected []float64
		want     float64
	}{
		{"perfect fit", []int{5, 5}, []float64{5, 5}, 0},
		{"one bin", []int{7}, []float64{4}, 9.0 / 4},
		{"unequal expected", []int{8, 2}, []float64{6, 4}, 4.0/6 + 4.0/4},
		{"fractional expected", []int{1, 2, 3}, []float64{2.5, 2, 1.5}, 2.25/2.5 + 0 + 2.25/1.5},
		{"empty bin", []int{0, 10}, []float64{5, 5}, 10},
	}
	for _, tc := range tests {
		got, err := ChiSquareGoodnessOfFit(tc.observed, tc.expected)
		if err != nil || !near(got, tc.want, 1e-13) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Oracle: the direct Σ (O − E)²/E in float64 on 300 random tables (fixed seed, k ≤ 50, counts ≤ 1000).
// Tolerance 1e-12 relative (the direct form is not compensated).
func TestChiSquareGoodnessOfFit_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		k := 1 + rng.IntN(50)
		obs := make([]int, k)
		exp := make([]float64, k)
		want := 0.0
		for i := range k {
			obs[i] = rng.IntN(1001)
			exp[i] = 0.5 + 1000*rng.Float64()
			d := float64(obs[i]) - exp[i]
			want += d * d / exp[i]
		}
		got, err := ChiSquareGoodnessOfFit(obs, exp)
		if err != nil || !near(got, want, 1e-12) {
			t.Fatalf("k %d: got (%v, %v), want %v", k, got, err, want)
		}
	}
}

func TestChiSquareGoodnessOfFit_errors(t *testing.T) {
	if _, err := ChiSquareGoodnessOfFit(nil, nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := ChiSquareGoodnessOfFit([]int{1, 2}, []float64{1}); !errors.Is(err, ErrLengthMismatch) {
		t.Errorf("mismatch: %v", err)
	}
	if _, err := ChiSquareGoodnessOfFit([]int{-1, 2}, []float64{1, 1}); !errors.Is(err, ErrNegative) {
		t.Errorf("negative observed: %v", err)
	}
	for _, e := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := ChiSquareGoodnessOfFit([]int{1, 2}, []float64{1, e}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("expected %v: %v", e, err)
		}
	}
	for _, e := range []float64{0, -2} {
		if _, err := ChiSquareGoodnessOfFit([]int{1, 2}, []float64{1, e}); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("expected %v: %v", e, err)
		}
	}
	if _, err := ChiSquareGoodnessOfFit([]int{math.MaxInt, math.MaxInt}, []float64{1e-300, 1e-300}); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

func TestChiSquareGoodnessOfFit_noAlloc(t *testing.T) {
	obs, exp := []int{10, 20, 30, 40}, []float64{25, 25, 25, 25}
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = ChiSquareGoodnessOfFit(obs, exp) }); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func BenchmarkChiSquareGoodnessOfFit(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		obs := make([]int, size.n)
		exp := make([]float64, size.n)
		for i := range obs {
			obs[i] = rng.IntN(100)
			exp[i] = 50
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = ChiSquareGoodnessOfFit(obs, exp)
			}
		})
	}
}
