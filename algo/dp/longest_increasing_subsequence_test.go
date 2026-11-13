package dp

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// checkIncreasing reports whether idx is ascending and values[idx[k]] strictly increases.
func checkIncreasing[T cmp.Ordered](values []T, idx []int) bool {
	for k := range idx {
		if idx[k] < 0 || idx[k] >= len(values) {
			return false
		}
		if k > 0 && (idx[k] <= idx[k-1] || cmp.Compare(values[idx[k-1]], values[idx[k]]) >= 0) {
			return false
		}
	}
	return true
}

// lisLengthSlow is the O(n²) recurrence L(i) = 1 + max{L(j) : j < i, v[j] < v[i]}.
func lisLengthSlow(values []int) int {
	best := 0
	l := make([]int, len(values))
	for i := range values {
		l[i] = 1
		for j := range i {
			if values[j] < values[i] {
				l[i] = max(l[i], l[j]+1)
			}
		}
		best = max(best, l[i])
	}
	return best
}

// Aldous, Diaconis 1999, §1.1 (patience sorting): the deck 7 2 8 1 3 4 10 6 9 5 deals into 5 piles,
// so its longest increasing subsequence has length 5. The trace by hand gives 1 3 4 6 9.
func TestLongestIncreasingSubsequence_aldousDiaconisDeck(t *testing.T) {
	deck := []int{7, 2, 8, 1, 3, 4, 10, 6, 9, 5}
	got := LongestIncreasingSubsequence(deck)
	if want := []int{3, 4, 5, 7, 8}; !slices.Equal(got, want) {
		t.Errorf("LongestIncreasingSubsequence(deck) = %v, want %v", got, want)
	}
}

func TestLongestIncreasingSubsequence_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty", nil, []int{}},
		{"one", []int{4}, []int{0}},
		{"two rising", []int{1, 2}, []int{0, 1}},
		{"two falling: the later, smaller one", []int{2, 1}, []int{1}},
		{"duplicates are not increasing", []int{3, 3, 3}, []int{2}},
		{"sorted", []int{1, 2, 3, 4}, []int{0, 1, 2, 3}},
		{"reversed", []int{4, 3, 2, 1}, []int{3}},
		{"negative and zero", []int{-5, 0, -3, -1, 2}, []int{0, 2, 3, 4}},
		{"extremes", []int{math.MaxInt, math.MinInt, 0}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LongestIncreasingSubsequence(tt.values)
			if got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("LongestIncreasingSubsequence(%v) = %#v, want %v", tt.values, got, tt.want)
			}
		})
	}
}

func TestLongestIncreasingSubsequence_floatsAndStrings(t *testing.T) {
	// cmp.Compare order: NaN is below every other float, so it can start the subsequence.
	floats := []float64{math.NaN(), 1, math.Inf(-1), 2, math.Inf(1)}
	if got := LongestIncreasingSubsequence(floats); len(got) != 4 || !checkIncreasing(floats, got) {
		t.Errorf("floats: got %v, want 4 increasing indices", got)
	}
	words := []string{"pear", "apple", "fig", "plum"}
	if got := LongestIncreasingSubsequence(words); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("strings: got %v, want [1 2 3]", got)
	}
}

func TestLongestIncreasingSubsequence_doesNotModifyInput(t *testing.T) {
	values := []int{5, 1, 4, 2, 3}
	orig := slices.Clone(values)
	LongestIncreasingSubsequence(values)
	if !slices.Equal(values, orig) {
		t.Errorf("input modified: %v", values)
	}
}

func TestLongestIncreasingSubsequence_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(60)
		values := make([]int, n)
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(2*n+1))
		}
		got := LongestIncreasingSubsequence(values)
		if !checkIncreasing(values, got) || len(got) != lisLengthSlow(values) {
			t.Fatalf("trial %d: %v → %v, want an increasing subsequence of length %d", trial, values, got, lisLengthSlow(values))
		}
	}
}

func BenchmarkLongestIncreasingSubsequence(b *testing.B) {
	for _, size := range linearSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		values := make([]int, size.n)
		for i := range values {
			values[i] = rng.IntN(size.n)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = LongestIncreasingSubsequence(values)
			}
		})
	}
}
