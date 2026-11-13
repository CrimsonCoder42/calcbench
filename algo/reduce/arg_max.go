package reduce

import (
	"cmp"
	"fmt"
)

// ArgMax returns the index of the largest value, the first one when several are equal. For float types
// -0 and +0 are equal (the first wins) and +Inf is an ordinary value.
// Errors: ErrEmpty when values is empty; ErrNotFinite (with the index) when a value is NaN, since a NaN
// has no place in the order.
//
// ID:      ALG-RED-009
// Time:    O(n)   n = len(values); n - 1 comparisons, which is optimal (Knuth §1.2.10)
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.2.10, Algorithm M (find the maximum)
// Tags:    argmax, maximum, index, reduction, first-index, no-alloc
func ArgMax[T cmp.Ordered](values []T) (int, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("%w: values", ErrEmpty)
	}
	best := 0
	for i, v := range values {
		if v != v { // only a NaN is not equal to itself.
			return 0, fmt.Errorf("%w: values[%d] is NaN", ErrNotFinite, i)
		}
		if v > values[best] {
			best = i
		}
	}
	return best, nil
}
