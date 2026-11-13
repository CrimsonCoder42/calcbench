package sorting

import (
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Our example (Knuth vol. 3 §5.2: sorting brings equal keys together): 3 1 3 2 1 3 has distinct values
// 1 2 3, first seen at 1 3 0, occurring 2 1 3 times; inverse codes 2 0 2 1 0 2.
func TestUniqueAll_example(t *testing.T) {
	u, first, inv, counts := UniqueAll([]int{3, 1, 3, 2, 1, 3})
	if !slices.Equal(u, []int{1, 2, 3}) || !slices.Equal(first, []int{1, 3, 0}) ||
		!slices.Equal(inv, []int{2, 0, 2, 1, 0, 2}) || !slices.Equal(counts, []int{2, 1, 3}) {
		t.Errorf("UniqueAll = %v %v %v %v", u, first, inv, counts)
	}
}

func TestUniqueAll_cases(t *testing.T) {
	u, first, inv, counts := UniqueAll([]float64{})
	if u == nil || first == nil || inv == nil || counts == nil || len(u)+len(first)+len(inv)+len(counts) != 0 {
		t.Errorf("empty: %v %v %v %v", u, first, inv, counts)
	}
	u, first, inv, counts = UniqueAll([]float64{7})
	if !slices.Equal(u, []float64{7}) || !slices.Equal(first, []int{0}) || !slices.Equal(inv, []int{0}) || !slices.Equal(counts, []int{1}) {
		t.Errorf("one: %v %v %v %v", u, first, inv, counts)
	}
	// NaNs merge into one value, sorted first; -0 and +0 merge and the first occurrence (-0) represents them.
	nz := math.Copysign(0, -1)
	u, first, inv, counts = UniqueAll([]float64{1, math.NaN(), nz, 0, math.NaN()})
	if len(u) != 3 || !math.IsNaN(u[0]) || !math.Signbit(u[1]) || u[2] != 1 ||
		!slices.Equal(first, []int{1, 2, 0}) || !slices.Equal(inv, []int{2, 0, 1, 1, 0}) || !slices.Equal(counts, []int{2, 2, 1}) {
		t.Errorf("NaN and zeros: %v %v %v %v", u, first, inv, counts)
	}
}

// Oracle: on 300 random inputs, uniques is the sorted deduplicated input, values[i] == uniques[inverse[i]],
// first is the first occurrence and counts sum to n.
func TestUniqueAll_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := r.IntN(100)
		values := make([]int, n)
		for i := range values {
			values[i] = r.IntN(15) - 7
		}
		u, first, inv, counts := UniqueAll(values)
		want := slices.Compact(slices.Sorted(slices.Values(values)))
		if !slices.Equal(u, want) {
			t.Fatalf("uniques %v, want %v", u, want)
		}
		total := 0
		for j, v := range u {
			if slices.Index(values, v) != first[j] {
				t.Fatalf("first[%d] = %d", j, first[j])
			}
			total += counts[j]
		}
		for i, v := range values {
			if u[inv[i]] != v {
				t.Fatalf("inverse[%d] wrong", i)
			}
		}
		if total != n {
			t.Fatalf("counts sum %d, want %d", total, n)
		}
	}
}

func benchInts(n int) []int {
	r := rand.New(rand.NewPCG(3, 4))
	values := make([]int, n)
	for i := range values {
		values[i] = r.IntN(n/10 + 1)
	}
	return values
}

func BenchmarkUniqueAll(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		values := benchInts(n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _, _ = UniqueAll(values)
			}
		})
	}
}
