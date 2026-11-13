package search

import "cmp"

// SearchRotatedDup returns an index of target in rotated, an ascending slice that may hold duplicates and
// has been rotated by an unknown amount, such as [2 2 2 3 4 2 2]; found is false and index is -1 when
// target is absent. When target equals rotated[0] the index is 0. rotated is read, never modified.
// Duplicates can hide the rotation point: in [1 1 … 1 0 1 … 1] the 0 can sit anywhere, so the worst
// case reads every element. It is O(log n) whenever rotated[0] differs from the last element.
//
// ID:      ALG-SEARCH-006
// Time:    O(n) worst case, O(log n) when rotated[0] ≠ rotated[n-1] (more exactly O(t + log n), t = length of the trailing run equal to rotated[0])   n = len(rotated)
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search), applied twice
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (first position where a monotone test holds)
// Tags:    search, binary-search, rotated, cyclic-shift, duplicates
func SearchRotatedDup[T cmp.Ordered](rotated []T, target T) (index int, found bool) {
	n := len(rotated)
	if n == 0 {
		return -1, false
	}
	if cmp.Compare(rotated[0], target) == 0 {
		return 0, true
	}
	// Elements equal to rotated[0] at the end belong to the second part and break the monotone test
	// rotated[i] < rotated[0]; after they are dropped it holds again, and they cannot be target.
	end := n
	for end > 1 && cmp.Compare(rotated[end-1], rotated[0]) == 0 {
		end--
	}
	return searchRotatedParts(rotated, end, target)
}
