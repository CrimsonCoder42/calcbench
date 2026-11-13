package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveAtMostKDistinct is the oracle: every run, longest first then earliest start, O(n³).
func naiveAtMostKDistinct(values []int, k int) (int, int) {
	for length := len(values); length >= 1; length-- {
		for start := 0; start+length <= len(values); start++ {
			seen := map[int]bool{}
			for _, v := range values[start : start+length] {
				seen[v] = true
			}
			if len(seen) <= k {
				return start, length
			}
		}
	}
	return 0, 0
}

// The two-pointers method of Laaksonen, Competitive Programmer's Handbook (2018), §8.1, p. 77, on
// hand-made data: in 1 2 1 3 4 3 5 3 with k = 2, the longest runs are 1 2 1 and 3 4 3 (length 3) and
// 3 5 3 (length 3); the first starts at 0. Hand-derived.
func TestLongestWithAtMostKDistinct_worked(t *testing.T) {
	start, length, err := LongestWithAtMostKDistinct([]int{1, 2, 1, 3, 4, 3, 5, 3}, 2)
	if err != nil || start != 0 || length != 3 {
		t.Errorf("LongestWithAtMostKDistinct(…, 2) = (%d, %d, %v), want (0, 3, nil)", start, length, err)
	}
}

func TestLongestWithAtMostKDistinct_cases(t *testing.T) {
	tests := []struct {
		name          string
		values        []int
		k             int
		start, length int
	}{
		{"empty", nil, 3, 0, 0},
		{"k zero", []int{1, 2}, 0, 0, 0},
		{"one", []int{9}, 1, 0, 1},
		{"two different, k 1", []int{8, 9}, 1, 0, 1},
		{"all equal", []int{4, 4, 4}, 1, 0, 3},
		{"k at least the distinct count", []int{1, 2, 3}, 3, 0, 3},
		{"k beyond n", []int{1, 2, 3}, math.MaxInt, 0, 3},
		{"later run longer", []int{1, 2, 3, 3, 3, 2}, 2, 1, 5},
		{"sorted distinct, k 2", []int{1, 2, 3, 4}, 2, 0, 2},
		{"negative and extremes", []int{math.MinInt, -1, math.MinInt, math.MaxInt}, 2, 0, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, length, err := LongestWithAtMostKDistinct(tt.values, tt.k)
			if err != nil || start != tt.start || length != tt.length {
				t.Errorf("LongestWithAtMostKDistinct(%v, %d) = (%d, %d, %v), want (%d, %d, nil)", tt.values, tt.k, start, length, err, tt.start, tt.length)
			}
		})
	}
}

// Each NaN is its own distinct value; -0 and +0 are one value.
func TestLongestWithAtMostKDistinct_floats(t *testing.T) {
	nan := math.NaN()
	start, length, err := LongestWithAtMostKDistinct([]float64{nan, nan, 1, math.Copysign(0, -1), 0, 0}, 2)
	if err != nil || start != 2 || length != 4 {
		t.Errorf("floats = (%d, %d, %v), want (2, 4, nil)", start, length, err)
	}
	start, length, _ = LongestWithAtMostKDistinct([]float64{nan, nan, nan}, 2)
	if start != 0 || length != 2 {
		t.Errorf("three NaNs, k 2 = (%d, %d), want (0, 2)", start, length)
	}
}

func TestLongestWithAtMostKDistinct_errors(t *testing.T) {
	for _, k := range []int{-1, math.MinInt} {
		if _, _, err := LongestWithAtMostKDistinct([]int{1}, k); !errors.Is(err, ErrNegativeK) {
			t.Errorf("k = %d: error = %v, want ErrNegativeK", k, err)
		}
	}
}

func TestLongestWithAtMostKDistinct_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int, rng.IntN(25))
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(6))
		}
		before := slices.Clone(values)
		k := rng.IntN(5)
		wantStart, wantLength := naiveAtMostKDistinct(values, k)
		start, length, err := LongestWithAtMostKDistinct(values, k)
		if err != nil || start != wantStart || length != wantLength {
			t.Fatalf("trial %d: LongestWithAtMostKDistinct(%v, %d) = (%d, %d, %v), want (%d, %d)", trial, values, k, start, length, err, wantStart, wantLength)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

// benchSmallAlphabet returns n fixed-seed values in [0, 26).
func benchSmallAlphabet(n int) []int {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]int, n)
	for i := range s {
		s[i] = rng.IntN(26)
	}
	return s
}

func BenchmarkLongestWithAtMostKDistinct(b *testing.B) {
	for _, size := range benchSizes {
		values := benchSmallAlphabet(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, length, _ := LongestWithAtMostKDistinct(values, 5)
				benchSinkDeque = length
			}
		})
	}
}
