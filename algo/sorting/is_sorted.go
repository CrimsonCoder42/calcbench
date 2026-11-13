package sorting

import "cmp"

// IsSorted reports whether values is ascending in cmp.Compare order (equal neighbours allowed), and
// returns firstDescent, the first index i with values[i] < values[i-1], or len(values) when there is
// none. It stops at the first descent. values is read, never modified.
//
// ID:      ALG-SORT-006
// Time:    O(n) worst case   n = len(values); O(i) when the first descent is at index i
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, ch. 5 introduction (sorted order K_p(1) ≤ K_p(2) ≤ … ≤ K_p(N))
// Tags:    sorting, sorted, check, precondition, descent
func IsSorted[T cmp.Ordered](values []T) (sorted bool, firstDescent int) {
	for i := 1; i < len(values); i++ {
		if cmp.Less(values[i], values[i-1]) {
			return false, i
		}
	}
	return true, len(values)
}
