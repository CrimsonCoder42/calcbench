package selection

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// checkSelected reports whether a[k] == want and a is partitioned around index k.
func checkSelected(t *testing.T, label string, a []int, k, want int) {
	t.Helper()
	if a[k] != want {
		t.Fatalf("%s: a[%d] = %d, want %d", label, k, a[k], want)
	}
	for i, v := range a {
		if (i < k && v > a[k]) || (i > k && v < a[k]) {
			t.Fatalf("%s: a[%d] = %d is on the wrong side of a[%d] = %d", label, i, v, k, a[k])
		}
	}
}

// The median-of-medians path is reached only on adversarial inputs, so it is forced here and checked
// against a sort on the same generated inputs as the public oracles.
func TestSelectRange_medianOfMediansOracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := 1 + rng.IntN(500)
		values := randomInts(rng, n, 1+rng.IntN(2*n))
		sorted := slices.Sorted(slices.Values(values))
		k := rng.IntN(n)
		for _, force := range []bool{false, true} {
			a := slices.Clone(values)
			selectRange(a, k, force)
			checkSelected(t, fmt.Sprintf("trial %d, forceMoM %v", trial, force), a, k, sorted[k])
		}
	}
}

// Inputs that defeat naive pivots: sorted, reversed, organ pipe, sawtooth, all equal. Each must give the
// right element with the partition property, on both paths.
func TestSelectRange_patterns(t *testing.T) {
	const n = 1000
	patterns := map[string]func(i int) int{
		"sorted":     func(i int) int { return i },
		"reversed":   func(i int) int { return n - i },
		"organ pipe": func(i int) int { return min(i, n-1-i) },
		"sawtooth":   func(i int) int { return i % 17 },
		"all equal":  func(int) int { return 5 },
	}
	for name, gen := range patterns {
		values := make([]int, n)
		for i := range values {
			values[i] = gen(i)
		}
		sorted := slices.Sorted(slices.Values(values))
		for _, k := range []int{0, 1, n / 3, n / 2, n - 2, n - 1} {
			for _, force := range []bool{false, true} {
				a := slices.Clone(values)
				selectRange(a, k, force)
				checkSelected(t, name, a, k, sorted[k])
			}
		}
	}
}

// Blum et al. (1973): at least 3/10 of the elements are <= the median of medians and at least 3/10 are
// >= it, less a constant for the partial last group (CLRS 3rd ed. §9.3 counts 3n/10 - 6).
func TestMedianOfMedians_guarantee(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		n := 5 + rng.IntN(1000)
		values := randomInts(rng, n, n)
		pivot := medianOfMedians(slices.Clone(values))
		notAbove, notBelow := 0, 0
		for _, v := range values {
			if v <= pivot {
				notAbove++
			}
			if v >= pivot {
				notBelow++
			}
		}
		if bound := 3*n/10 - 6; notAbove < bound || notBelow < bound {
			t.Fatalf("trial %d: n=%d, %d <= pivot, %d >= pivot, want both >= %d", trial, n, notAbove, notBelow, bound)
		}
	}
}

func TestPartitionThreeWay_blocks(t *testing.T) {
	s := []int{5, 1, 5, 9, 3, 5, 7, 5, 0}
	less, notGreater := partitionThreeWay(s, 5)
	if less != 3 || notGreater != 7 {
		t.Fatalf("partitionThreeWay boundaries = (%d, %d), want (3, 7): %v", less, notGreater, s)
	}
	for i, v := range s {
		if (i < less && v >= 5) || (i >= less && i < notGreater && v != 5) || (i >= notGreater && v <= 5) {
			t.Fatalf("partitionThreeWay: s[%d] = %d in the wrong block: %v", i, v, s)
		}
	}
}

func TestMedianOfThree_allOrders(t *testing.T) {
	for _, xyz := range [][3]int{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}, {2, 2, 1}, {1, 2, 2}, {2, 1, 2}} {
		want := slices.Sorted(slices.Values(xyz[:]))[1]
		if got := medianOfThree(xyz[0], xyz[1], xyz[2]); got != want {
			t.Errorf("medianOfThree%v = %d, want %d", xyz, got, want)
		}
	}
}

// Measures the constant-factor cost of the fallback path, the evidence for D-001.
func BenchmarkSelectRangeMedianOfMedians(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		work := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				copy(work, inputs[i&7])
				selectRange(work, size.n/2, true)
			}
		})
	}
}
