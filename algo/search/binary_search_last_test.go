package search

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1: 653 occurs once among the 16 sorted keys, so its last
// occurrence is also index 10. The unsuccessful search for 400 reports its insertion point, 5.
// Hand-derived.
func TestBinarySearchLast_knuthExample(t *testing.T) {
	if i, ok := BinarySearchLast(knuthSorted, 653); i != 10 || !ok {
		t.Errorf("BinarySearchLast(knuth, 653) = (%d, %v), want (10, true)", i, ok)
	}
	if i, ok := BinarySearchLast(knuthSorted, 400); i != 5 || ok {
		t.Errorf("BinarySearchLast(knuth, 400) = (%d, %v), want (5, false)", i, ok)
	}
}

func TestBinarySearchLast_cases(t *testing.T) {
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
		{"two elements, first", []int{1, 3}, 1, 0, true},
		{"duplicates, last of the run", []int{1, 2, 2, 2, 3}, 2, 3, true},
		{"all equal", []int{7, 7, 7, 7}, 7, 3, true},
		{"absent below every element", knuthSorted, 0, 0, false},
		{"negative duplicates", []int{-9, -5, -5, 0, 4}, -5, 2, true},
		{"minimum of the type", []int{math.MinInt, math.MinInt, math.MaxInt}, math.MinInt, 1, true},
		{"maximum of the type", []int{math.MinInt, math.MaxInt, math.MaxInt}, math.MaxInt, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, ok := BinarySearchLast(tt.sorted, tt.target)
			if i != tt.wantIndex || ok != tt.wantFound {
				t.Errorf("BinarySearchLast(%v, %d) = (%d, %v), want (%d, %v)",
					tt.sorted, tt.target, i, ok, tt.wantIndex, tt.wantFound)
			}
		})
	}
}

func TestBinarySearchLast_floats(t *testing.T) {
	tests := []struct {
		name      string
		target    float64
		wantIndex int
		wantFound bool
	}{
		{"NaN finds NaN", math.NaN(), 0, true},
		{"-0 finds the last zero", math.Copysign(0, -1), 4, true},
		{"+0 finds the last zero", 0, 4, true},
		{"-Inf", math.Inf(-1), 1, true},
		{"absent", 0.5, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, ok := BinarySearchLast(floatsSorted, tt.target)
			if i != tt.wantIndex || ok != tt.wantFound {
				t.Errorf("BinarySearchLast(floats, %v) = (%d, %v), want (%d, %v)", tt.target, i, ok, tt.wantIndex, tt.wantFound)
			}
		})
	}
}

func TestBinarySearchLast_noAllocation(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() { benchSink, _ = BinarySearchLast(knuthSorted, 653) })
	if allocs != 0 {
		t.Errorf("BinarySearchLast allocates %v times per call, want 0", allocs)
	}
}

// linearSearchLast is the slow obvious oracle: scan from the right.
func linearSearchLast(sorted []int, target int) (int, bool) {
	for i := len(sorted) - 1; i >= 0; i-- {
		if sorted[i] == target {
			return i, true
		}
		if sorted[i] < target {
			return i + 1, false
		}
	}
	return 0, false
}

func TestBinarySearchLast_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		sorted := randomSorted(rng, rng.IntN(64), 20)
		for target := -1; target <= 21; target++ {
			gi, gok := BinarySearchLast(sorted, target)
			wi, wok := linearSearchLast(sorted, target)
			if gi != wi || gok != wok {
				t.Fatalf("trial %d: BinarySearchLast(%v, %d) = (%d, %v), want (%d, %v)", trial, sorted, target, gi, gok, wi, wok)
			}
		}
	}
}

func BenchmarkBinarySearchLast(b *testing.B) {
	for _, size := range benchSizes {
		sorted, targets := benchData(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				j, _ := BinarySearchLast(sorted, targets[i&1023])
				sink += j
			}
			benchSink = sink
		})
	}
}
