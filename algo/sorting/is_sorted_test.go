package sorting

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., ch. 5: the 16 keys are unsorted in their original order (87 < 503 at
// index 1) and sorted after sorting.
func TestIsSorted_knuthExample(t *testing.T) {
	if ok, at := IsSorted(knuth16); ok || at != 1 {
		t.Errorf("IsSorted(knuth16) = (%v, %d), want (false, 1)", ok, at)
	}
	if ok, at := IsSorted(knuth16Sorted); !ok || at != 16 {
		t.Errorf("IsSorted(sorted) = (%v, %d), want (true, 16)", ok, at)
	}
}

func TestIsSorted_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		ok     bool
		at     int
	}{
		{"empty", nil, true, 0},
		{"one element", []int{1}, true, 1},
		{"two ascending", []int{1, 2}, true, 2},
		{"two descending", []int{2, 1}, false, 1},
		{"duplicates allowed", []int{1, 1, 2, 2}, true, 4},
		{"descent at the end", []int{1, 2, 3, 0}, false, 3},
		{"negative", []int{-3, -2, -2}, true, 3},
		{"extremes", []int{math.MinInt, math.MaxInt}, true, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok, at := IsSorted(tt.values); ok != tt.ok || at != tt.at {
				t.Errorf("IsSorted(%v) = (%v, %d), want (%v, %d)", tt.values, ok, at, tt.ok, tt.at)
			}
		})
	}
}

func TestIsSorted_floats(t *testing.T) {
	if ok, _ := IsSorted([]float64{math.NaN(), math.Inf(-1), math.Copysign(0, -1), 0, 1}); !ok {
		t.Error("NaN first, -0 equal to 0 should be sorted")
	}
	if ok, at := IsSorted([]float64{1, math.NaN()}); ok || at != 1 {
		t.Errorf("NaN after 1 = (%v, %d), want (false, 1)", ok, at)
	}
}

func TestIsSorted_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { _, benchSink = IsSorted(knuth16Sorted) }); allocs != 0 {
		t.Errorf("IsSorted allocates %v times per call, want 0", allocs)
	}
}

func TestIsSorted_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInts(rng, rng.IntN(10), 4)
		want := len(values)
		for i := 1; i < len(values); i++ {
			if values[i] < values[i-1] {
				want = i
				break
			}
		}
		ok, at := IsSorted(values)
		if at != want || ok != slices.IsSorted(values) {
			t.Fatalf("trial %d: IsSorted(%v) = (%v, %d), want first descent %d", trial, values, ok, at, want)
		}
	}
}

func BenchmarkIsSorted(b *testing.B) {
	for _, size := range benchSizes {
		values := make([]int, size.n)
		for i := range values {
			values[i] = i
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, benchSink = IsSorted(values)
			}
		})
	}
}
