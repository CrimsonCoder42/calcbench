package sorting

import (
	"cmp"
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// knuth16 is the 16-key running example of Knuth, TAOCP Vol. 3, 2nd ed. (ch. 5), in Knuth's original
// order. Shared by every test file in this package.
var knuth16 = []int{503, 87, 512, 61, 908, 170, 897, 275, 653, 426, 154, 509, 612, 677, 765, 703}

// knuth16Sorted is knuth16 ascending.
var knuth16Sorted = []int{61, 87, 154, 170, 275, 426, 503, 509, 512, 612, 653, 677, 703, 765, 897, 908}

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

// record is a key with a tag that shows whether equal keys kept their input order.
type record struct {
	key int
	tag int
}

func byKey(a, b record) int { return cmp.Compare(a.key, b.key) }

// Knuth, TAOCP Vol. 3, 2nd ed., ch. 5 keys sorted by merging (§5.2.4); expected order read off by hand.
func TestMergeSortStable_knuthExample(t *testing.T) {
	got, err := MergeSortStable(nil, knuth16, cmp.Compare[int])
	if err != nil || !slices.Equal(got, knuth16Sorted) {
		t.Errorf("MergeSortStable(knuth16) = (%v, %v), want %v", got, err, knuth16Sorted)
	}
}

// Sedgewick & Wayne §2.5 illustrates stability with records sorted by one field keeping the order of
// another. Our own data: deliveries listed by time, sorted by city; inside each city, time order stays.
func TestMergeSortStable_stability(t *testing.T) {
	cities := []record{{3, 1}, {1, 2}, {2, 3}, {1, 4}, {3, 5}, {2, 6}, {1, 7}, {3, 8}}
	want := []record{{1, 2}, {1, 4}, {1, 7}, {2, 3}, {2, 6}, {3, 1}, {3, 5}, {3, 8}}
	got, err := MergeSortStable(nil, cities, byKey)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("MergeSortStable(cities) = %v, want %v", got, want)
	}
}

func TestMergeSortStable_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty", []int{}, []int{}},
		{"one element", []int{5}, []int{5}},
		{"two elements reversed", []int{2, 1}, []int{1, 2}},
		{"duplicates", []int{3, 1, 3, 1, 2}, []int{1, 1, 2, 3, 3}},
		{"already sorted", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
		{"reversed past one run", []int{14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
		{"negative and zero", []int{0, -5, 3, -1}, []int{-5, -1, 0, 3}},
		{"extremes of the type", []int{math.MaxInt, 0, math.MinInt}, []int{math.MinInt, 0, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeSortStable(nil, tt.values, cmp.Compare[int])
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("MergeSortStable(%v) = (%v, %v), want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMergeSortStable_floats(t *testing.T) {
	values := []float64{1, math.NaN(), math.Inf(1), -2, math.Inf(-1), 0}
	got, err := MergeSortStable(nil, values, cmp.Compare[float64])
	if err != nil || !math.IsNaN(got[0]) || !slices.Equal(got[1:], []float64{math.Inf(-1), -2, 0, 1, math.Inf(1)}) {
		t.Errorf("MergeSortStable(floats) = (%v, %v)", got, err)
	}
}

func TestMergeSortStable_errors(t *testing.T) {
	dst := []int{9}
	got, err := MergeSortStable(dst, knuth16, nil)
	if !errors.Is(err, ErrNilCompare) || !slices.Equal(got, []int{9}) {
		t.Errorf("MergeSortStable(nil compare) = (%v, %v), want (dst, ErrNilCompare)", got, err)
	}
}

func TestMergeSortStable_inPlaceAndInputKept(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, err := MergeSortStable(nil, values, cmp.Compare[int]); err != nil || !slices.Equal(values, knuth16) {
		t.Errorf("MergeSortStable modified its input: %v", values)
	}
	got, err := MergeSortStable(values, values, cmp.Compare[int])
	if err != nil || !slices.Equal(values, knuth16Sorted) || &got[0] != &values[0] {
		t.Errorf("in place: values = %v", values)
	}
}

// The oracle is slices.SortStableFunc. Keys from a small range force long runs of equal keys; tags
// expose any reordering of them. Lengths cross several run widths.
func TestMergeSortStable_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	var dst []record
	for trial := range 400 {
		n := rng.IntN(500)
		values := make([]record, n)
		for i := range values {
			values[i] = record{rng.IntN(1 + n/4), i}
		}
		want := slices.Clone(values)
		slices.SortStableFunc(want, byKey)
		got, err := MergeSortStable(dst, values, byKey)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d (n=%d): MergeSortStable differs from slices.SortStableFunc", trial, n)
		}
		dst = got
	}
}

func BenchmarkMergeSortStable(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n)
		dst := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				out, _ := MergeSortStable(dst, values, cmp.Compare[int])
				benchSink = out[0]
			}
		})
	}
}
