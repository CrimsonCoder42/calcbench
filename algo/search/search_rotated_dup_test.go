package search

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1 keys with every key doubled, rotated by every amount: 653 is
// found at one of its two positions.
func TestSearchRotatedDup_knuthExample(t *testing.T) {
	doubled := make([]int, 0, 32)
	for _, k := range knuthSorted {
		doubled = append(doubled, k, k)
	}
	for r := range len(doubled) {
		rotated := rotate(doubled, r)
		got, found := SearchRotatedDup(rotated, 653)
		if !found || rotated[got] != 653 {
			t.Errorf("rotation %d: SearchRotatedDup(653) = (%d, %v)", r, got, found)
		}
		if _, found := SearchRotatedDup(rotated, 400); found {
			t.Errorf("rotation %d: 400 found", r)
		}
	}
}

// Our own data, answers by hand: duplicates around the turn point, one odd value hidden among equal
// ones, and an absent value.
func TestSearchRotatedDup_exercises(t *testing.T) {
	if got, found := SearchRotatedDup([]int{5, 5, 7, 8, 5, 5, 5}, 8); !found || got != 3 {
		t.Errorf("SearchRotatedDup(8) = (%d, %v), want (3, true)", got, found)
	}
	if got, found := SearchRotatedDup([]int{4, 4, 4, 4, 3, 4}, 3); !found || got != 4 {
		t.Errorf("SearchRotatedDup(3) = (%d, %v), want (4, true)", got, found)
	}
	if _, found := SearchRotatedDup([]int{6, 8, 9, 1, 1, 3, 6}, 5); found {
		t.Error("SearchRotatedDup(5) found an absent value")
	}
}

func TestSearchRotatedDup_cases(t *testing.T) {
	tests := []struct {
		name    string
		rotated []int
		target  int
		found   bool
	}{
		{"empty", nil, 1, false},
		{"one element", []int{5}, 5, true},
		{"all equal, present", []int{7, 7, 7}, 7, true},
		{"all equal, absent", []int{7, 7, 7}, 6, false},
		{"not rotated with duplicates", []int{1, 1, 2, 3, 3}, 3, true},
		{"equal ends hide the turn", []int{3, 3, 1, 2, 3}, 1, true},
		{"equal ends, absent", []int{3, 3, 1, 2, 3}, 0, false},
		{"negative", []int{-1, -1, -5, -1}, -5, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := SearchRotatedDup(tt.rotated, tt.target)
			if found != tt.found || (found && tt.rotated[got] != tt.target) || (!found && got != -1) {
				t.Errorf("SearchRotatedDup(%v, %d) = (%d, %v), want found %v", tt.rotated, tt.target, got, found, tt.found)
			}
		})
	}
}

func TestSearchRotatedDup_noAllocation(t *testing.T) {
	rotated := []int{2, 2, 2, 3, 4, 2, 2}
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _ = SearchRotatedDup(rotated, 4) }); allocs != 0 {
		t.Errorf("SearchRotatedDup allocates %v times per call, want 0", allocs)
	}
}

func TestSearchRotatedDup_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		n := rng.IntN(40)
		sorted := randomSorted(rng, n, 1+rng.IntN(8))
		r := 0
		if n > 0 {
			r = rng.IntN(n)
		}
		rotated := rotate(sorted, r)
		for target := -1; target <= 9; target++ {
			want := slices.Contains(rotated, target)
			got, found := SearchRotatedDup(rotated, target)
			if found != want || (found && rotated[got] != target) || (!found && got != -1) {
				t.Fatalf("trial %d: SearchRotatedDup(%v, %d) = (%d, %v), want found %v", trial, rotated, target, got, found, want)
			}
		}
	}
}

// Two shapes: distinct keys (the O(log n) path) and all-equal keys but one (the O(n) worst case).
func BenchmarkSearchRotatedDup(b *testing.B) {
	for _, size := range benchSizes {
		rotated, targets := rotatedBenchData(size.n)
		b.Run("distinct/"+size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				x, _ := SearchRotatedDup(rotated, targets[i&1023])
				sink += x
			}
			benchSink = sink
		})
		flat := make([]int, size.n)
		for i := range flat {
			flat[i] = 1
		}
		flat[size.n/2] = 0
		b.Run("worstCase/"+size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = SearchRotatedDup(flat, 0)
			}
		})
	}
}
