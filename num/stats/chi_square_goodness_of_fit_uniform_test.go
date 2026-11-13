package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Hand-checked from e-Handbook §1.3.5.15 with Eᵢ = N/k: 10, 20, 30, 40 → E = 25, χ² = 20.
func TestChiSquareGoodnessOfFitUniform_example(t *testing.T) {
	got, err := ChiSquareGoodnessOfFitUniform([]int{10, 20, 30, 40})
	if err != nil || !near(got, 20, 1e-13) {
		t.Errorf("= (%v, %v), want 20", got, err)
	}
}

// parity: ledger C2 (clustering.go:37), the source's day-part form: hour counts folded into 4 day parts,
// χ² = Σ (O − N/4)²/(N/4). Source case (built from its algorithm; it has no test file): parts 12, 3, 3, 2
// → N = 20, E = 5, χ² = (49 + 4 + 4 + 9)/5 = 13.2 (the source would flag it against 7.815, an app rule not
// ported). Also the 7-weekday form (clustering.go:121): 1, 1, 1, 1, 1, 1, 8 → E = 2, χ² = (6·1 + 36)/2 = 21.
func TestChiSquareGoodnessOfFitUniform_parity(t *testing.T) {
	for _, c := range []struct {
		obs  []int
		want float64
	}{
		{[]int{12, 3, 3, 2}, 13.2},
		{[]int{1, 1, 1, 1, 1, 1, 8}, 21},
	} {
		got, err := ChiSquareGoodnessOfFitUniform(c.obs)
		if err != nil || !near(got, c.want, 1e-13) {
			t.Errorf("%v: = (%v, %v), want %v", c.obs, got, err, c.want)
		}
	}
}

func TestChiSquareGoodnessOfFitUniform_cases(t *testing.T) {
	tests := []struct {
		name string
		obs  []int
		want float64
	}{
		{"one bin", []int{9}, 0},
		{"uniform", []int{4, 4, 4}, 0},
		{"all in one bin of two", []int{0, 6}, 6},
		{"non-integer expected", []int{1, 0, 0}, 2},
	}
	for _, tc := range tests {
		got, err := ChiSquareGoodnessOfFitUniform(tc.obs)
		if err != nil || !near(got, tc.want, 1e-13) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Oracle: NUM-STATS-033 with expected N/k on 300 random tables (fixed seed). Tolerance 1e-12 relative.
func TestChiSquareGoodnessOfFitUniform_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		k := 1 + rng.IntN(40)
		obs := make([]int, k)
		total := 0
		for i := range obs {
			obs[i] = rng.IntN(500)
			total += obs[i]
		}
		if total == 0 {
			continue
		}
		exp := make([]float64, k)
		for i := range exp {
			exp[i] = float64(total) / float64(k)
		}
		want, _ := ChiSquareGoodnessOfFit(obs, exp)
		got, err := ChiSquareGoodnessOfFitUniform(obs)
		if err != nil || !near(got, want, 1e-12) {
			t.Fatalf("%v: got (%v, %v), want %v", obs, got, err, want)
		}
	}
}

func TestChiSquareGoodnessOfFitUniform_errors(t *testing.T) {
	if _, err := ChiSquareGoodnessOfFitUniform(nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := ChiSquareGoodnessOfFitUniform([]int{3, -1}); !errors.Is(err, ErrNegative) {
		t.Errorf("negative: %v", err)
	}
	if _, err := ChiSquareGoodnessOfFitUniform([]int{0, 0}); !errors.Is(err, ErrZeroTotal) {
		t.Errorf("zero total: %v", err)
	}
	if _, err := ChiSquareGoodnessOfFitUniform([]int{math.MaxInt, 1}); !errors.Is(err, ErrOverflow) {
		t.Errorf("total overflow: %v", err)
	}
	if got, err := ChiSquareGoodnessOfFitUniform([]int{math.MaxInt / 2, 0}); err != nil || math.IsInf(got, 0) {
		t.Errorf("large counts: (%v, %v), want finite", got, err)
	}
}

func TestChiSquareGoodnessOfFitUniform_noAlloc(t *testing.T) {
	obs := []int{10, 20, 30, 40}
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = ChiSquareGoodnessOfFitUniform(obs) }); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func BenchmarkChiSquareGoodnessOfFitUniform(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		obs := make([]int, size.n)
		for i := range obs {
			obs[i] = rng.IntN(100)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = ChiSquareGoodnessOfFitUniform(obs)
			}
		})
	}
}
