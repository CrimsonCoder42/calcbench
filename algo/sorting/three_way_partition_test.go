package sorting

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// checkPartition verifies the three blocks and that the result is a permutation of values.
func checkPartition(t *testing.T, label string, values, got []int, pivot, less, notGreater int) {
	t.Helper()
	if less < 0 || less > notGreater || notGreater > len(got) {
		t.Fatalf("%s: bounds (%d, %d) for n = %d", label, less, notGreater, len(got))
	}
	for i, v := range got {
		c := cmp.Compare(v, pivot)
		if (i < less && c >= 0) || (i >= less && i < notGreater && c != 0) || (i >= notGreater && c <= 0) {
			t.Fatalf("%s: %v is not partitioned at %d around %d (%d, %d)", label, got, i, pivot, less, notGreater)
		}
	}
	if !slices.Equal(slices.Sorted(slices.Values(got)), slices.Sorted(slices.Values(values))) {
		t.Fatalf("%s: result is not a permutation of the input", label)
	}
}

// Dijkstra, A Discipline of Programming, 1976, ch. 14: buckets of three colours. Our own data, colours
// coded 0 (red), 1 (white), 2 (blue), pivot white.
func TestThreeWayPartition_dutchFlagExample(t *testing.T) {
	values := []int{2, 0, 2, 1, 1, 0}
	got, less, notGreater := ThreeWayPartition(nil, values, 1)
	if !slices.Equal(got, []int{0, 0, 1, 1, 2, 2}) || less != 2 || notGreater != 4 {
		t.Errorf("ThreeWayPartition = (%v, %d, %d), want ([0 0 1 1 2 2], 2, 4)", got, less, notGreater)
	}
}

func TestThreeWayPartition_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		pivot  int
	}{
		{"empty", []int{}, 0},
		{"one element equal", []int{5}, 5},
		{"one element less", []int{4}, 5},
		{"two elements", []int{6, 4}, 5},
		{"all equal", []int{3, 3, 3}, 3},
		{"pivot absent", []int{9, 1, 8, 2}, 5},
		{"pivot below all", []int{3, 1, 2}, -1},
		{"pivot above all", []int{3, 1, 2}, 10},
		{"negative", []int{-1, -3, -2, -3}, -3},
		{"extremes", []int{math.MaxInt, math.MinInt, 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, less, notGreater := ThreeWayPartition(nil, tt.values, tt.pivot)
			checkPartition(t, tt.name, tt.values, got, tt.pivot, less, notGreater)
		})
	}
}

func TestThreeWayPartition_inPlaceAndInputKept(t *testing.T) {
	values := []int{2, 0, 1}
	ThreeWayPartition(nil, values, 1)
	if !slices.Equal(values, []int{2, 0, 1}) {
		t.Errorf("ThreeWayPartition modified its input: %v", values)
	}
	got, _, _ := ThreeWayPartition(values, values, 1)
	if !slices.Equal(values, []int{0, 1, 2}) || &got[0] != &values[0] {
		t.Errorf("in place: %v", values)
	}
}

func TestThreeWayPartition_noAllocation(t *testing.T) {
	dst := make([]int, len(knuth16))
	allocs := testing.AllocsPerRun(100, func() { _, benchSink, _ = ThreeWayPartition(dst, knuth16, 512) })
	if allocs != 0 {
		t.Errorf("ThreeWayPartition with a large enough dst allocates %v times per call, want 0", allocs)
	}
}

func TestThreeWayPartition_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 400 {
		values := randomInts(rng, rng.IntN(100), 10)
		pivot := rng.IntN(12) - 1
		got, less, notGreater := ThreeWayPartition(nil, values, pivot)
		checkPartition(t, "oracle", values, got, pivot, less, notGreater)
	}
}

func BenchmarkThreeWayPartition(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, 3)
		dst := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, benchSink, _ = ThreeWayPartition(dst, values, 1)
			}
		})
	}
}
