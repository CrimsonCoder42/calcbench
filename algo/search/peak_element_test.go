package search

import (
	"math"
	"math/rand/v2"
	"testing"
)

// isPeak reports whether values[i] is not less than its neighbours (outside the slice counts as lower).
func isPeak(values []int, i int) bool {
	return (i == 0 || values[i] >= values[i-1]) && (i == len(values)-1 || values[i] >= values[i+1])
}

// Kleinberg & Tardos, Algorithm Design, ch. 5, Solved Exercise 1: the peak of a unimodal array. Our own
// unimodal data; the peak 9 is at index 3.
func TestPeakElement_unimodalExample(t *testing.T) {
	if got, ok := PeakElement([]int{1, 3, 7, 9, 6, 2}); !ok || got != 3 {
		t.Errorf("PeakElement(unimodal) = (%d, %v), want (3, true)", got, ok)
	}
}

// Our own data: a mountain-shaped list (peak 11 at index 3), and a list with local peaks at 1 (6) and
// 5 (10), where either is accepted.
func TestPeakElement_exercises(t *testing.T) {
	if got, _ := PeakElement([]int{2, 4, 9, 11, 7, 3}); got != 3 {
		t.Errorf("mountain peak = %d, want 3", got)
	}
	several := []int{3, 6, 2, 4, 8, 10, 5}
	if got, _ := PeakElement(several); got != 1 && got != 5 {
		t.Errorf("PeakElement(%v) = %d, want 1 or 5", several, got)
	}
}

func TestPeakElement_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"one element", []int{4}, 0},
		{"two, rising", []int{1, 2}, 1},
		{"two, falling", []int{2, 1}, 0},
		{"ascending", []int{1, 2, 3, 4, 5}, 4},
		{"descending", []int{5, 4, 3, 2, 1}, 0},
		{"negative", []int{-5, -1, -3}, 1},
		{"extremes", []int{math.MinInt, math.MaxInt, math.MinInt}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := PeakElement(tt.values); !ok || got != tt.want {
				t.Errorf("PeakElement(%v) = (%d, %v), want (%d, true)", tt.values, got, ok, tt.want)
			}
		})
	}
	if got, ok := PeakElement([]int{}); ok || got != 0 {
		t.Errorf("PeakElement(empty) = (%d, %v), want (0, false)", got, ok)
	}
	if got, _ := PeakElement([]int{3, 3, 3}); !isPeak([]int{3, 3, 3}, got) {
		t.Errorf("PeakElement(all equal) = %d, not a peak", got)
	}
}

func TestPeakElement_floats(t *testing.T) {
	if got, _ := PeakElement([]float64{math.Inf(-1), 0.5, math.Inf(1), 2}); got != 2 {
		t.Errorf("PeakElement(floats) = %d, want 2", got)
	}
}

func TestPeakElement_noAllocation(t *testing.T) {
	values := []int{1, 3, 7, 9, 6, 2}
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _ = PeakElement(values) }); allocs != 0 {
		t.Errorf("PeakElement allocates %v times per call, want 0", allocs)
	}
}

// The oracle checks the peak property directly on random inputs, with and without equal neighbours.
func TestPeakElement_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		n := 1 + rng.IntN(60)
		values := make([]int, n)
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(50))
		}
		got, ok := PeakElement(values)
		if !ok || got < 0 || got >= n || !isPeak(values, got) {
			t.Fatalf("trial %d: PeakElement(%v) = (%d, %v), not a peak", trial, values, got, ok)
		}
	}
}

// Unimodal input: rising to n/3, then falling.
func BenchmarkPeakElement(b *testing.B) {
	for _, size := range benchSizes {
		values := make([]int, size.n)
		for i := range values {
			values[i] = -abs(i - size.n/3)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = PeakElement(values)
			}
		})
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
