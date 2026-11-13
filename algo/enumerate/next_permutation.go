package enumerate

import (
	"cmp"
	"slices"
)

// NextPermutationInPlace rearranges values into the next permutation in lexicographic order and returns
// true, or, when values is already the last permutation (non-increasing), reverses it into the first
// (non-decreasing) and returns false. Repeated values are handled: each distinct arrangement of a
// multiset is visited once, so starting from sorted input and looping while it returns true visits every
// distinct permutation exactly once. values is modified in place.
//
// ID:      ALG-ENUM-001
// Time:    O(n) worst case per call, O(1) amortised over the n!/∏mᵢ! calls of a full run   n = len(values), mᵢ = multiplicities
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.2, Algorithm L (lexicographic permutation generation, steps L2–L4)
// Tags:    permutation, lexicographic, next-permutation, multiset, in-place
func NextPermutationInPlace[T cmp.Ordered](values []T) bool {
	n := len(values)
	// L2: the largest j with values[j] < values[j+1]; everything after j is non-increasing.
	j := n - 2
	for j >= 0 && !cmp.Less(values[j], values[j+1]) {
		j--
	}
	if j < 0 {
		slices.Reverse(values)
		return false
	}
	// L3: the largest l with values[j] < values[l]; swap them.
	l := n - 1
	for !cmp.Less(values[j], values[l]) {
		l--
	}
	values[j], values[l] = values[l], values[j]
	// L4: the suffix is still non-increasing; reverse it to make it the smallest.
	slices.Reverse(values[j+1:])
	return true
}
