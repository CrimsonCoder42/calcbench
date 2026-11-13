package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// nonAdjacentSlow tries every subset (n ≤ 16); ring adds the adjacency of the last and first.
func nonAdjacentSlow(values []int64, ring bool) int64 {
	n := len(values)
	var best int64
	for mask := 0; mask < 1<<n; mask++ {
		if mask&(mask>>1) != 0 || (ring && n > 1 && mask&1 == 1 && mask>>(n-1)&1 == 1) {
			continue
		}
		var s int64
		for i := range n {
			if mask>>i&1 == 1 {
				s += values[i]
			}
		}
		best = max(best, s)
	}
	return best
}

// Roughgarden, Algorithms Illuminated, Part 3, §16.1: the path with vertex weights 1, 4, 5, 4 has a
// maximum-weight independent set of weight 8 (the two 4s).
func TestMaxNonAdjacentSum_roughgardenPath(t *testing.T) {
	if got, err := MaxNonAdjacentSum([]int64{1, 4, 5, 4}); err != nil || got != 8 {
		t.Errorf("MaxNonAdjacentSum([1 4 5 4]) = (%d, %v), want (8, nil)", got, err)
	}
}

func TestMaxNonAdjacentSum_cases(t *testing.T) {
	tests := []struct {
		name         string
		values       []int64
		path, circle int64
	}{
		{"empty", nil, 0, 0},
		{"one", []int64{5}, 5, 5},
		{"one negative", []int64{-5}, 0, 0},
		{"two", []int64{5, 7}, 7, 7},
		{"three: ends on a path, middle on a ring", []int64{4, 5, 4}, 8, 5},
		{"four: ends apart on a path, adjacent on a ring", []int64{3, 1, 1, 3}, 6, 4},
		{"all negative", []int64{-1, -2, -3}, 0, 0},
		{"zeros", []int64{0, 0, 0}, 0, 0},
		{"sorted", []int64{1, 2, 3, 4, 5}, 9, 8},
		{"reversed", []int64{5, 4, 3, 2, 1}, 9, 8},
		{"extreme", []int64{math.MaxInt64, math.MinInt64}, math.MaxInt64, math.MaxInt64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MaxNonAdjacentSum(tt.values); err != nil || got != tt.path {
				t.Errorf("MaxNonAdjacentSum(%v) = (%d, %v), want %d", tt.values, got, err, tt.path)
			}
			if got, err := MaxNonAdjacentSumCircular(tt.values); err != nil || got != tt.circle {
				t.Errorf("MaxNonAdjacentSumCircular(%v) = (%d, %v), want %d", tt.values, got, err, tt.circle)
			}
		})
	}
}

func TestMaxNonAdjacentSum_errors(t *testing.T) {
	values := []int64{math.MaxInt64, 0, 1, 0}
	if _, err := MaxNonAdjacentSum(values); !errors.Is(err, ErrOverflow) {
		t.Errorf("MaxNonAdjacentSum overflow: err = %v", err)
	}
	if _, err := MaxNonAdjacentSumCircular(values); !errors.Is(err, ErrOverflow) {
		t.Errorf("MaxNonAdjacentSumCircular overflow: err = %v", err)
	}
}

func TestMaxNonAdjacentSum_allocations(t *testing.T) {
	values := []int64{1, 4, 5, 4, 2}
	if n := testing.AllocsPerRun(100, func() { _, _ = MaxNonAdjacentSumCircular(values) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func TestMaxNonAdjacentSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int64, rng.IntN(15))
		for i := range values {
			values[i] = rng.Int64N(41) - 10
		}
		if got, err := MaxNonAdjacentSum(values); err != nil || got != nonAdjacentSlow(values, false) {
			t.Fatalf("trial %d: MaxNonAdjacentSum(%v) = (%d, %v), want %d", trial, values, got, err, nonAdjacentSlow(values, false))
		}
		if got, err := MaxNonAdjacentSumCircular(values); err != nil || got != nonAdjacentSlow(values, true) {
			t.Fatalf("trial %d: MaxNonAdjacentSumCircular(%v) = (%d, %v), want %d", trial, values, got, err, nonAdjacentSlow(values, true))
		}
	}
}

func benchInt64s(n int, limit int64) []int64 {
	rng := rand.New(rand.NewPCG(1, 2))
	out := make([]int64, n)
	for i := range out {
		out[i] = rng.Int64N(limit)
	}
	return out
}

func BenchmarkMaxNonAdjacentSum(b *testing.B) {
	for _, size := range linearSizes {
		values := benchInt64s(size.n, 1000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MaxNonAdjacentSum(values)
			}
		})
	}
}

func BenchmarkMaxNonAdjacentSumCircular(b *testing.B) {
	for _, size := range linearSizes {
		values := benchInt64s(size.n, 1000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MaxNonAdjacentSumCircular(values)
			}
		})
	}
}
