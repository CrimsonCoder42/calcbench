package search

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// knuthSorted is the 16-key running example of Knuth, TAOCP Vol. 3, 2nd ed. (ch. 5 and §6.2.1),
// in ascending order. Shared by every test file in this package.
var knuthSorted = []int{61, 87, 154, 170, 275, 426, 503, 509, 512, 612, 653, 677, 703, 765, 897, 908}

// floatsSorted is in slices.Sort order: NaN first, -0 and +0 equal. Shared by the float tests.
var floatsSorted = []float64{math.NaN(), math.Inf(-1), -1, math.Copysign(0, -1), 0, 1, math.Inf(1)}

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchSink keeps benchmark results observable.
var benchSink int

// randomSorted returns n values drawn from [0, maxValue), ascending, so runs of duplicates are common.
func randomSorted(rng *rand.Rand, n, maxValue int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = rng.IntN(maxValue)
	}
	slices.Sort(s)
	return s
}

// benchData returns the even numbers 0, 2, …, 2(n-1) and 1024 fixed-seed targets in [0, 2n), half of
// them present and half absent.
func benchData(n int) (sorted, targets []int) {
	sorted = make([]int, n)
	for i := range sorted {
		sorted[i] = 2 * i
	}
	rng := rand.New(rand.NewPCG(1, 2))
	targets = make([]int, 1024)
	for i := range targets {
		targets[i] = rng.IntN(2 * n)
	}
	return sorted, targets
}

// Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1: binary search of the 16 sorted keys for K = 653, which is the
// 11th key (index 10 here). 400 is absent; its insertion point, before 426, is hand-derived.
func TestLowerBound_knuthExample(t *testing.T) {
	if got := LowerBound(knuthSorted, 653); got != 10 {
		t.Errorf("LowerBound(knuth, 653) = %d, want 10", got)
	}
	if got := LowerBound(knuthSorted, 400); got != 5 {
		t.Errorf("LowerBound(knuth, 400) = %d, want 5", got)
	}
}

func TestLowerBound_cases(t *testing.T) {
	tests := []struct {
		name   string
		sorted []int
		target int
		want   int
	}{
		{"empty", nil, 5, 0},
		{"one element, target below", []int{5}, 4, 0},
		{"one element, target equal", []int{5}, 5, 0},
		{"one element, target above", []int{5}, 6, 1},
		{"two elements, target between", []int{1, 3}, 2, 1},
		{"duplicates, first of the run", []int{1, 2, 2, 2, 3}, 2, 1},
		{"all equal", []int{7, 7, 7, 7}, 7, 0},
		{"target below every element", knuthSorted, 0, 0},
		{"target above every element", knuthSorted, 1000, 16},
		{"first element", knuthSorted, 61, 0},
		{"last element", knuthSorted, 908, 15},
		{"negative values", []int{-9, -5, -5, 0, 4}, -5, 1},
		{"minimum of the type", []int{math.MinInt, 0, math.MaxInt}, math.MinInt, 0},
		{"maximum of the type", []int{math.MinInt, 0, math.MaxInt}, math.MaxInt, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LowerBound(tt.sorted, tt.target); got != tt.want {
				t.Errorf("LowerBound(%v, %d) = %d, want %d", tt.sorted, tt.target, got, tt.want)
			}
		})
	}
}

func TestLowerBound_floats(t *testing.T) {
	tests := []struct {
		name   string
		target float64
		want   int
	}{
		{"NaN sorts first", math.NaN(), 0},
		{"-Inf", math.Inf(-1), 1},
		{"-0 equals +0", math.Copysign(0, -1), 3},
		{"+0 equals -0", 0, 3},
		{"between 0 and 1", 0.5, 5},
		{"+Inf", math.Inf(1), 6},
		{"MaxFloat64 lands before +Inf", math.MaxFloat64, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LowerBound(floatsSorted, tt.target); got != tt.want {
				t.Errorf("LowerBound(floats, %v) = %d, want %d", tt.target, got, tt.want)
			}
		})
	}
}

func TestLowerBound_strings(t *testing.T) {
	sorted := []string{"apple", "banana", "banana", "cherry"}
	for target, want := range map[string]int{"": 0, "banana": 1, "blueberry": 3, "zucchini": 4} {
		if got := LowerBound(sorted, target); got != want {
			t.Errorf("LowerBound(%v, %q) = %d, want %d", sorted, target, got, want)
		}
	}
}

func TestLowerBound_unsortedDoesNotPanic(t *testing.T) {
	unsorted := []int{5, 1, 4, 2, 3}
	for target := 0; target <= 6; target++ {
		if got := LowerBound(unsorted, target); got < 0 || got > len(unsorted) {
			t.Errorf("LowerBound(unsorted, %d) = %d, outside [0, %d]", target, got, len(unsorted))
		}
	}
}

func TestLowerBound_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSink = LowerBound(knuthSorted, 653) }); allocs != 0 {
		t.Errorf("LowerBound allocates %v times per call, want 0", allocs)
	}
}

// linearLowerBound is the slow obvious oracle: scan for the first element not less than target.
func linearLowerBound(sorted []int, target int) int {
	for i, v := range sorted {
		if v >= target {
			return i
		}
	}
	return len(sorted)
}

func TestLowerBound_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		sorted := randomSorted(rng, rng.IntN(64), 20)
		before := slices.Clone(sorted)
		for target := -1; target <= 21; target++ {
			if got, want := LowerBound(sorted, target), linearLowerBound(sorted, target); got != want {
				t.Fatalf("trial %d: LowerBound(%v, %d) = %d, want %d", trial, sorted, target, got, want)
			}
		}
		if !slices.Equal(sorted, before) {
			t.Fatalf("trial %d: LowerBound modified its input", trial)
		}
	}
}

func BenchmarkLowerBound(b *testing.B) {
	for _, size := range benchSizes {
		sorted, targets := benchData(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				sink += LowerBound(sorted, targets[i&1023])
			}
			benchSink = sink
		})
	}
}
