package stack

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived (catalog): 2 7 3 5 4 6 8 → −1 −1 1 1 3 1 −1.
func TestPreviousGreater_example(t *testing.T) {
	if got := PreviousGreater(nil, nearestExample); !slices.Equal(got, []int{-1, -1, 1, 1, 3, 1, -1}) {
		t.Errorf("PreviousGreater = %v", got)
	}
}

func TestPreviousGreater_cases(t *testing.T) {
	runNearestCases(t, "PreviousGreater", PreviousGreater[int], []nearestCase{
		{"empty", nil, []int{}},
		{"one", []int{4}, []int{-1}},
		{"two falling", []int{2, 1}, []int{-1, 0}},
		{"equal is not greater", []int{3, 3, 3}, []int{-1, -1, -1}},
		{"sorted", []int{1, 2, 3, 4}, []int{-1, -1, -1, -1}},
		{"reversed", []int{4, 3, 2, 1}, []int{-1, 0, 1, 2}},
		{"negative and zero", []int{0, -2, -1, -3}, []int{-1, 0, 0, 2}},
		{"extremes", []int{math.MaxInt, math.MinInt, math.MaxInt}, []int{-1, 0, -1}},
	})
}

// cmp.Compare order: NaN is the least value.
func TestPreviousGreater_floats(t *testing.T) {
	if got := PreviousGreater(nil, []float64{math.Inf(1), 0, math.Inf(-1), math.NaN()}); !slices.Equal(got, []int{-1, 0, 1, 2}) {
		t.Errorf("PreviousGreater(floats) = %v, want %v", got, []int{-1, 0, 1, 2})
	}
}

func TestPreviousGreater_noAllocWithBuffer(t *testing.T) {
	values := []int{5, 1, 4, 2, 3}
	buf := make([]int, 0, len(values))
	if allocs := testing.AllocsPerRun(100, func() { _ = PreviousGreater(buf[:0], values) }); allocs != 0 {
		t.Errorf("PreviousGreater with capacity allocated %v times, want 0", allocs)
	}
}

func TestPreviousGreater_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInts(rng, rng.IntN(120), 1+rng.IntN(30))
		want := nearestOracle(values, false, func(c, v int) bool { return c > v })
		if got := PreviousGreater(nil, values); !slices.Equal(got, want) {
			t.Fatalf("trial %d: PreviousGreater(%v) = %v, want %v", trial, values, got, want)
		}
	}
}

func BenchmarkPreviousGreater(b *testing.B) { benchmarkNearest(b, PreviousGreater[int]) }
