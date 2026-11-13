package window

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveWithoutRepeat is the oracle: every run, longest first then earliest start, O(n³).
func naiveWithoutRepeat(values []int) (int, int) {
	for length := len(values); length >= 1; length-- {
		for start := 0; start+length <= len(values); start++ {
			seen := map[int]bool{}
			ok := true
			for _, v := range values[start : start+length] {
				if seen[v] {
					ok = false
					break
				}
				seen[v] = true
			}
			if ok {
				return start, length
			}
		}
	}
	return 0, 0
}

// The two-pointers method of Laaksonen, Competitive Programmer's Handbook (2018), §8.1, p. 77, on
// hand-made data: in 4 1 4 2 3 1 2 the longest run without a repeat is 1 4 2 3 (start 1, length 4).
// Hand-derived.
func TestLongestWithoutRepeat_worked(t *testing.T) {
	if start, length := LongestWithoutRepeat([]int{4, 1, 4, 2, 3, 1, 2}); start != 1 || length != 4 {
		t.Errorf("LongestWithoutRepeat(4 1 4 2 3 1 2) = (%d, %d), want (1, 4)", start, length)
	}
}

func TestLongestWithoutRepeat_cases(t *testing.T) {
	tests := []struct {
		name          string
		values        []int
		start, length int
	}{
		{"empty", nil, 0, 0},
		{"one", []int{7}, 0, 1},
		{"two equal", []int{7, 7}, 0, 1},
		{"two different", []int{7, 8}, 0, 2},
		{"all equal", []int{3, 3, 3}, 0, 1},
		{"all distinct", []int{1, 2, 3, 4}, 0, 4},
		{"repeat seen before the window does not move it", []int{1, 2, 1, 3, 2}, 1, 3},
		{"tie: earliest start", []int{1, 2, 2, 1}, 0, 2},
		{"extremes", []int{math.MinInt, math.MaxInt, math.MinInt}, 0, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if start, length := LongestWithoutRepeat(tt.values); start != tt.start || length != tt.length {
				t.Errorf("LongestWithoutRepeat(%v) = (%d, %d), want (%d, %d)", tt.values, start, length, tt.start, tt.length)
			}
		})
	}
	if start, length := LongestWithoutRepeat([]string{"a", "b", "a", "c"}); start != 1 || length != 3 {
		t.Errorf("strings = (%d, %d), want (1, 3)", start, length)
	}
}

// NaN never repeats; -0 and +0 are the same value.
func TestLongestWithoutRepeat_floats(t *testing.T) {
	nan := math.NaN()
	if start, length := LongestWithoutRepeat([]float64{nan, nan, nan}); start != 0 || length != 3 {
		t.Errorf("three NaNs = (%d, %d), want (0, 3)", start, length)
	}
	if start, length := LongestWithoutRepeat([]float64{math.Copysign(0, -1), 0}); start != 0 || length != 1 {
		t.Errorf("-0 0 = (%d, %d), want (0, 1)", start, length)
	}
}

func TestLongestWithoutRepeat_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int, rng.IntN(25))
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(10))
		}
		before := slices.Clone(values)
		wantStart, wantLength := naiveWithoutRepeat(values)
		if start, length := LongestWithoutRepeat(values); start != wantStart || length != wantLength {
			t.Fatalf("trial %d: LongestWithoutRepeat(%v) = (%d, %d), want (%d, %d)", trial, values, start, length, wantStart, wantLength)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkLongestWithoutRepeat(b *testing.B) {
	for _, size := range benchSizes {
		values := benchSmallAlphabet(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, length := LongestWithoutRepeat(values)
				benchSinkDeque = length
			}
		})
	}
}
