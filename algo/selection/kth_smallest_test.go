package selection

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// knuth16 is the 16-key running example of Knuth, TAOCP Vol. 3, 2nd ed. (ch. 5), in Knuth's original
// unsorted order. Sorted: 61 87 154 170 275 426 503 509 512 612 653 677 703 765 897 908.
// Shared by every test file in this package.
var knuth16 = []int{503, 87, 512, 61, 908, 170, 897, 275, 653, 426, 154, 509, 612, 677, 765, 703}

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchSink keeps benchmark results observable.
var benchSink int

// randomInts returns n fixed-seed values in [0, maxValue).
func randomInts(rng *rand.Rand, n, maxValue int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = rng.IntN(maxValue)
	}
	return s
}

// Worked example on the data of Knuth, TAOCP Vol. 3, 2nd ed., ch. 5 (the 16 keys). Expected order
// statistics are read off the sorted list by hand.
func TestKthSmallest_knuthExample(t *testing.T) {
	for k, want := range map[int]int{1: 61, 5: 275, 8: 509, 9: 512, 16: 908} {
		got, err := KthSmallest(knuth16, k)
		if err != nil || got != want {
			t.Errorf("KthSmallest(knuth16, %d) = (%d, %v), want (%d, nil)", k, got, err, want)
		}
	}
}

func TestKthSmallest_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		k      int
		want   int
	}{
		{"one element", []int{42}, 1, 42},
		{"two elements, first", []int{9, 3}, 1, 3},
		{"two elements, second", []int{9, 3}, 2, 9},
		{"duplicates count separately", []int{4, 1, 4, 1, 4}, 3, 4},
		{"duplicates, rank inside the run", []int{4, 1, 4, 1, 4}, 2, 1},
		{"all equal", []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7}, 11, 7},
		{"already sorted", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, 13, 12},
		{"reversed", []int{19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, 13, 12},
		{"negative and zero", []int{-3, 0, -7, 5, -1}, 2, -3},
		{"extremes of the type", []int{math.MaxInt, math.MinInt, 0}, 1, math.MinInt},
		{"maximum of the type", []int{math.MaxInt, math.MinInt, 0}, 3, math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KthSmallest(tt.values, tt.k)
			if err != nil || got != tt.want {
				t.Errorf("KthSmallest(%v, %d) = (%d, %v), want (%d, nil)", tt.values, tt.k, got, err, tt.want)
			}
		})
	}
}

func TestKthSmallest_floats(t *testing.T) {
	values := []float64{1, math.NaN(), math.Inf(1), math.Copysign(0, -1), math.Inf(-1), 0.5}
	// cmp.Compare order: NaN, -Inf, -0, 0.5, 1, +Inf.
	want := []float64{math.NaN(), math.Inf(-1), 0, 0.5, 1, math.Inf(1)}
	for k := 1; k <= len(values); k++ {
		got, err := KthSmallest(values, k)
		w := want[k-1]
		same := got == w || (math.IsNaN(got) && math.IsNaN(w))
		if err != nil || !same {
			t.Errorf("KthSmallest(floats, %d) = (%v, %v), want (%v, nil)", k, got, err, w)
		}
	}
}

func TestKthSmallest_strings(t *testing.T) {
	got, err := KthSmallest([]string{"pear", "apple", "fig", "banana"}, 2)
	if err != nil || got != "banana" {
		t.Errorf("KthSmallest(fruit, 2) = (%q, %v), want (\"banana\", nil)", got, err)
	}
}

func TestKthSmallest_errors(t *testing.T) {
	tests := []struct {
		name    string
		values  []int
		k       int
		wantErr error
	}{
		{"nil input", nil, 1, ErrEmpty},
		{"empty input", []int{}, 1, ErrEmpty},
		{"k zero", []int{1, 2, 3}, 0, ErrRank},
		{"k negative", []int{1, 2, 3}, -1, ErrRank},
		{"k above n", []int{1, 2, 3}, 4, ErrRank},
		{"k at the minimum of int", []int{1, 2, 3}, math.MinInt, ErrRank},
		{"k at the maximum of int", []int{1, 2, 3}, math.MaxInt, ErrRank},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KthSmallest(tt.values, tt.k)
			if !errors.Is(err, tt.wantErr) || got != 0 {
				t.Errorf("KthSmallest(%v, %d) = (%d, %v), want (0, %v)", tt.values, tt.k, got, err, tt.wantErr)
			}
		})
	}
}

func TestKthSmallest_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, err := KthSmallest(values, 8); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, knuth16) {
		t.Errorf("KthSmallest modified its input: %v", values)
	}
}

// The oracle sorts a copy and indexes it. Small value ranges force long runs of duplicates; the
// lengths cross insertionMax so both the insertion-sort finish and the partition rounds are exercised.
func TestKthSmallest_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := 1 + rng.IntN(300)
		values := randomInts(rng, n, 1+rng.IntN(2*n))
		sorted := slices.Sorted(slices.Values(values))
		k := 1 + rng.IntN(n)
		got, err := KthSmallest(values, k)
		if err != nil || got != sorted[k-1] {
			t.Fatalf("trial %d: KthSmallest(n=%d, k=%d) = (%d, %v), want %d", trial, n, k, got, err, sorted[k-1])
		}
	}
}

// benchInputs returns 8 fixed-seed random inputs of length n. Benchmarks rotate through them: a
// single repeated input of 1e3 lets the branch predictor learn its whole pattern, which made 1e3
// look 2x cheaper per element than 1e4 and 1e5 (doc/reference_algorithms_2026-11-12/REPORT.md).
func benchInputs(n int) [8][]int {
	rng := rand.New(rand.NewPCG(1, 2))
	var inputs [8][]int
	for i := range inputs {
		inputs[i] = randomInts(rng, n, n)
	}
	return inputs
}

func BenchmarkKthSmallest(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				v, _ := KthSmallest(inputs[i&7], size.n/2)
				benchSink = v
			}
		})
	}
}
