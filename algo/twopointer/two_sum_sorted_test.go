package twopointer

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchSink keeps benchmark results observable.
var benchSink int

// benchSorted returns n fixed-seed ascending values in [0, 4n).
func benchSorted(n int) []int64 {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]int64, n)
	for i := range s {
		s[i] = int64(rng.IntN(4 * n))
	}
	slices.Sort(s)
	return s
}

// naiveTwoSum is the oracle: the smallest i and, for it, the largest j, sums in int.
func naiveTwoSum(sorted []int8, target int) (int, int, bool) {
	for i := range sorted {
		for j := len(sorted) - 1; j > i; j-- {
			if int(sorted[i])+int(sorted[j]) == target {
				return i, j, true
			}
		}
	}
	return -1, -1, false
}

// Laaksonen, Competitive Programmer's Handbook (2018), §8.1 "2SUM problem", pp. 78–79: in
// 1 4 5 6 7 9 9 10 with target 12 the pointers stop at 5 + 7 (indexes 2 and 4).
func TestTwoSumSorted_laaksonenExample(t *testing.T) {
	i, j, ok := TwoSumSorted([]int{1, 4, 5, 6, 7, 9, 9, 10}, 12)
	if !ok || i != 2 || j != 4 {
		t.Errorf("TwoSumSorted(laaksonen, 12) = (%d, %d, %v), want (2, 4, true)", i, j, ok)
	}
}

func TestTwoSumSorted_cases(t *testing.T) {
	tests := []struct {
		name   string
		sorted []int64
		target int64
		i, j   int
		ok     bool
	}{
		{"empty", nil, 0, -1, -1, false},
		{"one: no pair with itself", []int64{3}, 6, -1, -1, false},
		{"two", []int64{-2, 5}, 3, 0, 1, true},
		{"duplicates", []int64{3, 3}, 6, 0, 1, true},
		{"all equal: outermost pair", []int64{1, 1, 1, 1}, 2, 0, 3, true},
		{"none", []int64{1, 2, 4}, 7, -1, -1, false},
		{"negative and zero", []int64{-5, -1, 0, 4}, -1, 0, 3, true},
		{"extremes", []int64{math.MinInt64, 0, math.MaxInt64}, -1, 0, 2, true},
		{"overflow does not invent a pair", []int64{math.MaxInt64, math.MaxInt64}, -2, -1, -1, false},
		{"underflow does not invent a pair", []int64{math.MinInt64, math.MinInt64}, 0, -1, -1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, j, ok := TwoSumSorted(tt.sorted, tt.target)
			if i != tt.i || j != tt.j || ok != tt.ok {
				t.Errorf("TwoSumSorted(%v, %d) = (%d, %d, %v), want (%d, %d, %v)", tt.sorted, tt.target, i, j, ok, tt.i, tt.j, tt.ok)
			}
		})
	}
	if i, j, ok := TwoSumSorted([]uint8{100, 200}, 44); ok {
		t.Errorf("uint8 100+200 wraps to 44 but must not match: got (%d, %d, true)", i, j)
	}
	if i, j, ok := TwoSumSorted([]uint8{55, 200}, 255); !ok || i != 0 || j != 1 {
		t.Errorf("TwoSumSorted(uint8 55 200, 255) = (%d, %d, %v), want (0, 1, true)", i, j, ok)
	}
}

func TestTwoSumSorted_unsortedDoesNotPanic(t *testing.T) {
	_, _, _ = TwoSumSorted([]int{5, -3, 9, 0, 2}, 2)
}

func TestTwoSumSorted_noAllocation(t *testing.T) {
	s := benchSorted(1000)
	allocs := testing.AllocsPerRun(20, func() { _, _, _ = TwoSumSorted(s, -1) })
	if allocs != 0 {
		t.Errorf("TwoSumSorted allocates %v times per call, want 0", allocs)
	}
}

// Oracle: random sorted int8 data, targets across the whole int8 range, against the O(n²) search.
func TestTwoSumSorted_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		s := make([]int8, rng.IntN(15))
		for i := range s {
			s[i] = int8(rng.IntN(256) - 128)
			if rng.IntN(2) == 0 {
				s[i] /= 16
			}
		}
		slices.Sort(s)
		before := slices.Clone(s)
		target := int8(rng.IntN(256) - 128)
		wi, wj, wok := naiveTwoSum(s, int(target))
		if i, j, ok := TwoSumSorted(s, target); i != wi || j != wj || ok != wok {
			t.Fatalf("trial %d: TwoSumSorted(%v, %d) = (%d, %d, %v), want (%d, %d, %v)", trial, s, target, i, j, ok, wi, wj, wok)
		}
		if !slices.Equal(s, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkTwoSumSorted(b *testing.B) {
	for _, size := range benchSizes {
		s := benchSorted(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				i, _, _ := TwoSumSorted(s, -1) // absent: a full scan
				benchSink = i
			}
		})
	}
}
