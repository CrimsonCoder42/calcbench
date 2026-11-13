package search

import "cmp"

// BinarySearchFirst returns the index of the first element of sorted equal to target and true, or,
// when there is none, the index where target would be inserted and false.
// Equality is cmp.Compare(x, target) == 0, so NaN finds NaN. sorted must be ascending in cmp.Compare
// order; it is read, never modified.
//
// ID:      ALG-SEARCH-003
// Time:    O(log n)   n = len(sorted)
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (binary search for the first occurrence)
// Tags:    search, binary-search, sorted, first-occurrence
func BinarySearchFirst[T cmp.Ordered](sorted []T, target T) (index int, found bool) {
	i := LowerBound(sorted, target)
	return i, i < len(sorted) && cmp.Compare(sorted[i], target) == 0
}
