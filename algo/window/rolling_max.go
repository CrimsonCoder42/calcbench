package window

import "cmp"

// RollingMax appends to dst the greatest element of each window of windowLen consecutive values and
// returns the extended slice: n-windowLen+1 maxima for n values, none when windowLen > n.
// Order is cmp.Compare's, so NaN is the least value and is the maximum only of an all-NaN window.
// values is read, never modified.
// Errors: ErrWindowLen when windowLen < 1; dst is then returned unchanged.
//
// ID:      ALG-WINDOW-004
// Time:    O(n)   n = len(values); at most 2n comparisons
// Space:   O(w) extra (w = windowLen), no allocation when w <= 64; O(n-w+1) for the result, one allocation only if dst lacks capacity
// Source:  Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element", Nordic Journal of Computing, 2006, 13(4):328–339, arXiv:cs/0610046
// Tags:    window, sliding-window, rolling, maximum, monotonic-deque, filter, time-series
func RollingMax[T cmp.Ordered](dst, values []T, windowLen int) ([]T, error) {
	return appendWindowExtremes(dst, values, windowLen, true)
}
