package selection

import (
	"cmp"
	"fmt"
)

// KthLargest returns the k-th largest element of values, k = 1 being the maximum and k = len(values)
// the minimum: the element a sorted copy would hold at index len(values)-k. Duplicates count separately.
// values is read, never modified: a copy is reordered.
// Errors: ErrEmpty when values is empty; ErrRank when k < 1 or k > len(values).
//
// ID:      ALG-SELECT-003
// Time:    O(n) worst case   n = len(values)
// Space:   O(n) extra (one allocation, the copy) plus O(log n) stack
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 9 introduction (order statistics) and §9.2–9.3
// Source:  Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience, 1997, 27(8):983–993, doi:10.1002/(SICI)1097-024X(199708)27:8<983::AID-SPE117>3.0.CO;2-#
// Tags:    selection, order-statistic, kth-largest, introselect, quickselect
func KthLargest[T cmp.Ordered](values []T, k int) (T, error) {
	var zero T
	n := len(values)
	if n == 0 {
		return zero, fmt.Errorf("%w: values", ErrEmpty)
	}
	if k < 1 || k > n {
		return zero, fmt.Errorf("%w: k = %d, want 1..%d", ErrRank, k, n)
	}
	// The k-th largest is the (n-k+1)-th smallest; n-k+1 is in 1..n, so KthSmallest cannot fail here.
	return KthSmallest(values, n-k+1)
}
