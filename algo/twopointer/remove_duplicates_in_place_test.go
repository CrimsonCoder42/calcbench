package twopointer

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// The "Duplicates" application of Sedgewick & Wayne, Algorithms 4th ed., §2.5: after sorting, equal
// values are adjacent and one pass keeps the first of each run. Hand-made data.
func TestRemoveDuplicatesInPlace_worked(t *testing.T) {
	values := []int{1, 1, 2, 3, 3, 3, 7}
	n := RemoveDuplicatesInPlace(values)
	if want := []int{1, 2, 3, 7}; n != 4 || !slices.Equal(values[:n], want) {
		t.Errorf("RemoveDuplicatesInPlace = %d, %v, want 4, %v", n, values[:n], want)
	}
}

func TestRemoveDuplicatesInPlace_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty", nil, []int{}},
		{"one", []int{5}, []int{5}},
		{"two equal", []int{5, 5}, []int{5}},
		{"two different", []int{5, 6}, []int{5, 6}},
		{"all equal", []int{2, 2, 2, 2}, []int{2}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"reversed: adjacent runs only", []int{3, 3, 2, 1, 1}, []int{3, 2, 1}},
		{"unsorted: non-adjacent repeats kept", []int{1, 2, 1}, []int{1, 2, 1}},
		{"negative and zero", []int{-1, -1, 0, 0}, []int{-1, 0}},
		{"extremes", []int{math.MinInt, math.MinInt, math.MaxInt}, []int{math.MinInt, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := slices.Clone(tt.values)
			n := RemoveDuplicatesInPlace(values)
			if n != len(tt.want) || !slices.Equal(values[:n], tt.want) {
				t.Errorf("RemoveDuplicatesInPlace(%v) = %d, %v, want %v", tt.values, n, values[:n], tt.want)
			}
		})
	}
}

func TestRemoveDuplicatesInPlace_floats(t *testing.T) {
	negZero := math.Copysign(0, -1)
	values := []float64{math.NaN(), math.NaN(), negZero, 0, 1}
	n := RemoveDuplicatesInPlace(values)
	// NaN != NaN, so both stay; -0 == +0, so +0 goes and -0 (the first) stays.
	if n != 4 || !math.IsNaN(values[0]) || !math.IsNaN(values[1]) || values[2] != 0 || !math.Signbit(values[2]) || values[3] != 1 {
		t.Errorf("RemoveDuplicatesInPlace(NaN NaN -0 0 1) = %d, %v", n, values[:n])
	}
}

func TestRemoveDuplicatesInPlace_noAllocation(t *testing.T) {
	src := make([]int, 1000)
	for i := range src {
		src[i] = i / 3
	}
	values := make([]int, len(src))
	allocs := testing.AllocsPerRun(20, func() {
		copy(values, src)
		benchSink = RemoveDuplicatesInPlace(values)
	})
	if allocs != 0 {
		t.Errorf("RemoveDuplicatesInPlace allocates %v times per call, want 0", allocs)
	}
}

// Oracle: slices.Compact on a copy keeps the same prefix.
func TestRemoveDuplicatesInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int, rng.IntN(30))
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(6))
		}
		if rng.IntN(2) == 0 {
			slices.Sort(values)
		}
		want := slices.Compact(slices.Clone(values))
		original := slices.Clone(values)
		n := RemoveDuplicatesInPlace(values)
		if !slices.Equal(values[:n], want) {
			t.Fatalf("trial %d: RemoveDuplicatesInPlace(%v) = %v, want %v", trial, original, values[:n], want)
		}
	}
}

func BenchmarkRemoveDuplicatesInPlace(b *testing.B) {
	for _, size := range benchSizes {
		src := benchSorted(size.n)
		values := make([]int64, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				copy(values, src)
				benchSink = RemoveDuplicatesInPlace(values)
			}
		})
	}
}
