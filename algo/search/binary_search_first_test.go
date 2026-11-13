package search

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1: the search for K = 653 in the 16 sorted keys succeeds at the
// 11th key (index 10). The unsuccessful search for 400 reports its insertion point, 5. Hand-derived.
func TestBinarySearchFirst_knuthExample(t *testing.T) {
	if i, ok := BinarySearchFirst(knuthSorted, 653); i != 10 || !ok {
		t.Errorf("BinarySearchFirst(knuth, 653) = (%d, %v), want (10, true)", i, ok)
	}
	if i, ok := BinarySearchFirst(knuthSorted, 400); i != 5 || ok {
		t.Errorf("BinarySearchFirst(knuth, 400) = (%d, %v), want (5, false)", i, ok)
	}
}

func TestBinarySearchFirst_cases(t *testing.T) {
	tests := []struct {
		name      string
		sorted    []int
		target    int
		wantIndex int
		wantFound bool
	}{
		{"empty", nil, 5, 0, false},
		{"one element, found", []int{5}, 5, 0, true},
		{"one element, absent below", []int{5}, 4, 0, false},
		{"one element, absent above", []int{5}, 6, 1, false},
		{"two elements, second", []int{1, 3}, 3, 1, true},
		{"duplicates, first of the run", []int{1, 2, 2, 2, 3}, 2, 1, true},
		{"all equal", []int{7, 7, 7, 7}, 7, 0, true},
		{"absent above every element", knuthSorted, 1000, 16, false},
		{"negative duplicates", []int{-9, -5, -5, 0, 4}, -5, 1, true},
		{"minimum of the type", []int{math.MinInt, math.MinInt, math.MaxInt}, math.MinInt, 0, true},
		{"maximum of the type", []int{math.MinInt, math.MaxInt, math.MaxInt}, math.MaxInt, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, ok := BinarySearchFirst(tt.sorted, tt.target)
			if i != tt.wantIndex || ok != tt.wantFound {
				t.Errorf("BinarySearchFirst(%v, %d) = (%d, %v), want (%d, %v)",
					tt.sorted, tt.target, i, ok, tt.wantIndex, tt.wantFound)
			}
		})
	}
}

func TestBinarySearchFirst_floats(t *testing.T) {
	tests := []struct {
		name      string
		target    float64
		wantIndex int
		wantFound bool
	}{
		{"NaN finds NaN", math.NaN(), 0, true},
		{"-0 finds the first zero", math.Copysign(0, -1), 3, true},
		{"+0 finds the first zero", 0, 3, true},
		{"+Inf", math.Inf(1), 6, true},
		{"absent", 0.5, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, ok := BinarySearchFirst(floatsSorted, tt.target)
			if i != tt.wantIndex || ok != tt.wantFound {
				t.Errorf("BinarySearchFirst(floats, %v) = (%d, %v), want (%d, %v)", tt.target, i, ok, tt.wantIndex, tt.wantFound)
			}
		})
	}
}

func TestBinarySearchFirst_noAllocation(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() { benchSink, _ = BinarySearchFirst(knuthSorted, 653) })
	if allocs != 0 {
		t.Errorf("BinarySearchFirst allocates %v times per call, want 0", allocs)
	}
}

// linearSearchFirst is the slow obvious oracle.
func linearSearchFirst(sorted []int, target int) (int, bool) {
	for i, v := range sorted {
		if v == target {
			return i, true
		}
		if v > target {
			return i, false
		}
	}
	return len(sorted), false
}

func TestBinarySearchFirst_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		sorted := randomSorted(rng, rng.IntN(64), 20)
		for target := -1; target <= 21; target++ {
			gi, gok := BinarySearchFirst(sorted, target)
			wi, wok := linearSearchFirst(sorted, target)
			if gi != wi || gok != wok {
				t.Fatalf("trial %d: BinarySearchFirst(%v, %d) = (%d, %v), want (%d, %v)", trial, sorted, target, gi, gok, wi, wok)
			}
		}
	}
}

func BenchmarkBinarySearchFirst(b *testing.B) {
	for _, size := range benchSizes {
		sorted, targets := benchData(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				j, _ := BinarySearchFirst(sorted, targets[i&1023])
				sink += j
			}
			benchSink = sink
		})
	}
}
