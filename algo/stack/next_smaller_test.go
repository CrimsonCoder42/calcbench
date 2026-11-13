package stack

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived (catalog): 2 7 3 5 4 6 8 → −1 2 −1 4 −1 −1 −1.
func TestNextSmaller_example(t *testing.T) {
	if got := NextSmaller(nil, nearestExample); !slices.Equal(got, []int{-1, 2, -1, 4, -1, -1, -1}) {
		t.Errorf("NextSmaller = %v", got)
	}
}

func TestNextSmaller_cases(t *testing.T) {
	runNearestCases(t, "NextSmaller", NextSmaller[int], []nearestCase{
		{"empty", nil, []int{}},
		{"one", []int{4}, []int{-1}},
		{"two falling", []int{2, 1}, []int{1, -1}},
		{"equal is not smaller", []int{3, 3, 3}, []int{-1, -1, -1}},
		{"sorted", []int{1, 2, 3, 4}, []int{-1, -1, -1, -1}},
		{"reversed", []int{4, 3, 2, 1}, []int{1, 2, 3, -1}},
		{"negative and zero", []int{0, -2, -1, -3}, []int{1, 3, 3, -1}},
		{"extremes", []int{math.MinInt, math.MaxInt, math.MinInt}, []int{-1, 2, -1}},
	})
}

// cmp.Compare order: NaN is the least value.
func TestNextSmaller_floats(t *testing.T) {
	if got := NextSmaller(nil, []float64{math.Inf(1), 0, math.Inf(-1), math.NaN()}); !slices.Equal(got, []int{1, 2, 3, -1}) {
		t.Errorf("NextSmaller(floats) = %v, want %v", got, []int{1, 2, 3, -1})
	}
}

func TestNextSmaller_noAllocWithBuffer(t *testing.T) {
	values := []int{5, 1, 4, 2, 3}
	buf := make([]int, 0, len(values))
	if allocs := testing.AllocsPerRun(100, func() { _ = NextSmaller(buf[:0], values) }); allocs != 0 {
		t.Errorf("NextSmaller with capacity allocated %v times, want 0", allocs)
	}
}

func TestNextSmaller_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInts(rng, rng.IntN(120), 1+rng.IntN(30))
		want := nearestOracle(values, true, func(c, v int) bool { return c < v })
		if got := NextSmaller(nil, values); !slices.Equal(got, want) {
			t.Fatalf("trial %d: NextSmaller(%v) = %v, want %v", trial, values, got, want)
		}
	}
}

func BenchmarkNextSmaller(b *testing.B) { benchmarkNearest(b, NextSmaller[int]) }
