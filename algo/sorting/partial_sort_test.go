package sorting

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., ch. 5 keys: the 5 smallest, ascending, read off the sorted list.
func TestPartialSort_knuthExample(t *testing.T) {
	got, err := PartialSort(nil, knuth16, 5)
	if err != nil || !slices.Equal(got, []int{61, 87, 154, 170, 275}) {
		t.Errorf("PartialSort(knuth16, 5) = (%v, %v), want [61 87 154 170 275]", got, err)
	}
}

func TestPartialSort_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		k      int
		want   []int
	}{
		{"empty", []int{}, 3, []int{}},
		{"k zero", []int{3, 1}, 0, []int{}},
		{"one element", []int{4}, 1, []int{4}},
		{"k above n is clamped", []int{3, 1, 2}, 10, []int{1, 2, 3}},
		{"duplicates count separately", []int{5, 1, 5, 1, 5}, 3, []int{1, 1, 5}},
		{"reversed", []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -1, -2, -3, -4, -5, -6, -7, -8, -9, -10}, 4, []int{-10, -9, -8, -7}},
		{"extremes", []int{math.MaxInt, 0, math.MinInt}, 2, []int{math.MinInt, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PartialSort(nil, tt.values, tt.k)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("PartialSort(%v, %d) = (%v, %v), want %v", tt.values, tt.k, got, err, tt.want)
			}
		})
	}
}

func TestPartialSort_errors(t *testing.T) {
	for _, k := range []int{-1, math.MinInt} {
		dst := []int{9}
		got, err := PartialSort(dst, knuth16, k)
		if !errors.Is(err, ErrNegativeK) || !slices.Equal(got, []int{9}) {
			t.Errorf("PartialSort(k=%d) = (%v, %v), want (dst, ErrNegativeK)", k, got, err)
		}
	}
}

func TestPartialSort_inPlaceAndInputKept(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, err := PartialSort(nil, values, 3); err != nil || !slices.Equal(values, knuth16) {
		t.Errorf("PartialSort modified its input: %v", values)
	}
	got, err := PartialSort(values, values, 3)
	if err != nil || !slices.Equal(got, []int{61, 87, 154}) || &got[0] != &values[0] {
		t.Errorf("in place: %v", got)
	}
}

func TestPartialSort_noAllocation(t *testing.T) {
	dst := make([]int, len(knuth16))
	allocs := testing.AllocsPerRun(100, func() {
		out, _ := PartialSort(dst, knuth16, 5)
		benchSink = out[0]
	})
	if allocs != 0 {
		t.Errorf("PartialSort with a large enough dst allocates %v times per call, want 0", allocs)
	}
}

func TestPartialSort_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(300)
		values := randomInts(rng, n, 1+rng.IntN(2*n+1))
		k := rng.IntN(n + 2)
		want := slices.Sorted(slices.Values(values))[:min(k, n)]
		got, err := PartialSort(nil, values, k)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: PartialSort(n=%d, k=%d) differs from sort-and-take", trial, n, k)
		}
	}
}

// k = 100 fixed: the O(n) selection dominates the O(k log k) sort.
func BenchmarkPartialSort(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n)
		dst := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				out, _ := PartialSort(dst, values, 100)
				benchSink = out[0]
			}
		})
	}
}
