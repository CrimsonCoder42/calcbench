package selection

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrEmpty: the input has no elements.
	ErrEmpty = errors.New("selection: empty input")
	// ErrRank: the requested rank is outside 1..n.
	ErrRank = errors.New("selection: rank out of range")
)

// KthSmallest returns the k-th smallest element of values, k = 1 being the minimum and k = len(values)
// the maximum: the element a sorted copy would hold at index k-1. Duplicates count separately.
// values is read, never modified: a copy is reordered.
// Errors: ErrEmpty when values is empty; ErrRank when k < 1 or k > len(values).
//
// ID:      ALG-SELECT-001
// Time:    O(n) worst case   n = len(values)
// Space:   O(n) extra (one allocation, the copy) plus O(log n) stack
// Source:  Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience, 1997, 27(8):983–993, doi:10.1002/(SICI)1097-024X(199708)27:8<983::AID-SPE117>3.0.CO;2-#
// Source:  Blum, Floyd, Pratt, Rivest, Tarjan, "Time bounds for selection", Journal of Computer and System Sciences, 1973, 7(4):448–461, doi:10.1016/S0022-0000(73)80033-9
// Source:  Hoare, "Algorithm 65: Find", Communications of the ACM, 1961, 4(7):321–322, doi:10.1145/366622.366647
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §9.2–9.3
// Tags:    selection, order-statistic, quickselect, introselect, median-of-medians, kth
func KthSmallest[T cmp.Ordered](values []T, k int) (T, error) {
	var zero T
	if len(values) == 0 {
		return zero, fmt.Errorf("%w: values", ErrEmpty)
	}
	if k < 1 || k > len(values) {
		return zero, fmt.Errorf("%w: k = %d, want 1..%d", ErrRank, k, len(values))
	}
	work := slices.Clone(values)
	selectInPlace(work, k-1)
	return work[k-1], nil
}
