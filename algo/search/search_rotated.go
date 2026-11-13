package search

import "cmp"

// SearchRotated returns the index of target in rotated, an ascending slice of distinct values that has
// been rotated (cyclically shifted) by an unknown amount, such as [40 50 60 10 20 30]; found is false
// and index is -1 when target is absent. It finds the rotation point by binary search, then binary
// searches the one ascending part that can hold target. rotated is read, never modified.
// Distinct values are a precondition (use SearchRotatedDup otherwise); on input that breaks it the
// result is an index in [-1, n) and nothing panics.
//
// ID:      ALG-SEARCH-005
// Time:    O(log n)   n = len(rotated); two binary searches
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search), applied twice
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (first position where a monotone test holds: here rotated[i] < rotated[0])
// Tags:    search, binary-search, rotated, cyclic-shift, sorted
func SearchRotated[T cmp.Ordered](rotated []T, target T) (index int, found bool) {
	if len(rotated) == 0 {
		return -1, false
	}
	return searchRotatedParts(rotated, len(rotated), target)
}

// searchRotatedParts searches rotated[:end], where rotated[0] ≤ every element of the first ascending
// part and every element of the second part is < rotated[0]. The caller guarantees 1 ≤ end ≤ len(rotated).
func searchRotatedParts[T cmp.Ordered](rotated []T, end int, target T) (int, bool) {
	first := rotated[0]
	// The rotation point p: rotated[i] < first is false on [1, p) and true on [p, end).
	p, _ := SearchPredicate(1, end, func(i int) bool { return cmp.Less(rotated[i], first) })
	lo, hi := 0, p
	if cmp.Less(target, first) {
		lo, hi = p, end
	}
	i := lo + LowerBound(rotated[lo:hi], target)
	if i < hi && cmp.Compare(rotated[i], target) == 0 {
		return i, true
	}
	return -1, false
}
