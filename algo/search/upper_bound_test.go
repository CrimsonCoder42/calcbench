package search

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1: the sorted 16 keys. 653 is the 11th key, so the first index
// greater than it is 11. 400 is absent, so its upper and lower bounds agree (5). Hand-derived.
func TestUpperBound_knuthExample(t *testing.T) {
	if got := UpperBound(knuthSorted, 653); got != 11 {
		t.Errorf("UpperBound(knuth, 653) = %d, want 11", got)
	}
	if got := UpperBound(knuthSorted, 400); got != 5 {
		t.Errorf("UpperBound(knuth, 400) = %d, want 5", got)
	}
}

func TestUpperBound_cases(t *testing.T) {
	tests := []struct {
		name   string
		sorted []int
		target int
		want   int
	}{
		{"empty", nil, 5, 0},
		{"one element, target below", []int{5}, 4, 0},
		{"one element, target equal", []int{5}, 5, 1},
		{"one element, target above", []int{5}, 6, 1},
		{"two elements, target between", []int{1, 3}, 2, 1},
		{"duplicates, past the run", []int{1, 2, 2, 2, 3}, 2, 4},
		{"all equal", []int{7, 7, 7, 7}, 7, 4},
		{"target below every element", knuthSorted, 0, 0},
		{"target above every element", knuthSorted, 1000, 16},
		{"last element", knuthSorted, 908, 16},
		{"negative values", []int{-9, -5, -5, 0, 4}, -5, 3},
		{"minimum of the type", []int{math.MinInt, 0, math.MaxInt}, math.MinInt, 1},
		{"maximum of the type", []int{math.MinInt, 0, math.MaxInt}, math.MaxInt, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpperBound(tt.sorted, tt.target); got != tt.want {
				t.Errorf("UpperBound(%v, %d) = %d, want %d", tt.sorted, tt.target, got, tt.want)
			}
		})
	}
}

func TestUpperBound_floats(t *testing.T) {
	tests := []struct {
		name   string
		target float64
		want   int
	}{
		{"NaN sorts first", math.NaN(), 1},
		{"-Inf", math.Inf(-1), 2},
		{"-0 equals +0", math.Copysign(0, -1), 5},
		{"+0 equals -0", 0, 5},
		{"+Inf", math.Inf(1), 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpperBound(floatsSorted, tt.target); got != tt.want {
				t.Errorf("UpperBound(floats, %v) = %d, want %d", tt.target, got, tt.want)
			}
		})
	}
}

func TestUpperBound_unsortedDoesNotPanic(t *testing.T) {
	unsorted := []int{5, 1, 4, 2, 3}
	for target := 0; target <= 6; target++ {
		if got := UpperBound(unsorted, target); got < 0 || got > len(unsorted) {
			t.Errorf("UpperBound(unsorted, %d) = %d, outside [0, %d]", target, got, len(unsorted))
		}
	}
}

func TestUpperBound_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSink = UpperBound(knuthSorted, 653) }); allocs != 0 {
		t.Errorf("UpperBound allocates %v times per call, want 0", allocs)
	}
}

// linearUpperBound is the slow obvious oracle: scan for the first element greater than target.
func linearUpperBound(sorted []int, target int) int {
	for i, v := range sorted {
		if v > target {
			return i
		}
	}
	return len(sorted)
}

func TestUpperBound_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		sorted := randomSorted(rng, rng.IntN(64), 20)
		for target := -1; target <= 21; target++ {
			if got, want := UpperBound(sorted, target), linearUpperBound(sorted, target); got != want {
				t.Fatalf("trial %d: UpperBound(%v, %d) = %d, want %d", trial, sorted, target, got, want)
			}
		}
	}
}

func BenchmarkUpperBound(b *testing.B) {
	for _, size := range benchSizes {
		sorted, targets := benchData(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				sink += UpperBound(sorted, targets[i&1023])
			}
			benchSink = sink
		})
	}
}
