package stats

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// medianInPlace and selectFloatRange against a sorted copy, on random inputs with many ties and on
// sorted, reversed and constant runs, with and without forced median-of-medians pivots.
func TestMedianInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	shapes := []func(n int) []float64{
		func(n int) []float64 { return randomFloats(rng, n) },
		func(n int) []float64 { // many ties
			v := make([]float64, n)
			for i := range v {
				v[i] = float64(rng.IntN(5))
			}
			return v
		},
		func(n int) []float64 { // ascending
			v := make([]float64, n)
			for i := range v {
				v[i] = float64(i)
			}
			return v
		},
		func(n int) []float64 { // descending
			v := make([]float64, n)
			for i := range v {
				v[i] = float64(n - i)
			}
			return v
		},
	}
	for trial := range 400 {
		n := 1 + rng.IntN(400)
		values := shapes[trial%len(shapes)](n)
		sorted := slices.Sorted(slices.Values(values))
		want := midpoint(sorted[(n-1)/2], sorted[n/2])
		if got := medianInPlace(slices.Clone(values)); got != want {
			t.Fatalf("trial %d: medianInPlace = %v; want %v", trial, got, want)
		}
		k := rng.IntN(n)
		for _, force := range []bool{false, true} {
			a := slices.Clone(values)
			selectFloatRange(a, k, force)
			if a[k] != sorted[k] {
				t.Fatalf("trial %d (forceMoM %v): a[%d] = %v; want %v", trial, force, k, a[k], sorted[k])
			}
			for i := range a {
				if (i < k && a[i] > a[k]) || (i > k && a[i] < a[k]) {
					t.Fatalf("trial %d (forceMoM %v): a[%d] = %v on the wrong side of a[%d] = %v", trial, force, i, a[i], k, a[k])
				}
			}
		}
	}
}
