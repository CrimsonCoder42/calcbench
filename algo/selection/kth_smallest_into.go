package selection

import (
	"cmp"
	"fmt"
)

// KthSmallestInto returns the k-th smallest element of values (k = 1 being the minimum), selecting in a
// caller-supplied buffer so that it does not allocate. values is copied into dst, which is reused when
// cap(dst) >= len(values) and otherwise replaced by one new slice; work is that buffer, length
// len(values), reordered so that work[:k] holds the k smallest elements (in no particular order) and
// work[k-1] is the k-th. dst may be values itself to select in place; any other overlap is unsupported.
// values is read, never modified, unless it is dst.
// Errors: ErrEmpty when values is empty; ErrRank when k < 1 or k > len(values). On error work is dst,
// unwritten.
//
// ID:      ALG-SELECT-004
// Time:    O(n) worst case   n = len(values)
// Space:   O(1) extra heap when cap(dst) >= n (no allocation), else O(n) (one allocation, the buffer); plus O(log n) stack
// Source:  Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience, 1997, 27(8):983–993, doi:10.1002/(SICI)1097-024X(199708)27:8<983::AID-SPE117>3.0.CO;2-#
// Source:  Blum, Floyd, Pratt, Rivest, Tarjan, "Time bounds for selection", Journal of Computer and System Sciences, 1973, 7(4):448–461, doi:10.1016/S0022-0000(73)80033-9
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §9.2–9.3
// Tags:    selection, order-statistic, quickselect, introselect, no-alloc, dst
func KthSmallestInto[T cmp.Ordered](dst, values []T, k int) (kth T, work []T, err error) {
	n := len(values)
	if n == 0 {
		return kth, dst, fmt.Errorf("%w: values", ErrEmpty)
	}
	if k < 1 || k > n {
		return kth, dst, fmt.Errorf("%w: k = %d, want 1..%d", ErrRank, k, n)
	}
	if cap(dst) >= n {
		work = dst[:n]
	} else {
		work = make([]T, n)
	}
	copy(work, values) // a no-op when work and values share their first element.
	selectInPlace(work, k-1)
	return work[k-1], work, nil
}
