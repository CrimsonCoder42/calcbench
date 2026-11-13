package sorting

import (
	"slices"
	"strconv"
	"testing"
)

// Our example (Knuth vol. 3 §5.2): "b a c a b a" has a×3, b×2, c×1.
func TestUniqueCounts_example(t *testing.T) {
	u, c := UniqueCounts([]string{"b", "a", "c", "a", "b", "a"})
	if !slices.Equal(u, []string{"a", "b", "c"}) || !slices.Equal(c, []int{3, 2, 1}) {
		t.Errorf("UniqueCounts = %v %v", u, c)
	}
}

func TestUniqueCounts_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		u, c   []int
	}{
		{"empty", []int{}, []int{}, []int{}},
		{"one", []int{4}, []int{4}, []int{1}},
		{"all equal", []int{2, 2, 2}, []int{2}, []int{3}},
		{"sorted distinct", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 1, 1}},
		{"negative", []int{-1, 5, -1}, []int{-1, 5}, []int{2, 1}},
	}
	for _, tc := range tests {
		u, c := UniqueCounts(tc.values)
		if u == nil || c == nil || !slices.Equal(u, tc.u) || !slices.Equal(c, tc.c) {
			t.Errorf("%s: = %v %v, want %v %v", tc.name, u, c, tc.u, tc.c)
		}
	}
}

func BenchmarkUniqueCounts(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchInts(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = UniqueCounts(values)
			}
		})
	}
}
