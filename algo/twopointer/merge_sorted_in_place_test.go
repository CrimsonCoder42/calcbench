package twopointer

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Cormen et al., Introduction to Algorithms, 3rd ed., §2.3.1, Figure 2.3: MERGE of the sorted halves
// 2 4 5 7 and 1 2 3 6 gives 1 2 2 3 4 5 6 7.
func TestMergeSortedInPlace_clrsExample(t *testing.T) {
	buf := []int{2, 4, 5, 7, 0, 0, 0, 0}
	if err := MergeSortedInPlace(buf, 4, []int{1, 2, 3, 6}); err != nil {
		t.Fatal(err)
	}
	if want := []int{1, 2, 2, 3, 4, 5, 6, 7}; !slices.Equal(buf, want) {
		t.Errorf("MergeSortedInPlace = %v, want %v", buf, want)
	}
}

func TestMergeSortedInPlace_cases(t *testing.T) {
	tests := []struct {
		name  string
		buf   []int
		m     int
		other []int
		want  []int
	}{
		{"both empty", nil, 0, nil, nil},
		{"other empty", []int{1, 2}, 2, nil, []int{1, 2}},
		{"buf part empty", []int{0, 0}, 0, []int{3, 4}, []int{3, 4}},
		{"one each", []int{5, 0}, 1, []int{4}, []int{4, 5}},
		{"all of other first", []int{7, 8, 0, 0}, 2, []int{1, 2}, []int{1, 2, 7, 8}},
		{"all of buf first", []int{1, 2, 0, 0}, 2, []int{7, 8}, []int{1, 2, 7, 8}},
		{"duplicates", []int{1, 1, 0, 0}, 2, []int{1, 1}, []int{1, 1, 1, 1}},
		{"negative and zero", []int{-3, 0, 0, 0}, 2, []int{-5, -1}, []int{-5, -3, -1, 0}},
		{"extremes", []int{math.MinInt, math.MaxInt, 0}, 2, []int{0}, []int{math.MinInt, 0, math.MaxInt}},
		{"room beyond the result is untouched", []int{2, 0, 9}, 1, []int{1}, []int{1, 2, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := slices.Clone(tt.buf)
			other := slices.Clone(tt.other)
			if err := MergeSortedInPlace(buf, tt.m, other); err != nil || !slices.Equal(buf, tt.want) {
				t.Errorf("MergeSortedInPlace(%v, %d, %v) = %v, %v; want %v", tt.buf, tt.m, tt.other, buf, err, tt.want)
			}
			if !slices.Equal(other, tt.other) {
				t.Errorf("other modified: %v", other)
			}
		})
	}
}

// Stability is visible with -0 and +0, which compare equal: buf's element must come first. NaN sorts
// before everything in cmp.Compare order.
func TestMergeSortedInPlace_stableAndNaN(t *testing.T) {
	negZero := math.Copysign(0, -1)
	buf := []float64{negZero, 1, 0, 0, 0}
	if err := MergeSortedInPlace(buf, 2, []float64{math.NaN(), 0, 2}); err != nil {
		t.Fatal(err)
	}
	if !math.IsNaN(buf[0]) || !math.Signbit(buf[1]) || buf[2] != 0 || math.Signbit(buf[2]) || buf[3] != 1 || buf[4] != 2 {
		t.Errorf("MergeSortedInPlace(-0 1 | NaN 0 2) = %v, want [NaN -0 0 1 2] with -0 before +0", buf)
	}
}

func TestMergeSortedInPlace_errors(t *testing.T) {
	tests := []struct {
		name  string
		buf   []int
		m     int
		other []int
	}{
		{"negative m", []int{1, 2}, -1, nil},
		{"m beyond buf", []int{1, 2}, 3, nil},
		{"no room", []int{1, 2, 0}, 2, []int{3, 4}},
		{"m at the minimum of int", []int{1}, math.MinInt, nil},
		{"m at the maximum of int", []int{1}, math.MaxInt, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := slices.Clone(tt.buf)
			if err := MergeSortedInPlace(buf, tt.m, tt.other); !errors.Is(err, ErrLength) || !slices.Equal(buf, tt.buf) {
				t.Errorf("MergeSortedInPlace(%v, %d, %v) = %v, buf %v; want ErrLength, buf unchanged", tt.buf, tt.m, tt.other, err, buf)
			}
		})
	}
}

func TestMergeSortedInPlace_noAllocation(t *testing.T) {
	a, other := benchSorted(500), benchSorted(500)
	buf := make([]int64, 1000)
	allocs := testing.AllocsPerRun(20, func() {
		copy(buf, a)
		_ = MergeSortedInPlace(buf, len(a), other)
	})
	if allocs != 0 {
		t.Errorf("MergeSortedInPlace allocates %v times per call, want 0", allocs)
	}
}

// Oracle: the merge equals sorting the concatenation.
func TestMergeSortedInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		m, k := rng.IntN(15), rng.IntN(15)
		a := make([]int, m)
		other := make([]int, k)
		for i := range a {
			a[i] = rng.IntN(20) - 10
		}
		for i := range other {
			other[i] = rng.IntN(20) - 10
		}
		slices.Sort(a)
		slices.Sort(other)
		want := slices.Sorted(slices.Values(append(slices.Clone(a), other...)))
		buf := append(slices.Clone(a), make([]int, k+rng.IntN(3))...)
		if err := MergeSortedInPlace(buf, m, other); err != nil || !slices.Equal(buf[:m+k], want) {
			t.Fatalf("trial %d: merge %v + %v = %v (%v), want %v", trial, a, other, buf[:m+k], err, want)
		}
	}
}

func BenchmarkMergeSortedInPlace(b *testing.B) {
	for _, size := range benchSizes {
		a, other := benchSorted(size.n/2), benchSorted(size.n-size.n/2)
		buf := make([]int64, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				copy(buf, a)
				_ = MergeSortedInPlace(buf, len(a), other)
			}
		})
	}
}
