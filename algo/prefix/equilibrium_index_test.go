package prefix

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveEquilibriumIndex is the oracle: both side sums recomputed in int for every index, O(n²).
func naiveEquilibriumIndex(values []int8) int {
	for i := range values {
		left, right := 0, 0
		for _, v := range values[:i] {
			left += int(v)
		}
		for _, v := range values[i+1:] {
			right += int(v)
		}
		if left == right {
			return i
		}
	}
	return -1
}

// Hand-derived from the prefix-sum identity (Laaksonen 2018, §9.1): in -7 1 5 2 -4 3 0, index 3 has
// -7+1+5 = -1 on the left and -4+3+0 = -1 on the right; index 6 also balances (0 = 0), but 3 is first.
func TestEquilibriumIndex_worked(t *testing.T) {
	if got := EquilibriumIndex([]int{-7, 1, 5, 2, -4, 3, 0}); got != 3 {
		t.Errorf("EquilibriumIndex(-7 1 5 2 -4 3 0) = %d, want 3", got)
	}
}

func TestEquilibriumIndex_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   int
	}{
		{"empty", nil, -1},
		{"one: both sides empty", []int64{5}, 0},
		{"two, unequal", []int64{1, 2}, -1},
		{"two zeros", []int64{0, 0}, 0},
		{"last index", []int64{1, -1, 7}, 2},
		{"duplicates", []int64{2, 9, 2}, 1},
		{"none", []int64{1, 2, 3}, -1},
		{"extremes, sides beyond int64", []int64{math.MaxInt64, math.MaxInt64, 0, math.MaxInt64, math.MaxInt64}, 2},
		{"minimums", []int64{math.MinInt64, math.MinInt64, 5, math.MinInt64, math.MinInt64}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EquilibriumIndex(tt.values); got != tt.want {
				t.Errorf("EquilibriumIndex(%v) = %d, want %d", tt.values, got, tt.want)
			}
		})
	}
	if got := EquilibriumIndex([]uint8{255, 3, 255}); got != 1 {
		t.Errorf("EquilibriumIndex(uint8 255 3 255) = %d, want 1", got)
	}
}

func TestEquilibriumIndex_noAllocation(t *testing.T) {
	values := benchInts(1000)
	allocs := testing.AllocsPerRun(20, func() { benchSink = EquilibriumIndex(values) })
	if allocs != 0 {
		t.Errorf("EquilibriumIndex allocates %v times per call, want 0", allocs)
	}
}

func TestEquilibriumIndex_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int8, rng.IntN(12))
		for i := range values {
			values[i] = int8(rng.IntN(7) - 3)
		}
		if rng.IntN(4) == 0 {
			values = randomInt8s(rng, len(values))
		}
		before := slices.Clone(values)
		if got, want := EquilibriumIndex(values), naiveEquilibriumIndex(values); got != want {
			t.Fatalf("trial %d: EquilibriumIndex(%v) = %d, want %d", trial, values, got, want)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkEquilibriumIndex(b *testing.B) {
	for _, size := range benchSizes {
		values := benchInts(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = EquilibriumIndex(values)
			}
		})
	}
}
