package window

import "cmp"

// RollingMin appends to dst the least element of each window of windowLen consecutive values and
// returns the extended slice: n-windowLen+1 minima for n values, none when windowLen > n.
// Order is cmp.Compare's, so a window holding a NaN has NaN as its minimum.
// values is read, never modified.
// Errors: ErrWindowLen when windowLen < 1; dst is then returned unchanged.
//
// ID:      ALG-WINDOW-003
// Time:    O(n)   n = len(values); at most 2n comparisons
// Space:   O(w) extra (w = windowLen), no allocation when w <= 64; O(n-w+1) for the result, one allocation only if dst lacks capacity
// Source:  Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element", Nordic Journal of Computing, 2006, 13(4):328–339, arXiv:cs/0610046
// Tags:    window, sliding-window, rolling, minimum, monotonic-deque, filter, time-series
func RollingMin[T cmp.Ordered](dst, values []T, windowLen int) ([]T, error) {
	return appendWindowExtremes(dst, values, windowLen, false)
}
