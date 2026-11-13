package search

import "cmp"

// LowerBound returns the first index i with sorted[i] >= target, or len(sorted) when every element is
// less than target: the insertion point of target before any elements equal to it.
// sorted must be ascending in cmp.Compare order; it is read, never modified.
//
// ID:      ALG-SEARCH-001
// Time:    O(log n)   n = len(sorted); at most ⌈log2(n+1)⌉ comparisons
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (binary search for the first occurrence)
// Tags:    search, binary-search, sorted, lower-bound, insertion-point
func LowerBound[T cmp.Ordered](sorted []T, target T) int {
	// Invariant: sorted[:lo] < target and sorted[hi:] >= target.
	lo, hi := 0, len(sorted)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if cmp.Less(sorted[mid], target) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}
