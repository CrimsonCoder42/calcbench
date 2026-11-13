package sorting

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// CLRS 3rd ed. §8.2 (Figure 8.2) runs COUNTING-SORT on eight keys in 0..5. Our own eight keys of the
// same shape; the expected output is hand-derived.
func TestCountingSort_clrsShapedExample(t *testing.T) {
	got, err := CountingSort(nil, []int{2, 5, 3, 0, 2, 3, 0, 3}, 5)
	want := []int{0, 0, 2, 2, 3, 3, 3, 5}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("CountingSort = (%v, %v), want %v", got, err, want)
	}
}

func TestCountingSort_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		maxKey int
		want   []int
	}{
		{"empty", []int{}, 0, []int{}},
		{"one element", []int{4}, 4, []int{4}},
		{"two elements reversed", []int{1, 0}, 1, []int{0, 1}},
		{"all zero", []int{0, 0, 0}, 0, []int{0, 0, 0}},
		{"already sorted", []int{0, 1, 2, 3}, 9, []int{0, 1, 2, 3}},
		{"reversed", []int{3, 2, 1, 0}, 3, []int{0, 1, 2, 3}},
		{"duplicates", []int{7, 7, 1, 7}, 7, []int{1, 7, 7, 7}},
		// The count table is maxKey+1 ints whatever the input, so maxCountingKey itself (128 MiB) is not
		// tested here; 1<<16 (512 KiB) exercises a wide key range.
		{"wide key range", []int{1 << 16, 0}, 1 << 16, []int{0, 1 << 16}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountingSort(nil, tt.values, tt.maxKey)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("CountingSort(%v, %d) = (%v, %v), want %v", tt.values, tt.maxKey, got, err, tt.want)
			}
		})
	}
}

func TestCountingSort_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		maxKey int
	}{
		{"maxKey negative", []int{1}, -1},
		{"maxKey above the cap", []int{1}, maxCountingKey + 1},
		{"maxKey at the maximum of int", []int{1}, math.MaxInt},
		{"value negative", []int{1, -1}, 5},
		{"value above maxKey", []int{1, 6}, 5},
		{"value at the minimum of int", []int{math.MinInt}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []int{9, 9}
			got, err := CountingSort(dst, tt.values, tt.maxKey)
			if !errors.Is(err, ErrKeyRange) || !slices.Equal(got, []int{9, 9}) {
				t.Errorf("CountingSort(%v, %d) = (%v, %v), want (dst unwritten, ErrKeyRange)", tt.values, tt.maxKey, got, err)
			}
		})
	}
}

func TestCountingSort_inPlaceAndInputKept(t *testing.T) {
	values := []int{3, 1, 2}
	if _, err := CountingSort(nil, values, 3); err != nil || !slices.Equal(values, []int{3, 1, 2}) {
		t.Errorf("CountingSort modified its input: %v", values)
	}
	if _, err := CountingSort(values, values, 3); err != nil || !slices.Equal(values, []int{1, 2, 3}) {
		t.Errorf("in place: %v", values)
	}
}

func TestCountingSort_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		maxKey := rng.IntN(100)
		values := randomInts(rng, rng.IntN(300), maxKey+1)
		want := slices.Sorted(slices.Values(values))
		got, err := CountingSort(nil, values, maxKey)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: CountingSort differs from slices.Sort", trial)
		}
	}
}

// Keys in [0, n), so k = n and the cost is O(n).
func BenchmarkCountingSort(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n)
		dst := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				out, _ := CountingSort(dst, values, size.n-1)
				benchSink = out[0]
			}
		})
	}
}
