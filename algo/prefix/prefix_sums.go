package prefix

import (
	"errors"
	"fmt"
	"slices"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrOverflow: a result does not fit the element type.
	ErrOverflow = errors.New("prefix: result does not fit the integer type")
	// ErrRange: a bound is outside the array or the lower bound exceeds the upper one.
	ErrRange = errors.New("prefix: index out of range")
	// ErrRagged: the rows of a grid have different lengths.
	ErrRagged = errors.New("prefix: grid rows have different lengths")
	// ErrEmpty: the input has no elements.
	ErrEmpty = errors.New("prefix: empty input")
)

// PrefixSums appends to dst the n+1 prefix sums of values and returns the extended slice: element i is
// values[0] + … + values[i-1], so element 0 is 0 and element n is the total. The sum of values[lo:hi]
// is then RangeSum(prefix, lo, hi) in O(1).
// values is read, never modified. dst must not overlap values.
// Errors: ErrOverflow when a prefix sum does not fit T (the index of the first value that overflows is
// in the message). On error dst is returned with its length unchanged.
//
// ID:      ALG-PREFIX-001
// Time:    O(n)   n = len(values); n checked additions
// Space:   O(1) extra; O(n+1) for the result, one allocation only if dst lacks capacity
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.1 "Static array queries", "Sum queries", p. 84, https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Blelloch, "Prefix Sums and Their Applications", in Reif (ed.), Synthesis of Parallel Algorithms, Morgan Kaufmann, 1993, ch. 1, §1.1 (the scan and prescan definitions); also CMU-CS-90-190, 1990
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8 "Algorithm Design Techniques", §8.2 (the cumulative array cumarr)
// Tags:    prefix-sum, cumulative-sum, scan, range-query, integer, overflow-checked
func PrefixSums[T integer](dst, values []T) ([]T, error) {
	out := slices.Grow(dst, len(values)+1)
	var sum T
	out = append(out, sum)
	for i, v := range values {
		s, ok := checkedAdd(sum, v)
		if !ok {
			return dst, fmt.Errorf("%w: prefix sum through values[%d]", ErrOverflow, i)
		}
		sum = s
		out = append(out, sum)
	}
	return out, nil
}
