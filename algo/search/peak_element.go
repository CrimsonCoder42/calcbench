package search

import "cmp"

// PeakElement returns the index of a peak of values: an element not less than its neighbours, where a
// position outside the slice counts as lower than everything. On a unimodal slice (rising, then
// falling) it is the maximum; when no two adjacent values are equal it is a strict local maximum.
// ok is false only for an empty slice (index 0). values is read, never modified.
// It binary searches on values[i] < values[i+1]: if the next value is larger, a peak lies to the right.
//
// ID:      ALG-SEARCH-007
// Time:    O(log n)   n = len(values); at most ⌈log2 n⌉ comparisons
// Space:   O(1) extra, no allocation
// Source:  Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, ch. 5, Solved Exercise 1 (peak entry of a unimodal array with O(log n) probes)
// Source:  Kiefer, "Sequential minimax search for a maximum", Proceedings of the American Mathematical Society, 1953, 4(3):502–506, doi:10.1090/S0002-9939-1953-0055639-3
// Tags:    search, binary-search, peak, local-maximum, unimodal, bitonic
func PeakElement[T cmp.Ordered](values []T) (index int, ok bool) {
	if len(values) == 0 {
		return 0, false
	}
	// Invariant: lo is 0 or values[lo-1] < values[lo]; hi is n-1 or values[hi] >= values[hi+1].
	// So a peak lies in [lo, hi], and when lo == hi it is lo.
	lo, hi := 0, len(values)-1
	for lo < hi {
		mid := lo + (hi-lo)/2
		if cmp.Less(values[mid], values[mid+1]) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo, true
}
