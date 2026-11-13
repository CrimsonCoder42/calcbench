package hashing

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived: 10 4 20 1 3 2 11 holds the run 1 2 3 4.
func TestLongestConsecutive_example(t *testing.T) {
	if start, n := LongestConsecutive([]int{10, 4, 20, 1, 3, 2, 11}); start != 1 || n != 4 {
		t.Errorf("LongestConsecutive = (%d, %d), want (1, 4)", start, n)
	}
}

func TestLongestConsecutive_cases(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		wantStart int
		wantLen   int
	}{
		{"empty", nil, 0, 0},
		{"one", []int{7}, 7, 1},
		{"two apart", []int{7, 9}, 7, 1},
		{"two adjacent", []int{8, 7}, 7, 2},
		{"duplicates count once", []int{2, 2, 3, 3, 1}, 1, 3},
		{"tie goes to smallest start", []int{20, 21, 5, 6, -3, -2}, -3, 2},
		{"negative through zero", []int{1, -1, 0, -2}, -2, 4},
		{"sorted", []int{1, 2, 3, 4, 5}, 1, 5},
		{"reversed", []int{5, 4, 3, 2, 1}, 1, 5},
		{"top of the type", []int{math.MaxInt, math.MaxInt - 1, math.MinInt}, math.MaxInt - 1, 2},
		{"bottom of the type", []int{math.MinInt + 1, math.MinInt, math.MaxInt}, math.MinInt, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if start, n := LongestConsecutive(tt.values); start != tt.wantStart || n != tt.wantLen {
				t.Errorf("LongestConsecutive(%v) = (%d, %d), want (%d, %d)", tt.values, start, n, tt.wantStart, tt.wantLen)
			}
		})
	}
}

func TestLongestConsecutive_unsignedEdges(t *testing.T) {
	if start, n := LongestConsecutive([]uint8{255, 0, 254, 1}); start != 0 || n != 2 {
		t.Errorf("LongestConsecutive(uint8 edges) = (%d, %d), want (0, 2): runs never wrap", start, n)
	}
	if start, n := LongestConsecutive([]int8{-128, 127, 126}); start != 126 || n != 2 {
		t.Errorf("LongestConsecutive(int8 edges) = (%d, %d), want (126, 2)", start, n)
	}
}

// Oracle: sort the distinct values and scan for the longest run (first run wins ties, which is the smallest start).
func TestLongestConsecutive_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInts(rng, rng.IntN(150), 1+rng.IntN(200))
		for i := range values {
			values[i] -= 50
		}
		s := slices.Compact(slices.Sorted(slices.Values(values)))
		wantStart, wantLen := 0, 0
		for i := 0; i < len(s); {
			j := i + 1
			for j < len(s) && s[j] == s[j-1]+1 {
				j++
			}
			if j-i > wantLen {
				wantStart, wantLen = s[i], j-i
			}
			i = j
		}
		if start, n := LongestConsecutive(values); start != wantStart || n != wantLen {
			t.Fatalf("trial %d: LongestConsecutive = (%d, %d), want (%d, %d)", trial, start, n, wantStart, wantLen)
		}
	}
}

func BenchmarkLongestConsecutive(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, 2*size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = LongestConsecutive(values)
			}
		})
	}
}
