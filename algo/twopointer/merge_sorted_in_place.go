package twopointer

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrLength: a length or buffer size does not fit the call.
var ErrLength = errors.New("twopointer: length out of range")

// MergeSortedInPlace merges the ascending slice other into the ascending prefix buf[:m], writing the
// merged ascending result to buf[:m+len(other)]. buf must have room: len(buf) >= m + len(other);
// elements of buf beyond m are overwritten. The merge runs back to front, so no element of buf[:m] is
// overwritten before it is read and no scratch space is needed. It is stable: equal elements of buf
// come before those of other. Order is cmp.Compare's, so NaN sorts first.
// buf IS MODIFIED; other is read, never modified, and must not overlap buf.
// Errors: ErrLength when m < 0 or m + len(other) > len(buf); buf is then unchanged.
//
// ID:      ALG-TWOPTR-004
// Time:    O(m + k)   m = sorted elements already in buf, k = len(other); at most m + k − 1 comparisons
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.4 "Sorting by Merging", Algorithm M (two-way merge), p. 158 (run here from the high ends)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §2.3.1, procedure MERGE
// Tags:    two-pointers, merge, sorted, in-place, stable
func MergeSortedInPlace[T cmp.Ordered](buf []T, m int, other []T) error {
	k := len(other)
	if m < 0 || m > len(buf) || k > len(buf)-m {
		return fmt.Errorf("%w: m = %d, len(other) = %d, len(buf) = %d", ErrLength, m, k, len(buf))
	}
	i, j := m-1, k-1
	for w := m + k - 1; j >= 0; w-- {
		// Take from buf only when strictly greater, so equal elements of other land after buf's.
		if i >= 0 && cmp.Compare(buf[i], other[j]) > 0 {
			buf[w] = buf[i]
			i--
		} else {
			buf[w] = other[j]
			j--
		}
	}
	// Whatever is left of buf[:m] is already in place.
	return nil
}
