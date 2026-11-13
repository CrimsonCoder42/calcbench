package enumerate

import (
	"cmp"
	"iter"
	"slices"
)

// PermutationsUnique returns an iterator over the distinct permutations of the multiset values, in
// lexicographic order (cmp.Compare). [1 1 2] yields 112, 121, 211. Empty input yields one empty
// permutation. The yielded slice is reused between steps: copy it to keep it, never modify it. values is
// copied when iteration starts and never modified.
//
// ID:      ALG-ENUM-003
// Time:    O(n log n) to sort, then O(1) amortised per permutation; P = n!/∏mᵢ! permutations   n = len(values), mᵢ = multiplicity of each distinct value
// Space:   O(n) extra: one buffer, allocated once per iteration
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.2, Algorithm L (stated for a multiset a₁ ≤ … ≤ aₙ); calls NextPermutationInPlace (ALG-ENUM-001)
// Tags:    permutation, multiset, duplicates, enumeration, iterator, lexicographic
func PermutationsUnique[T cmp.Ordered](values []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		buf := slices.Clone(values)
		slices.Sort(buf)
		for {
			if !yield(buf) {
				return
			}
			if !NextPermutationInPlace(buf) {
				return
			}
		}
	}
}
