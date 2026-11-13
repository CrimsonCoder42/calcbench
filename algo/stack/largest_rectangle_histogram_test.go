package stack

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// histogramOracle is the O(n²) scan: for every start, extend right keeping the running minimum.
func histogramOracle(heights []int64) int64 {
	var best int64
	for i := range heights {
		low := int64(math.MaxInt64)
		for j := i; j < len(heights); j++ {
			low = min(low, heights[j])
			best = max(best, low*int64(j-i+1))
		}
	}
	return best
}

// Hand-derived (catalog), our own heights: 3 1 4 4 3 5 1 → 12, bars 2–5 at height 3 (beats 4×2 = 8 and 1×7 = 7).
func TestLargestRectangleHistogram_example(t *testing.T) {
	if got, err := LargestRectangleHistogram([]int64{3, 1, 4, 4, 3, 5, 1}); err != nil || got != 12 {
		t.Errorf("LargestRectangleHistogram = (%d, %v), want (12, nil)", got, err)
	}
}

func TestLargestRectangleHistogram_cases(t *testing.T) {
	tests := []struct {
		name    string
		heights []int64
		want    int64
	}{
		{"empty", nil, 0},
		{"one", []int64{7}, 7},
		{"one zero", []int64{0}, 0},
		{"two", []int64{2, 4}, 4},
		{"all equal", []int64{3, 3, 3, 3}, 12},
		{"sorted", []int64{1, 2, 3, 4, 5}, 9},
		{"reversed", []int64{5, 4, 3, 2, 1}, 9},
		{"zero splits", []int64{4, 4, 0, 3, 3, 3}, 9},
		{"wide and low beats tall", []int64{1, 1, 1, 1, 1, 6}, 6},
		{"max height, width 1", []int64{math.MaxInt64}, math.MaxInt64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := LargestRectangleHistogram(tt.heights); err != nil || got != tt.want {
				t.Errorf("LargestRectangleHistogram(%v) = (%d, %v), want (%d, nil)", tt.heights, got, err, tt.want)
			}
		})
	}
}

func TestLargestRectangleHistogram_errors(t *testing.T) {
	tests := []struct {
		name    string
		heights []int64
		want    error
	}{
		{"negative", []int64{3, -1}, ErrNegativeHeight},
		{"min int64", []int64{math.MinInt64}, ErrNegativeHeight},
		{"area overflow", []int64{math.MaxInt64, math.MaxInt64}, ErrOverflow},
		{"area overflow, half max wide", []int64{math.MaxInt64/2 + 1, math.MaxInt64/2 + 1}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := LargestRectangleHistogram(tt.heights); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("LargestRectangleHistogram(%v) = (%d, %v), want (0, %v)", tt.heights, got, err, tt.want)
			}
		})
	}
}

func TestLargestRectangleHistogram_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		heights := make([]int64, rng.IntN(100))
		for i := range heights {
			heights[i] = int64(rng.IntN(1 + rng.IntN(20)))
		}
		want := histogramOracle(heights)
		if got, err := LargestRectangleHistogram(heights); err != nil || got != want {
			t.Fatalf("trial %d: LargestRectangleHistogram(%v) = (%d, %v), want %d", trial, heights, got, err, want)
		}
	}
}

func BenchmarkLargestRectangleHistogram(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		heights := make([]int64, size.n)
		for i := range heights {
			heights[i] = int64(rng.IntN(size.n))
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = LargestRectangleHistogram(heights)
			}
		})
	}
}
