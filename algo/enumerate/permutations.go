package enumerate

import (
	"iter"
	"slices"
)

// Permutations returns an iterator over all n! orderings of values, one per arrangement of the
// positions, in lexicographic order of the positions (so for strictly ascending values, in lexicographic
// order of the values). Equal values are not merged: [1 1] yields two (equal) orderings; use
// PermutationsUnique for distinct ones. Empty input yields one empty permutation.
// The yielded slice is reused between steps: copy it to keep it, never modify it. values is copied when
// iteration starts and never modified.
//
// ID:      ALG-ENUM-002
// Time:    O(1) amortised per permutation (O(n) worst step), O(n!) in all   n = len(values)
// Space:   O(n) extra: the yielded buffer and the index array, allocated once per iteration
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.2, Algorithm L, run on the index array 0…n−1 with every interchange mirrored on the values
// Tags:    permutation, enumeration, iterator, lexicographic, backtracking
func Permutations[T any](values []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(values)
		buf := slices.Clone(values)
		idx := make([]int, n)
		for i := range idx {
			idx[i] = i
		}
		for {
			if !yield(buf) {
				return
			}
			j := n - 2
			for j >= 0 && idx[j] > idx[j+1] {
				j--
			}
			if j < 0 {
				return
			}
			l := n - 1
			for idx[j] > idx[l] {
				l--
			}
			idx[j], idx[l] = idx[l], idx[j]
			buf[j], buf[l] = buf[l], buf[j]
			slices.Reverse(idx[j+1:])
			slices.Reverse(buf[j+1:])
		}
	}
}
