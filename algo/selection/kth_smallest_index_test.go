package selection

import (
	"cmp"
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Our example (CLRS §9.2 selection, applied to indexes): in 50 10 40 10 30 the 2nd smallest is the 10 at
// index 3; perm lists the indexes below it (1), equal to it (1, 3 in index order), then above (0, 2, 4).
func TestKthSmallestIndex_example(t *testing.T) {
	idx, perm, err := KthSmallestIndex(nil, []int{50, 10, 40, 10, 30}, 2)
	if err != nil || idx != 3 || !slices.Equal(perm, []int{1, 3, 0, 2, 4}) {
		t.Errorf("KthSmallestIndex = (%d, %v, %v), want (3, [1 3 0 2 4])", idx, perm, err)
	}
}

func TestKthSmallestIndex_cases(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		k        int
		wantIdx  int
		wantPerm []int
	}{
		{"one", []float64{7}, 1, 0, []int{0}},
		{"min", []float64{3, 1, 2}, 1, 1, []int{1, 0, 2}},
		{"max", []float64{3, 1, 2}, 3, 0, []int{1, 2, 0}},
		{"all equal", []float64{5, 5, 5}, 2, 1, []int{0, 1, 2}},
		{"NaN sorts first", []float64{1, math.NaN(), 0}, 1, 1, []int{1, 0, 2}},
		{"signed zeros equal", []float64{0, math.Copysign(0, -1), -1}, 2, 0, []int{2, 0, 1}},
	}
	for _, tc := range tests {
		idx, perm, err := KthSmallestIndex(nil, tc.values, tc.k)
		if err != nil || idx != tc.wantIdx || !slices.Equal(perm, tc.wantPerm) {
			t.Errorf("%s: = (%d, %v, %v), want (%d, %v)", tc.name, idx, perm, err, tc.wantIdx, tc.wantPerm)
		}
	}
}

func TestKthSmallestIndex_errors(t *testing.T) {
	if _, _, err := KthSmallestIndex(nil, []int{}, 1); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: err = %v", err)
	}
	for _, k := range []int{0, 4, -1} {
		if _, _, err := KthSmallestIndex(nil, []int{1, 2, 3}, k); !errors.Is(err, ErrRank) {
			t.Errorf("k = %d: err = %v", k, err)
		}
	}
}

// Oracle: perm equals the indexes grouped into <, ==, > blocks around the k-th value, each block in ascending
// index order; the k-th value comes from a stable sort of the indexes, on 300 random inputs with many ties.
func TestKthSmallestIndex_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := 1 + r.IntN(100)
		values := make([]int, n)
		for i := range values {
			values[i] = r.IntN(10)
		}
		k := 1 + r.IntN(n)
		idx, perm, err := KthSmallestIndex(nil, values, k)
		if err != nil {
			t.Fatal(err)
		}
		sorted := make([]int, n)
		for i := range sorted {
			sorted[i] = i
		}
		slices.SortStableFunc(sorted, func(a, b int) int { return cmp.Compare(values[a], values[b]) })
		kth := values[sorted[k-1]]
		if values[idx] != kth {
			t.Fatalf("values[%d] = %d, want %d", idx, values[idx], kth)
		}
		// Each block (<, ==, >) is in ascending index order: walk i = 0…n−1 and append to its block.
		var lessBlock, equalBlock, greaterBlock []int
		for i := range n {
			switch cmp.Compare(values[i], kth) {
			case -1:
				lessBlock = append(lessBlock, i)
			case 0:
				equalBlock = append(equalBlock, i)
			default:
				greaterBlock = append(greaterBlock, i)
			}
		}
		want := slices.Concat(lessBlock, equalBlock, greaterBlock)
		if !slices.Equal(perm, want) {
			t.Fatalf("perm = %v, want %v", perm, want)
		}
	}
}

func BenchmarkKthSmallestIndex(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		r := rand.New(rand.NewPCG(3, 4))
		values := make([]float64, n)
		for i := range values {
			values[i] = r.Float64()
		}
		dst := make([]int, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = KthSmallestIndex(dst, values, n/2)
			}
		})
	}
}
