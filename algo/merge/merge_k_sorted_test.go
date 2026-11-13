package merge

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchLists is the number of lists in every benchmark: fixed, so growth in n is visible on its own.
const benchLists = 16

// Knuth, TAOCP Vol. 3, 2nd ed., ch. 5: the 16 keys, cut in Knuth's order into four runs of four, each
// sorted. Their merge is the sorted 16 keys (read off by hand).
func TestMergeKSorted_knuthExample(t *testing.T) {
	runs := [][]int{{61, 87, 503, 512}, {170, 275, 897, 908}, {154, 426, 509, 653}, {612, 677, 703, 765}}
	want := []int{61, 87, 154, 170, 275, 426, 503, 509, 512, 612, 653, 677, 703, 765, 897, 908}
	if got := MergeKSorted(runs); !slices.Equal(got, want) {
		t.Errorf("MergeKSorted(knuth runs) = %v, want %v", got, want)
	}
}

func TestMergeKSorted_cases(t *testing.T) {
	tests := []struct {
		name  string
		lists [][]int
		want  []int
	}{
		{"no lists", nil, []int{}},
		{"only empty lists", [][]int{nil, {}, nil}, []int{}},
		{"one list", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
		{"one non-empty among empties", [][]int{nil, {4, 5}, {}}, []int{4, 5}},
		{"two lists", [][]int{{1, 4, 9}, {2, 3, 10, 11}}, []int{1, 2, 3, 4, 9, 10, 11}},
		{"two lists among empties", [][]int{{}, {5}, nil, {1}}, []int{1, 5}},
		{"three lists, one element each", [][]int{{3}, {1}, {2}}, []int{1, 2, 3}},
		{"duplicates across lists", [][]int{{1, 2, 2}, {2, 3}, {0, 2}}, []int{0, 1, 2, 2, 2, 2, 3}},
		{"all equal", [][]int{{7, 7}, {7}, {7, 7, 7}}, []int{7, 7, 7, 7, 7, 7}},
		{"disjoint ranges, reversed list order", [][]int{{7, 8}, {4, 5}, {1, 2}}, []int{1, 2, 4, 5, 7, 8}},
		{"negative and zero", [][]int{{-5, 0}, {-3}, {-9, -1}}, []int{-9, -5, -3, -1, 0}},
		{"extremes of the type", [][]int{{math.MinInt, math.MaxInt}, {0}, {math.MinInt}}, []int{math.MinInt, math.MinInt, 0, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeKSorted(tt.lists); got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("MergeKSorted(%v) = %#v, want %v", tt.lists, got, tt.want)
			}
		})
	}
}

// -0 and +0 compare equal, so their sign shows which list each output element came from. Stability
// means the lower list index comes first, on both the two-list path and the heap path.
func TestMergeKSorted_stable(t *testing.T) {
	negZero := math.Copysign(0, -1)
	tests := []struct {
		name      string
		lists     [][]float64
		wantSigns []bool // math.Signbit of each output element
	}{
		{"two lists, -0 first", [][]float64{{negZero}, {0}}, []bool{true, false}},
		{"two lists, +0 first", [][]float64{{0}, {negZero}}, []bool{false, true}},
		{"three lists", [][]float64{{-1, 0}, {negZero, negZero}, {0, 1}}, []bool{true, false, true, true, false, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeKSorted(tt.lists)
			if len(got) != len(tt.wantSigns) {
				t.Fatalf("MergeKSorted(%v) = %v, want %d elements", tt.lists, got, len(tt.wantSigns))
			}
			for i, v := range got {
				if math.Signbit(v) != tt.wantSigns[i] {
					t.Errorf("element %d = %v, sign bit %v, want %v", i, v, math.Signbit(v), tt.wantSigns[i])
				}
			}
		})
	}
}

func TestMergeKSorted_floats(t *testing.T) {
	got := MergeKSorted([][]float64{{math.NaN(), 1}, {math.Inf(-1), math.Inf(1)}, {math.NaN(), 0}})
	// cmp.Compare order: NaN first, then -Inf, 0, 1, +Inf.
	if len(got) != 6 || !math.IsNaN(got[0]) || !math.IsNaN(got[1]) || !math.IsInf(got[2], -1) ||
		got[3] != 0 || got[4] != 1 || !math.IsInf(got[5], 1) {
		t.Errorf("MergeKSorted(floats) = %v, want [NaN NaN -Inf 0 1 +Inf]", got)
	}
}

// Unsorted input breaks the order (a precondition), but no element is lost or repeated.
func TestMergeKSorted_unsortedKeepsEveryElement(t *testing.T) {
	lists := [][]int{{3, 1}, {9, 2, 8}, {5}}
	got := MergeKSorted(lists)
	slices.Sort(got)
	if want := []int{1, 2, 3, 5, 8, 9}; !slices.Equal(got, want) {
		t.Errorf("sorted(MergeKSorted(unsorted)) = %v, want %v", got, want)
	}
}

// Oracle: concatenate and sort. Lists are cloned before and compared after, so the test also proves
// the input is not modified.
func TestMergeKSorted_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		lists := make([][]int, rng.IntN(12))
		var all []int
		for i := range lists {
			lists[i] = make([]int, rng.IntN(20))
			for j := range lists[i] {
				lists[i][j] = rng.IntN(30)
			}
			slices.Sort(lists[i])
			all = append(all, lists[i]...)
		}
		before := make([][]int, len(lists))
		for i := range lists {
			before[i] = slices.Clone(lists[i])
		}
		slices.Sort(all)
		if got := MergeKSorted(lists); !slices.Equal(got, all) {
			t.Fatalf("trial %d: MergeKSorted(%v) = %v, want %v", trial, lists, got, all)
		}
		for i := range lists {
			if !slices.Equal(lists[i], before[i]) {
				t.Fatalf("trial %d: MergeKSorted modified list %d", trial, i)
			}
		}
	}
}

func BenchmarkMergeKSorted(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		lists := make([][]int, benchLists)
		for i := range lists {
			lists[i] = make([]int, size.n/benchLists)
			for j := range lists[i] {
				lists[i][j] = rng.IntN(size.n)
			}
			slices.Sort(lists[i])
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = MergeKSorted(lists)
			}
		})
	}
}
