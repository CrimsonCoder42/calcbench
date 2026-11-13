package sorting

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., §5.2.1 sorts the ch. 5 keys by straight insertion (Algorithm S);
// expected order read off by hand.
func TestInsertionSort_knuthExample(t *testing.T) {
	if got := InsertionSort(nil, knuth16); !slices.Equal(got, knuth16Sorted) {
		t.Errorf("InsertionSort(knuth16) = %v, want %v", got, knuth16Sorted)
	}
}

func TestInsertionSort_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty", []int{}, []int{}},
		{"one element", []int{5}, []int{5}},
		{"two elements reversed", []int{2, 1}, []int{1, 2}},
		{"duplicates", []int{3, 1, 3, 1, 2}, []int{1, 1, 2, 3, 3}},
		{"already sorted", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{"reversed", []int{4, 3, 2, 1}, []int{1, 2, 3, 4}},
		{"negative and zero", []int{0, -5, 3, -1}, []int{-5, -1, 0, 3}},
		{"extremes of the type", []int{math.MaxInt, 0, math.MinInt}, []int{math.MinInt, 0, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InsertionSort(nil, tt.values); !slices.Equal(got, tt.want) {
				t.Errorf("InsertionSort(%v) = %v, want %v", tt.values, got, tt.want)
			}
		})
	}
}

func TestInsertionSort_floats(t *testing.T) {
	got := InsertionSort(nil, []float64{1, math.Inf(-1), math.NaN(), 0})
	if !math.IsNaN(got[0]) || !slices.Equal(got[1:], []float64{math.Inf(-1), 0, 1}) {
		t.Errorf("InsertionSort(floats) = %v, want [NaN -Inf 0 1]", got)
	}
}

func TestInsertionSort_inPlaceAndInputKept(t *testing.T) {
	values := slices.Clone(knuth16)
	InsertionSort(nil, values)
	if !slices.Equal(values, knuth16) {
		t.Errorf("InsertionSort modified its input: %v", values)
	}
	InsertionSort(values, values)
	if !slices.Equal(values, knuth16Sorted) {
		t.Errorf("in place: %v", values)
	}
}

func TestInsertionSort_noAllocation(t *testing.T) {
	dst := make([]int, len(knuth16))
	allocs := testing.AllocsPerRun(100, func() { benchSink = InsertionSort(dst, knuth16)[0] })
	if allocs != 0 {
		t.Errorf("InsertionSort with a large enough dst allocates %v times per call, want 0", allocs)
	}
}

func TestInsertionSort_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInts(rng, rng.IntN(64), 20)
		want := slices.Sorted(slices.Values(values))
		if got := InsertionSort(nil, values); !slices.Equal(got, want) {
			t.Fatalf("trial %d: InsertionSort(%v) = %v, want %v", trial, values, got, want)
		}
	}
}

// Insertion sort is O(n²) on random input, so the sizes are the small n it is meant for (16, 64)
// plus nearly sorted inputs of 1e3/1e4/1e5 (10 swapped adjacent pairs), where it is O(n).
func BenchmarkInsertionSort(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, n := range []int{16, 64} {
		values := randomInts(rng, n, 1000)
		dst := make([]int, n)
		b.Run("random/"+map[int]string{16: "16", 64: "64"}[n], func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = InsertionSort(dst, values)[0]
			}
		})
	}
	for _, size := range benchSizes {
		values := make([]int, size.n)
		for i := range values {
			values[i] = i
		}
		for range 10 {
			i := rng.IntN(size.n - 1)
			values[i], values[i+1] = values[i+1], values[i]
		}
		dst := make([]int, size.n)
		b.Run("nearlySorted/"+size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = InsertionSort(dst, values)[0]
			}
		})
	}
}
