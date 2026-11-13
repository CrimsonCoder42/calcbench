package enumerate

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 4A, §7.2.1.2: Algorithm L on the multiset {1, 1, 2} visits 112, 121, 211.
func TestPermutationsUnique_knuthExample(t *testing.T) {
	got := collect(PermutationsUnique([]int{2, 1, 1}))
	want := [][]int{{1, 1, 2}, {1, 2, 1}, {2, 1, 1}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPermutationsUnique_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"empty", nil, 1},
		{"one", []int{3}, 1},
		{"all equal", []int{7, 7, 7, 7}, 1},
		{"distinct", []int{3, 1, 2}, 6},
		{"two pairs", []int{1, 2, 1, 2}, 6},
		{"negative", []int{-1, -1, 0}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.in)
			if got := len(collect(PermutationsUnique(in))); got != tt.want {
				t.Errorf("PermutationsUnique(%v) yields %d, want %d", tt.in, got, tt.want)
			}
			if !slices.Equal(in, tt.in) {
				t.Errorf("input modified: %v", in)
			}
		})
	}
}

func TestPermutationsUnique_earlyBreak(t *testing.T) {
	count := 0
	for range PermutationsUnique([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}) {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

// TestPermutationsUnique_oracle: the slow version deduplicates all n! orderings with a set and sorts.
func TestPermutationsUnique_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		n := rng.IntN(7)
		vals := make([]int, n)
		for i := range vals {
			vals[i] = rng.IntN(3)
		}
		seen := map[string][]int{}
		for p := range Permutations(vals) {
			seen[fmt.Sprint(p)] = slices.Clone(p)
		}
		var want [][]int
		for _, p := range seen {
			want = append(want, p)
		}
		slices.SortFunc(want, slices.Compare)
		got := collect(PermutationsUnique(vals))
		if !slices.EqualFunc(got, want, slices.Equal) {
			t.Fatalf("trial %d: PermutationsUnique(%v) = %v, want %v", trial, vals, got, want)
		}
	}
}

func BenchmarkPermutationsUnique(b *testing.B) {
	for _, n := range []int{6, 8, 10} {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i / 2 // pairs: n!/2^(n/2) permutations
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				for p := range PermutationsUnique(vals) {
					sink += len(p)
				}
			}
			benchSink = sink
		})
	}
}
