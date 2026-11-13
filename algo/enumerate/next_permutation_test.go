package enumerate

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// benchSink keeps benchmark results observable. Shared by every benchmark in this package.
var benchSink int

// factorial is n! for the small n of these tests.
func factorial(n int) int {
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}
	return f
}

// Knuth, TAOCP Vol. 4A, §7.2.1.2, Algorithm L: from a₁ ≤ a₂ ≤ a₃ = 1 2 3 the algorithm visits the six
// permutations in lexicographic order and then terminates (hand-run of steps L2–L4).
func TestNextPermutationInPlace_knuthExample(t *testing.T) {
	want := [][]int{{1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}
	p := []int{1, 2, 3}
	for i, w := range want {
		if !NextPermutationInPlace(p) || !slices.Equal(p, w) {
			t.Fatalf("step %d: got %v, want %v and true", i, p, w)
		}
	}
	if NextPermutationInPlace(p) || !slices.Equal(p, []int{1, 2, 3}) {
		t.Fatalf("after the last: got %v, want false and [1 2 3]", p)
	}
}

func TestNextPermutationInPlace_cases(t *testing.T) {
	tests := []struct {
		name   string
		in     []int
		want   []int
		wantOK bool
	}{
		{"empty", nil, nil, false},
		{"one element", []int{7}, []int{7}, false},
		{"two ascending", []int{1, 2}, []int{2, 1}, true},
		{"two descending wraps", []int{2, 1}, []int{1, 2}, false},
		{"all equal", []int{5, 5, 5}, []int{5, 5, 5}, false},
		{"duplicates", []int{1, 1, 2}, []int{1, 2, 1}, true},
		{"duplicates last wraps", []int{2, 1, 1}, []int{1, 1, 2}, false},
		{"middle pivot", []int{1, 3, 5, 4, 2}, []int{1, 4, 2, 3, 5}, true},
		{"negative values", []int{-3, -1, -2}, []int{-2, -3, -1}, true},
		{"extremes", []int{math.MinInt, math.MaxInt}, []int{math.MaxInt, math.MinInt}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := slices.Clone(tt.in)
			ok := NextPermutationInPlace(p)
			if ok != tt.wantOK || !slices.Equal(p, tt.want) {
				t.Errorf("NextPermutationInPlace(%v) = %v, %v; want %v, %v", tt.in, p, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestNextPermutationInPlace_floatsAndStrings(t *testing.T) {
	f := []float64{math.NaN(), 1, math.Inf(1)}
	if !NextPermutationInPlace(f) || !math.IsNaN(f[0]) || f[1] != math.Inf(1) || f[2] != 1 {
		t.Errorf("floats: got %v, want [NaN +Inf 1]", f)
	}
	s := []string{"a", "c", "b"}
	if !NextPermutationInPlace(s) || !slices.Equal(s, []string{"b", "a", "c"}) {
		t.Errorf("strings: got %v, want [b a c]", s)
	}
}

func TestNextPermutationInPlace_noAllocation(t *testing.T) {
	p := []int{1, 2, 3, 4, 5, 6}
	if allocs := testing.AllocsPerRun(100, func() { NextPermutationInPlace(p) }); allocs != 0 {
		t.Errorf("NextPermutationInPlace allocates %v times per call, want 0", allocs)
	}
}

// TestNextPermutationInPlace_oracle runs full cycles on random multisets: every arrangement appears
// once, strictly increasing, and the count is n!/∏mᵢ! (counted by a set of keys).
func TestNextPermutationInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(7)
		p := make([]int, n)
		for i := range p {
			p[i] = rng.IntN(4)
		}
		slices.Sort(p)
		seen := map[string]bool{}
		prev := slices.Clone(p)
		seen[fmt.Sprint(p)] = true
		for NextPermutationInPlace(p) {
			if slices.Compare(prev, p) >= 0 {
				t.Fatalf("trial %d: %v does not follow %v", trial, p, prev)
			}
			key := fmt.Sprint(p)
			if seen[key] {
				t.Fatalf("trial %d: %v repeated", trial, p)
			}
			seen[key] = true
			copy(prev, p)
		}
		want := factorial(n)
		counts := map[int]int{}
		for _, v := range p {
			counts[v]++
		}
		for _, c := range counts {
			want /= factorial(c)
		}
		if len(seen) != want {
			t.Fatalf("trial %d: %d permutations, want %d", trial, len(seen), want)
		}
	}
}

func BenchmarkNextPermutationInPlace(b *testing.B) {
	for _, n := range []int{8, 10, 12} {
		p := make([]int, n)
		for i := range p {
			p[i] = i
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				if NextPermutationInPlace(p) {
					sink++
				}
			}
			benchSink = sink
		})
	}
}
