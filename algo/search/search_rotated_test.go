package search

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// rotate returns s rotated left by r: s[r:] followed by s[:r].
func rotate(s []int, r int) []int {
	return append(slices.Clone(s[r:]), s[:r]...)
}

// Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1 searches the 16 sorted keys for 653. Rotated by every amount,
// 653 is found where the rotation put it; 400 is never found.
func TestSearchRotated_knuthExample(t *testing.T) {
	for r := range len(knuthSorted) {
		rotated := rotate(knuthSorted, r)
		want := (10 - r + 16) % 16
		if got, found := SearchRotated(rotated, 653); !found || got != want {
			t.Errorf("rotation %d: SearchRotated(653) = (%d, %v), want (%d, true)", r, got, found, want)
		}
		if got, found := SearchRotated(rotated, 400); found || got != -1 {
			t.Errorf("rotation %d: SearchRotated(400) = (%d, %v), want (-1, false)", r, got, found)
		}
	}
}

// Our own wording of the common exercise: a sorted list of distinct numbers turned at an unknown point.
func TestSearchRotated_exercises(t *testing.T) {
	rotated := []int{40, 50, 60, 10, 20, 30}
	if got, found := SearchRotated(rotated, 20); !found || got != 4 {
		t.Errorf("SearchRotated(20) = (%d, %v), want (4, true)", got, found)
	}
	if got, found := SearchRotated(rotated, 35); found || got != -1 {
		t.Errorf("SearchRotated(35) = (%d, %v), want (-1, false)", got, found)
	}
}

func TestSearchRotated_cases(t *testing.T) {
	tests := []struct {
		name    string
		rotated []int
		target  int
		want    int
	}{
		{"empty", nil, 1, -1},
		{"one element, present", []int{5}, 5, 0},
		{"one element, absent", []int{5}, 4, -1},
		{"two, rotated", []int{9, 3}, 3, 1},
		{"two, not rotated", []int{3, 9}, 9, 1},
		{"not rotated", []int{1, 2, 3, 4}, 1, 0},
		{"first element", []int{4, 5, 1, 2, 3}, 4, 0},
		{"last element", []int{4, 5, 1, 2, 3}, 3, 4},
		{"rotation point", []int{4, 5, 1, 2, 3}, 1, 2},
		{"below every element", []int{4, 5, 1, 2, 3}, 0, -1},
		{"above every element", []int{4, 5, 1, 2, 3}, 9, -1},
		{"negative", []int{-1, 0, -5, -3}, -3, 3},
		{"extremes", []int{math.MaxInt, math.MinInt, 0}, math.MinInt, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := SearchRotated(tt.rotated, tt.target)
			if got != tt.want || found != (tt.want >= 0) {
				t.Errorf("SearchRotated(%v, %d) = (%d, %v), want %d", tt.rotated, tt.target, got, found, tt.want)
			}
		})
	}
}

func TestSearchRotated_duplicatesDoNotPanic(t *testing.T) {
	rotated := []int{2, 2, 3, 1, 2, 2}
	for target := 0; target <= 4; target++ {
		if got, _ := SearchRotated(rotated, target); got < -1 || got >= len(rotated) {
			t.Errorf("SearchRotated(dups, %d) = %d, outside [-1, %d)", target, got, len(rotated))
		}
	}
}

func TestSearchRotated_noAllocation(t *testing.T) {
	rotated := rotate(knuthSorted, 5)
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _ = SearchRotated(rotated, 653) }); allocs != 0 {
		t.Errorf("SearchRotated allocates %v times per call, want 0", allocs)
	}
}

func TestSearchRotated_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		n := rng.IntN(40)
		sorted := make([]int, n) // distinct: strictly ascending.
		for i := range sorted {
			sorted[i] = 3*i + rng.IntN(2)
		}
		r := 0
		if n > 0 {
			r = rng.IntN(n)
		}
		rotated := rotate(sorted, r)
		for target := -1; target <= 3*n+1; target++ {
			want := slices.Index(rotated, target)
			if got, found := SearchRotated(rotated, target); got != want || found != (want >= 0) {
				t.Fatalf("trial %d: SearchRotated(%v, %d) = (%d, %v), want %d", trial, rotated, target, got, found, want)
			}
		}
	}
}

// rotatedBenchData is benchData's even numbers rotated by n/3.
func rotatedBenchData(n int) (rotated, targets []int) {
	sorted, targets := benchData(n)
	return rotate(sorted, n/3), targets
}

func BenchmarkSearchRotated(b *testing.B) {
	for _, size := range benchSizes {
		rotated, targets := rotatedBenchData(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				x, _ := SearchRotated(rotated, targets[i&1023])
				sink += x
			}
			benchSink = sink
		})
	}
}
