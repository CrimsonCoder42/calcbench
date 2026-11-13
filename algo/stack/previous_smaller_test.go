package stack

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived (catalog): 2 7 3 5 4 6 8 → −1 0 0 2 2 4 5.
func TestPreviousSmaller_example(t *testing.T) {
	if got := PreviousSmaller(nil, nearestExample); !slices.Equal(got, []int{-1, 0, 0, 2, 2, 4, 5}) {
		t.Errorf("PreviousSmaller = %v", got)
	}
}

func TestPreviousSmaller_cases(t *testing.T) {
	runNearestCases(t, "PreviousSmaller", PreviousSmaller[int], []nearestCase{
		{"empty", nil, []int{}},
		{"one", []int{4}, []int{-1}},
		{"two rising", []int{1, 2}, []int{-1, 0}},
		{"equal is not smaller", []int{3, 3, 3}, []int{-1, -1, -1}},
		{"sorted", []int{1, 2, 3, 4}, []int{-1, 0, 1, 2}},
		{"reversed", []int{4, 3, 2, 1}, []int{-1, -1, -1, -1}},
		{"negative and zero", []int{-2, 0, -1, 0}, []int{-1, 0, 0, 2}},
		{"extremes", []int{math.MinInt, math.MaxInt, math.MinInt}, []int{-1, 0, -1}},
	})
}

// cmp.Compare order: NaN is the least value.
func TestPreviousSmaller_floats(t *testing.T) {
	if got := PreviousSmaller(nil, []float64{math.NaN(), math.Inf(-1), 0, math.Inf(1)}); !slices.Equal(got, []int{-1, 0, 1, 2}) {
		t.Errorf("PreviousSmaller(floats) = %v, want %v", got, []int{-1, 0, 1, 2})
	}
}

func TestPreviousSmaller_noAllocWithBuffer(t *testing.T) {
	values := []int{5, 1, 4, 2, 3}
	buf := make([]int, 0, len(values))
	if allocs := testing.AllocsPerRun(100, func() { _ = PreviousSmaller(buf[:0], values) }); allocs != 0 {
		t.Errorf("PreviousSmaller with capacity allocated %v times, want 0", allocs)
	}
}

func TestPreviousSmaller_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInts(rng, rng.IntN(120), 1+rng.IntN(30))
		want := nearestOracle(values, false, func(c, v int) bool { return c < v })
		if got := PreviousSmaller(nil, values); !slices.Equal(got, want) {
			t.Fatalf("trial %d: PreviousSmaller(%v) = %v, want %v", trial, values, got, want)
		}
	}
}

func BenchmarkPreviousSmaller(b *testing.B) { benchmarkNearest(b, PreviousSmaller[int]) }
