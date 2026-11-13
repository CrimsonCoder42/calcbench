package sorting

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/selection"
)

// ErrNegativeK: k is less than 0.
var ErrNegativeK = errors.New("sorting: k must not be negative")

// PartialSort returns the k smallest elements of values in ascending cmp.Compare order, or all of them,
// sorted, when k >= len(values). Duplicates count separately. dst is the working buffer: it is reused
// when cap(dst) >= len(values), otherwise one new slice of len(values) is allocated; the result is its
// first k elements, and the rest of the buffer holds the other elements in no particular order. dst may
// be values itself to work in place. values is never modified unless it is dst. k = 0 or empty values
// returns dst[:0].
// Errors: ErrNegativeK when k < 0 (dst is returned unwritten).
//
// ID:      ALG-SORT-008
// Time:    O(n + k log k) worst case   n = len(values), k = min(k, n)
// Space:   O(1) extra heap when cap(dst) >= n, else O(n) (one allocation, the buffer); plus O(log n) stack
// Source:  Chambers, "Algorithm 410: Partial sorting", Communications of the ACM, 1971, 14(5):357–358, doi:10.1145/362588.362602
// Source:  Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience, 1997, 27(8):983–993, doi:10.1002/(SICI)1097-024X(199708)27:8<983::AID-SPE117>3.0.CO;2-#
// Tags:    sorting, partial-sort, selection, introselect, smallest-k, top-k
func PartialSort[T cmp.Ordered](dst, values []T, k int) ([]T, error) {
	if k < 0 {
		return dst, fmt.Errorf("%w: k = %d", ErrNegativeK, k)
	}
	k = min(k, len(values))
	if k == 0 {
		return dst[:0], nil
	}
	// KthSmallestInto (ALG-SELECT-004) leaves the k smallest in work[:k] in O(n) worst case.
	// 1 <= k <= len(values) here, so it cannot fail.
	_, work, err := selection.KthSmallestInto(dst, values, k)
	if err != nil {
		return dst, err
	}
	slices.Sort(work[:k])
	return work[:k], nil
}
