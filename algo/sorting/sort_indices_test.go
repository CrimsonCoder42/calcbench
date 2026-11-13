package sorting

import (
	"cmp"
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Our own example: values 30 10 20 10 sort as 10 (index 1), 10 (index 3), 20 (index 2), 30 (index 0);
// the two 10s keep ascending index order (Sedgewick & Wayne §2.5, pointer sorting; stability §2.2).
func TestSortIndices_example(t *testing.T) {
	got, err := SortIndices(nil, []int{30, 10, 20, 10}, cmp.Compare[int])
	if err != nil || !slices.Equal(got, []int{1, 3, 2, 0}) {
		t.Errorf("SortIndices = (%v, %v), want [1 3 2 0]", got, err)
	}
}

func TestSortIndices_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   []int
	}{
		{"empty", []string{}, []int{}},
		{"one element", []string{"a"}, []int{0}},
		{"two reversed", []string{"b", "a"}, []int{1, 0}},
		{"all equal keeps index order", []string{"x", "x", "x"}, []int{0, 1, 2}},
		{"already sorted", []string{"a", "b", "c"}, []int{0, 1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortIndices(nil, tt.values, cmp.Compare[string])
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("SortIndices(%v) = (%v, %v), want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

// Newest first with equal times in input order: the shape of the consumer's feed ordering.
func TestSortIndices_descendingByTime(t *testing.T) {
	times := []int{100, 300, 200, 300, 100}
	newestFirst := func(a, b int) int { return cmp.Compare(b, a) }
	got, err := SortIndices(nil, times, newestFirst)
	if err != nil || !slices.Equal(got, []int{1, 3, 2, 0, 4}) {
		t.Errorf("SortIndices(newest first) = (%v, %v), want [1 3 2 0 4]", got, err)
	}
}

func TestSortIndices_errors(t *testing.T) {
	got, err := SortIndices([]int{9}, []int{1, 2}, nil)
	if !errors.Is(err, ErrNilCompare) || !slices.Equal(got, []int{9}) {
		t.Errorf("SortIndices(nil compare) = (%v, %v), want (dst, ErrNilCompare)", got, err)
	}
}

func TestSortIndices_inputKept(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, err := SortIndices(nil, values, cmp.Compare[int]); err != nil || !slices.Equal(values, knuth16) {
		t.Errorf("SortIndices modified its input: %v", values)
	}
}

// Oracle: a stable sort of (value, index) pairs.
func TestSortIndices_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(300)
		values := randomInts(rng, n, 1+n/5)
		want := make([]int, n)
		for i := range want {
			want[i] = i
		}
		slices.SortStableFunc(want, func(i, j int) int { return cmp.Compare(values[i], values[j]) })
		got, err := SortIndices(nil, values, cmp.Compare[int])
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d (n=%d): SortIndices differs from the stable oracle", trial, n)
		}
	}
}

func BenchmarkSortIndices(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n)
		dst := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				out, _ := SortIndices(dst, values, cmp.Compare[int])
				benchSink = out[0]
			}
		})
	}
}
