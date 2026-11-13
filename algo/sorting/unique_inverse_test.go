package sorting

import (
	"slices"
	"strconv"
	"testing"
)

// Our example (label encoding by sorting, Knuth vol. 3 §5.2): red green blue green encode as blue=0,
// green=1, red=2, so the codes are 2 1 0 1.
func TestUniqueInverse_example(t *testing.T) {
	u, inv := UniqueInverse([]string{"red", "green", "blue", "green"})
	if !slices.Equal(u, []string{"blue", "green", "red"}) || !slices.Equal(inv, []int{2, 1, 0, 1}) {
		t.Errorf("UniqueInverse = %v %v", u, inv)
	}
}

func TestUniqueInverse_cases(t *testing.T) {
	tests := []struct {
		name       string
		values     []int
		u, inverse []int
	}{
		{"empty", []int{}, []int{}, []int{}},
		{"one", []int{9}, []int{9}, []int{0}},
		{"reversed", []int{3, 2, 1}, []int{1, 2, 3}, []int{2, 1, 0}},
		{"duplicates", []int{5, 5, 1}, []int{1, 5}, []int{1, 1, 0}},
	}
	for _, tc := range tests {
		u, inv := UniqueInverse(tc.values)
		if u == nil || inv == nil || !slices.Equal(u, tc.u) || !slices.Equal(inv, tc.inverse) {
			t.Errorf("%s: = %v %v, want %v %v", tc.name, u, inv, tc.u, tc.inverse)
		}
	}
}

func BenchmarkUniqueInverse(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchInts(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = UniqueInverse(values)
			}
		})
	}
}
