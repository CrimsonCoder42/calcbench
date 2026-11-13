package heap

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// knuth16 is the 16-key running example of Knuth, TAOCP Vol. 3, 2nd ed. (ch. 5), in Knuth's order.
var knuth16 = []int{503, 87, 512, 61, 908, 170, 897, 275, 653, 426, 154, 509, 612, 677, 765, 703}

var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchK is the k of every benchmark: fixed, so growth in n is visible on its own.
const benchK = 100

// TopM client of Sedgewick & Wayne, Algorithms, 4th ed., §2.4 (the M largest, in descending order),
// applied to Knuth's 16 keys; the three largest are read off the sorted keys by hand.
func TestTopK_knuthExample(t *testing.T) {
	got, err := TopK(knuth16, 3)
	if want := []int{908, 897, 765}; err != nil || !slices.Equal(got, want) {
		t.Errorf("TopK(knuth16, 3) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestTopK_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		k      int
		want   []int
	}{
		{"empty, k 0", nil, 0, []int{}},
		{"empty, k 3", nil, 3, []int{}},
		{"k 0", []int{4, 1}, 0, []int{}},
		{"one element", []int{4}, 1, []int{4}},
		{"two elements, k 1", []int{4, 9}, 1, []int{9}},
		{"k above n returns all, descending", []int{2, 9, 4}, 10, []int{9, 4, 2}},
		{"k equals n", []int{2, 9, 4}, 3, []int{9, 4, 2}},
		{"duplicates count separately", []int{5, 1, 5, 3, 5}, 2, []int{5, 5}},
		{"duplicates at the boundary", []int{1, 3, 2, 3, 3}, 4, []int{3, 3, 3, 2}},
		{"already sorted", []int{1, 2, 3, 4, 5, 6}, 3, []int{6, 5, 4}},
		{"reversed", []int{6, 5, 4, 3, 2, 1}, 3, []int{6, 5, 4}},
		{"negative and zero", []int{-3, 0, -7, -1}, 2, []int{0, -1}},
		{"extremes of the type", []int{0, math.MinInt, math.MaxInt}, 3, []int{math.MaxInt, 0, math.MinInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TopK(tt.values, tt.k)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("TopK(%v, %d) = (%#v, %v), want (%v, nil)", tt.values, tt.k, got, err, tt.want)
			}
		})
	}
}

func TestTopK_floats(t *testing.T) {
	got, err := TopK([]float64{1, math.NaN(), math.Inf(-1), math.Inf(1), 2}, 5)
	// cmp.Compare order: NaN is the least value, so it comes last.
	if err != nil || len(got) != 5 || got[0] != math.Inf(1) || got[1] != 2 || got[2] != 1 ||
		!math.IsInf(got[3], -1) || !math.IsNaN(got[4]) {
		t.Errorf("TopK(floats, 5) = (%v, %v), want ([+Inf 2 1 -Inf NaN], nil)", got, err)
	}
}

func TestTopK_strings(t *testing.T) {
	got, err := TopK([]string{"fig", "pear", "apple", "plum"}, 2)
	if want := []string{"plum", "pear"}; err != nil || !slices.Equal(got, want) {
		t.Errorf("TopK(fruit, 2) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestTopK_errors(t *testing.T) {
	for _, k := range []int{-1, math.MinInt} {
		if got, err := TopK(knuth16, k); !errors.Is(err, ErrNegativeK) || got != nil {
			t.Errorf("TopK(knuth16, %d) = (%v, %v), want (nil, ErrNegativeK)", k, got, err)
		}
	}
}

func TestTopK_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, err := TopK(values, 5); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, knuth16) {
		t.Errorf("TopK modified its input: %v", values)
	}
}

// Oracle: sort a copy descending and take the first k.
func TestTopK_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(300)
		values := make([]int, n)
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(2*n+1))
		}
		k := rng.IntN(n + 3)
		desc := slices.Sorted(slices.Values(values))
		slices.Reverse(desc)
		want := desc[:min(k, n)]
		got, err := TopK(values, k)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: TopK(n=%d, k=%d) = (%v, %v), want %v", trial, n, k, got, err, want)
		}
	}
}

// Random order is the usual case: most values lose to the root with one comparison. Ascending order
// is the worst case: every value enters the heap and costs a sift of O(log k).
func BenchmarkTopK(b *testing.B) {
	for _, order := range []string{"", "ascending-"} {
		for _, size := range benchSizes {
			values := make([]int, size.n)
			rng := rand.New(rand.NewPCG(1, 2))
			for i := range values {
				values[i] = rng.IntN(size.n)
			}
			if order != "" {
				slices.Sort(values)
			}
			b.Run(order+size.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					top, _ := TopK(values, benchK)
					_ = top
				}
			})
		}
	}
}
