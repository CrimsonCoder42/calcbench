package enumerate

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
)

// collect copies every yielded slice, so the reused buffer cannot alias the results.
func collect[T any](seq iter.Seq[[]T]) [][]T {
	var out [][]T
	for v := range seq {
		out = append(out, slices.Clone(v))
	}
	return out
}

// Knuth, TAOCP Vol. 4A, §7.2.1.2, Algorithm L on 1 2 3 (positions 0 1 2): the six permutations in
// lexicographic order, here with letters standing for the positions.
func TestPermutations_knuthExample(t *testing.T) {
	got := collect(Permutations([]string{"a", "b", "c"}))
	want := [][]string{{"a", "b", "c"}, {"a", "c", "b"}, {"b", "a", "c"}, {"b", "c", "a"}, {"c", "a", "b"}, {"c", "b", "a"}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("Permutations(abc) = %v, want %v", got, want)
	}
}

func TestPermutations_cases(t *testing.T) {
	if got := collect(Permutations([]int(nil))); len(got) != 1 || len(got[0]) != 0 {
		t.Errorf("empty: got %v, want one empty permutation", got)
	}
	if got := collect(Permutations([]int{9})); len(got) != 1 || !slices.Equal(got[0], []int{9}) {
		t.Errorf("one: got %v", got)
	}
	// Equal values are not merged: two positions give two orderings.
	if got := collect(Permutations([]int{4, 4})); len(got) != 2 {
		t.Errorf("duplicates: got %d permutations, want 2", len(got))
	}
	// Order follows positions, not values: reversed input yields reversed-first.
	got := collect(Permutations([]int{3, 2, 1}))
	if !slices.Equal(got[0], []int{3, 2, 1}) || !slices.Equal(got[5], []int{1, 2, 3}) {
		t.Errorf("reversed: first %v last %v", got[0], got[5])
	}
}

// T needs no ordering: a struct type works.
func TestPermutations_anyType(t *testing.T) {
	type point struct{ x, y int }
	if got := collect(Permutations([]point{{1, 2}, {3, 4}, {5, 6}})); len(got) != 6 {
		t.Errorf("got %d permutations, want 6", len(got))
	}
}

func TestPermutations_inputNotModified(t *testing.T) {
	in := []int{1, 2, 3, 4}
	for range Permutations(in) {
	}
	if !slices.Equal(in, []int{1, 2, 3, 4}) {
		t.Errorf("input modified: %v", in)
	}
}

func TestPermutations_earlyBreak(t *testing.T) {
	count := 0
	for range Permutations(make([]int, 20)) { // 20! would never finish
		count++
		if count == 5 {
			break
		}
	}
	if count != 5 {
		t.Errorf("count = %d, want 5", count)
	}
}

// TestPermutations_oracle: for distinct values the output equals the Algorithm L run of
// NextPermutationInPlace from sorted input, and has n! elements.
func TestPermutations_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		n := rng.IntN(7)
		vals := rng.Perm(100)[:n]
		slices.Sort(vals)
		got := collect(Permutations(vals))
		var want [][]int
		p := slices.Clone(vals)
		for {
			want = append(want, slices.Clone(p))
			if !NextPermutationInPlace(p) {
				break
			}
		}
		if len(got) != factorial(n) || !slices.EqualFunc(got, want, slices.Equal) {
			t.Fatalf("trial %d: Permutations(%v) disagrees with Algorithm L (%d vs %d)", trial, vals, len(got), len(want))
		}
	}
}

func BenchmarkPermutations(b *testing.B) {
	for _, n := range []int{6, 8, 10} {
		vals := make([]int, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				for p := range Permutations(vals) {
					sink += len(p)
				}
			}
			benchSink = sink
		})
	}
}
