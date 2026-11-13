package stack

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// nearestExample is the hand-derived example of doc/algorithms_core_2026-11-12/catalog/stack.md.
var nearestExample = []int{2, 7, 3, 5, 4, 6, 8}

// nearestOracle is the O(n²) scan: for each i, walk outward on one side until qualifies holds.
func nearestOracle[T cmp.Ordered](values []T, forward bool, qualifies func(candidate, v T) bool) []int {
	out := make([]int, len(values))
	for i := range values {
		out[i] = -1
		if forward {
			for j := i + 1; j < len(values); j++ {
				if qualifies(values[j], values[i]) {
					out[i] = j
					break
				}
			}
		} else {
			for j := i - 1; j >= 0; j-- {
				if qualifies(values[j], values[i]) {
					out[i] = j
					break
				}
			}
		}
	}
	return out
}

// nearestCase is one table case shared by the four nearest-element tests.
type nearestCase struct {
	name   string
	values []int
	want   []int
}

func runNearestCases(t *testing.T, name string, f func([]int, []int) []int, tests []nearestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f(nil, tt.values); !slices.Equal(got, tt.want) {
				t.Errorf("%s(%v) = %v, want %v", name, tt.values, got, tt.want)
			}
		})
	}
}

// Hand-derived (catalog): 2 7 3 5 4 6 8 → 1 6 3 5 5 6 −1.
func TestNextGreater_example(t *testing.T) {
	if got := NextGreater(nil, nearestExample); !slices.Equal(got, []int{1, 6, 3, 5, 5, 6, -1}) {
		t.Errorf("NextGreater = %v", got)
	}
}

func TestNextGreater_cases(t *testing.T) {
	runNearestCases(t, "NextGreater", NextGreater[int], []nearestCase{
		{"empty", nil, []int{}},
		{"one", []int{4}, []int{-1}},
		{"two rising", []int{1, 2}, []int{1, -1}},
		{"equal is not greater", []int{3, 3, 3}, []int{-1, -1, -1}},
		{"sorted", []int{1, 2, 3, 4}, []int{1, 2, 3, -1}},
		{"reversed", []int{4, 3, 2, 1}, []int{-1, -1, -1, -1}},
		{"negative and zero", []int{-2, 0, -1, 0}, []int{1, -1, 3, -1}},
		{"extremes", []int{math.MaxInt, math.MinInt, math.MaxInt}, []int{-1, 2, -1}},
	})
}

func TestNextGreater_floats(t *testing.T) {
	// cmp.Compare order: NaN is least, so anything else is greater than a NaN.
	got := NextGreater(nil, []float64{math.NaN(), math.Inf(-1), 0, math.Inf(1)})
	if !slices.Equal(got, []int{1, 2, 3, -1}) {
		t.Errorf("NextGreater(floats) = %v, want [1 2 3 -1]", got)
	}
}

func TestNextGreater_appendsToDst(t *testing.T) {
	buf := make([]int, 2, 16)
	buf[0], buf[1] = 99, 98
	got := NextGreater(buf, []int{1, 0})
	if !slices.Equal(got, []int{99, 98, -1, -1}) || &got[0] != &buf[0] {
		t.Errorf("NextGreater(dst) = %v, want [99 98 -1 -1] in the same buffer", got)
	}
}

func TestNextGreater_noAllocWithBuffer(t *testing.T) {
	values := []int{5, 1, 4, 2, 3}
	buf := make([]int, 0, len(values))
	if allocs := testing.AllocsPerRun(100, func() { _ = NextGreater(buf[:0], values) }); allocs != 0 {
		t.Errorf("NextGreater with capacity allocated %v times, want 0", allocs)
	}
}

func TestNextGreater_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInts(rng, rng.IntN(120), 1+rng.IntN(30))
		want := nearestOracle(values, true, func(c, v int) bool { return c > v })
		if got := NextGreater(nil, values); !slices.Equal(got, want) {
			t.Fatalf("trial %d: NextGreater(%v) = %v, want %v", trial, values, got, want)
		}
	}
}

// randomInts returns n values drawn from [0, maxValue). Shared by the tests of this package.
func randomInts(rng *rand.Rand, n, maxValue int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = rng.IntN(maxValue)
	}
	return s
}

// benchInputs returns random values and a descending run (the deepest chain walks). Shared by benchmarks.
func benchInputs(n int) map[string][]int {
	random := randomInts(rand.New(rand.NewPCG(1, 2)), n, n)
	desc := make([]int, n)
	for i := range desc {
		desc[i] = n - i
	}
	return map[string][]int{"random": random, "descending": desc}
}

func benchmarkNearest(b *testing.B, f func([]int, []int) []int) {
	for _, order := range []string{"random", "descending"} {
		for _, size := range benchSizes {
			values := benchInputs(size.n)[order]
			buf := make([]int, 0, size.n)
			b.Run(order+"-"+size.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					buf = f(buf[:0], values)
				}
			})
		}
	}
}

func BenchmarkNextGreater(b *testing.B) { benchmarkNearest(b, NextGreater[int]) }
