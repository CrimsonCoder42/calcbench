package interval

import (
	"errors"
	"fmt"
)

// ErrNotSorted: a slice that must be ascending has a value smaller than the one before it.
var ErrNotSorted = errors.New("interval: values are not in ascending order")

// integer is the set of types whose consecutive values the integer range functions can count.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// MissingRanges returns the integers of [lower, upper] that are not in values, as closed ranges sorted
// and disjoint. values must be ascending; duplicates and values outside [lower, upper] are allowed and
// ignored. Every integer in range present gives an empty, non-nil slice. values is read, never modified.
// Works up to the extremes of T with no overflow.
// Errors: ErrInverted when upper < lower; ErrNotSorted, with the index, when values is not ascending.
//
// ID:      ALG-INTVL-005
// Time:    O(n)   n = len(values)
// Space:   O(1) extra; O(g) for the result, g = number of gaps ≤ n + 1 (amortized appends)
// Source:  Klee, "Can the measure of ∪[aᵢ,bᵢ] be computed in less than O(n log n) steps?", American Mathematical Monthly, 1977, 84(4):284–285, doi:10.1080/00029890.1977.11994336 (the gaps of a union of intervals, read in one sweep once sorted; here the unit intervals [v, v+1))
// Source:  Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4 (one-pass merge of two sorted sequences: the range and values)
// Tags:    interval, ranges, gaps, missing, complement, integers, sorted
func MissingRanges[T integer](values []T, lower, upper T) ([]Interval[T], error) {
	if upper < lower {
		return nil, fmt.Errorf("%w: lower = %v, upper = %v", ErrInverted, lower, upper)
	}
	out := []Interval[T]{}
	next := lower // the smallest integer of the range not yet accounted for
	for i, v := range values {
		if i > 0 && v < values[i-1] {
			return nil, fmt.Errorf("%w: values[%d] = %v < values[%d] = %v", ErrNotSorted, i, v, i-1, values[i-1])
		}
		if v < next || v > upper {
			continue // a duplicate, below the range, or above it (the sort check still runs)
		}
		if v > next {
			out = append(out, Interval[T]{next, v - 1})
		}
		if v == upper {
			return checkRestSorted(values, i, out)
		}
		next = v + 1
	}
	return append(out, Interval[T]{next, upper}), nil
}

// checkRestSorted finishes the ascending check of values after index i once the result is complete.
func checkRestSorted[T integer](values []T, i int, out []Interval[T]) ([]Interval[T], error) {
	for j := i + 1; j < len(values); j++ {
		if values[j] < values[j-1] {
			return nil, fmt.Errorf("%w: values[%d] = %v < values[%d] = %v", ErrNotSorted, j, values[j], j-1, values[j-1])
		}
	}
	return out, nil
}
