package sorting

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Sedgewick & Wayne §2.5 "Duplicates": sort, then one pass noting equal neighbours. Our own example.
func TestDedupeSorted_example(t *testing.T) {
	if got := DedupeSorted(nil, []int{1, 1, 2, 3, 3, 3, 4}); !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Errorf("DedupeSorted = %v, want [1 2 3 4]", got)
	}
}

// Our own data: a sorted list of scores with repeats keeps one of each.
func TestDedupeSorted_exercises(t *testing.T) {
	if got := DedupeSorted(nil, []int{2, 2, 5, 5, 5, 8, 9, 9}); !slices.Equal(got, []int{2, 5, 8, 9}) {
		t.Errorf("got %v, want [2 5 8 9]", got)
	}
}

func TestDedupeSorted_cases(t *testing.T) {
	tests := []struct {
		name   string
		sorted []int
		want   []int
	}{
		{"empty", []int{}, []int{}},
		{"one element", []int{5}, []int{5}},
		{"two equal", []int{5, 5}, []int{5}},
		{"two distinct", []int{4, 5}, []int{4, 5}},
		{"all equal", []int{7, 7, 7, 7}, []int{7}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"negative", []int{-2, -2, -1, 0, 0}, []int{-2, -1, 0}},
		{"extremes", []int{math.MinInt, math.MinInt, math.MaxInt, math.MaxInt}, []int{math.MinInt, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DedupeSorted(nil, tt.sorted); !slices.Equal(got, tt.want) {
				t.Errorf("DedupeSorted(%v) = %v, want %v", tt.sorted, got, tt.want)
			}
		})
	}
}

func TestDedupeSorted_floats(t *testing.T) {
	got := DedupeSorted(nil, []float64{math.NaN(), math.NaN(), math.Copysign(0, -1), 0, 1})
	if len(got) != 3 || !math.IsNaN(got[0]) || !math.Signbit(got[1]) || got[2] != 1 {
		t.Errorf("DedupeSorted(floats) = %v, want [NaN -0 1]", got)
	}
}

func TestDedupeSorted_inPlaceAndInputKept(t *testing.T) {
	values := []int{1, 1, 2, 2, 2, 3}
	if got := DedupeSorted(nil, values); cap(got) != 3 || !slices.Equal(values, []int{1, 1, 2, 2, 2, 3}) {
		t.Errorf("DedupeSorted modified its input or over-allocated: cap %d, input %v", cap(got), values)
	}
	got := DedupeSorted(values, values)
	if !slices.Equal(got, []int{1, 2, 3}) || &got[0] != &values[0] {
		t.Errorf("in place: %v", got)
	}
}

func TestDedupeSorted_noAllocation(t *testing.T) {
	dst := make([]int, len(knuth16Sorted))
	if allocs := testing.AllocsPerRun(100, func() { benchSink = len(DedupeSorted(dst, knuth16Sorted)) }); allocs != 0 {
		t.Errorf("DedupeSorted with a large enough dst allocates %v times per call, want 0", allocs)
	}
}

func TestDedupeSorted_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		sorted := slices.Sorted(slices.Values(randomInts(rng, rng.IntN(100), 20)))
		want := slices.Compact(slices.Clone(sorted))
		if got := DedupeSorted(nil, sorted); !slices.Equal(got, want) {
			t.Fatalf("trial %d: DedupeSorted(%v) = %v, want %v", trial, sorted, got, want)
		}
		if got := DedupeSorted(sorted, sorted); !slices.Equal(got, want) {
			t.Fatalf("trial %d: in place = %v, want %v", trial, got, want)
		}
	}
}

func BenchmarkDedupeSorted(b *testing.B) {
	for _, size := range benchSizes {
		sorted := slices.Sorted(slices.Values(randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n/4)))
		dst := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = len(DedupeSorted(dst, sorted))
			}
		})
	}
}
