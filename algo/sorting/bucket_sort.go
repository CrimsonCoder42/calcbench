package sorting

import (
	"errors"
	"fmt"
	"slices"
)

// ErrOutOfRange: a value is not in the half-open range [0, 1), or is NaN.
var ErrOutOfRange = errors.New("sorting: value out of range")

// BucketSort returns values sorted ascending, for values in [0, 1). It spreads the n values over n
// equal-width buckets (bucket ⌊n·x⌋), then sorts each bucket: by insertion sort when it holds 12 or
// fewer, otherwise by slices.Sort, so a crowded bucket cannot make the whole sort quadratic. -0 is
// accepted and sorts with +0. The result is written into dst (reused when cap(dst) >= len(values), else
// one new slice); dst may be values itself (a copy of values is then made), but must not otherwise
// overlap values: a partial overlap is unsupported and gives an unspecified result. values is never
// modified unless it is dst.
// Errors: ErrOutOfRange when a value is < 0, >= 1 or NaN; checked before dst is written, and on error dst
// is returned unwritten.
//
// ID:      ALG-SORT-005
// Time:    O(n) expected when values are independent and uniform on [0, 1); O(n log n) worst   n = len(values)
// Space:   O(n) extra (one allocation, n+1 bucket bounds; a second, the copy, when dst is values); O(n) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.4 (BUCKET-SORT) and Exercise 8.4-2 (worst case)
// Tags:    sorting, bucket-sort, distribution-sort, floats, linear-time-expected
func BucketSort(dst, values []float64) ([]float64, error) {
	for i, v := range values {
		if !(v >= 0 && v < 1) { // false for NaN as well.
			return dst, fmt.Errorf("%w: values[%d] = %v, want [0, 1)", ErrOutOfRange, i, v)
		}
	}
	n := len(values)
	out := resize(dst, n)
	if n == 0 {
		return out, nil
	}
	src := values
	if &out[0] == &values[0] {
		src = slices.Clone(values) // the scatter below reads src while writing out.
	}
	// ends[b+1] counts bucket b; after the prefix sum ends[b] is where bucket b starts.
	ends := make([]int, n+1)
	for _, v := range src {
		ends[bucketIndex(v, n)+1]++
	}
	for b := 1; b <= n; b++ {
		ends[b] += ends[b-1]
	}
	for _, v := range src {
		b := bucketIndex(v, n)
		out[ends[b]] = v
		ends[b]++
	}
	// Now ends[b] is where bucket b ends; bucket b is out[start:ends[b]], start being the previous end.
	start := 0
	for b := range n {
		bucket := out[start:ends[b]]
		if len(bucket) <= smallRun {
			insertionSortInPlace(bucket)
		} else {
			slices.Sort(bucket)
		}
		start = ends[b]
	}
	return out, nil
}

// bucketIndex returns ⌊n·v⌋ for v in [0, 1), clamped to n-1: when v is just below 1, n·v can round up
// to n in floating point.
func bucketIndex(v float64, n int) int {
	return min(int(v*float64(n)), n-1)
}
