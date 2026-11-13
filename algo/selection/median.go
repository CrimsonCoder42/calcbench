package selection

import (
	"cmp"
	"fmt"
	"slices"
)

// Median returns the lower and upper medians of values: the elements a sorted copy would hold at
// indexes (n-1)/2 and n/2. They are equal when n is odd. For numbers, the arithmetic median is
// (lower+upper)/2, left to the caller so that any ordered type works.
// values is read, never modified: a copy is reordered.
// Errors: ErrEmpty when values is empty.
//
// ID:      ALG-SELECT-002
// Time:    O(n) worst case   n = len(values)
// Space:   O(n) extra (one allocation, the copy) plus O(log n) stack
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 9 introduction (lower and upper median)
// Source:  Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience, 1997, 27(8):983–993, doi:10.1002/(SICI)1097-024X(199708)27:8<983::AID-SPE117>3.0.CO;2-#
// Tags:    selection, median, order-statistic, introselect
func Median[T cmp.Ordered](values []T) (lower, upper T, err error) {
	n := len(values)
	if n == 0 {
		return lower, upper, fmt.Errorf("%w: values", ErrEmpty)
	}
	work := slices.Clone(values)
	mid := (n - 1) / 2
	selectInPlace(work, mid)
	lower = work[mid]
	if n%2 == 1 {
		return lower, lower, nil
	}
	// Everything right of mid is >= lower, so the upper median is the least of it.
	upper = work[mid+1]
	for _, v := range work[mid+2:] {
		if cmp.Less(v, upper) {
			upper = v
		}
	}
	return lower, upper, nil
}
