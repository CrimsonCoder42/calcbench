package enumerate

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// binomial is C(n, k) for the small n of these tests.
func binomial(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	c := 1
	for i := 1; i <= k; i++ {
		c = c * (n - k + i) / i
	}
	return c
}

// Kreher & Stinson, Combinatorial Algorithms, §2.3.1: the 2-subsets of {1,2,3,4} in lexicographic order
// are 12 13 14 23 24 34.
func TestCombinations_lexExample(t *testing.T) {
	got := collect(Combinations([]int{1, 2, 3, 4}, 2))
	want := [][]int{{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCombinations_cases(t *testing.T) {
	five := []int{10, 20, 30, 40, 50}
	tests := []struct {
		name string
		in   []int
		k    int
		want int
	}{
		{"empty, k=0", nil, 0, 1},
		{"empty, k=1", nil, 1, 0},
		{"k negative", five, -1, 0},
		{"k above n", five, 6, 0},
		{"k=0", five, 0, 1},
		{"k=n", five, 5, 1},
		{"k=1", five, 1, 5},
		{"k=3", five, 3, 10},
		{"duplicates not merged", []int{1, 1, 1}, 2, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(collect(Combinations(tt.in, tt.k))); got != tt.want {
				t.Errorf("Combinations(%v, %d) yields %d, want %d", tt.in, tt.k, got, tt.want)
			}
		})
	}
}

// Choosing 2 of 3 people for a shift (our own example): Ann+Bo, Ann+Cy, Bo+Cy.
func TestCombinations_strings(t *testing.T) {
	got := collect(Combinations([]string{"Ann", "Bo", "Cy"}, 2))
	want := [][]string{{"Ann", "Bo"}, {"Ann", "Cy"}, {"Bo", "Cy"}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCombinations_earlyBreak(t *testing.T) {
	count := 0
	for range Combinations(make([]int, 60), 30) {
		count++
		if count == 4 {
			break
		}
	}
	if count != 4 {
		t.Errorf("count = %d, want 4", count)
	}
}

// TestCombinations_oracle: the slow version filters all n-bit masks with k bits set, as position lists,
// then sorts them lexicographically.
func TestCombinations_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(10)
		k := rng.IntN(n+2) - 1
		pos := make([]int, n)
		for i := range pos {
			pos[i] = i
		}
		var want [][]int
		for mask := 0; mask < 1<<n; mask++ {
			var c []int
			for i := range n {
				if mask>>i&1 == 1 {
					c = append(c, i)
				}
			}
			if len(c) == k {
				want = append(want, c)
			}
		}
		slices.SortFunc(want, slices.Compare)
		got := collect(Combinations(pos, k))
		if len(got) != binomial(n, k) || len(got) != len(want) {
			t.Fatalf("trial %d: n=%d k=%d: %d combinations, want %d", trial, n, k, len(got), binomial(n, k))
		}
		for i := range got {
			if !slices.Equal(got[i], want[i]) && !(len(got[i]) == 0 && len(want[i]) == 0) {
				t.Fatalf("trial %d: n=%d k=%d: [%d] = %v, want %v", trial, n, k, i, got[i], want[i])
			}
		}
	}
}

func BenchmarkCombinations(b *testing.B) {
	for _, n := range []int{10, 16, 20} {
		vals := make([]int, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				for c := range Combinations(vals, n/2) {
					sink += len(c)
				}
			}
			benchSink = sink
		})
	}
}
