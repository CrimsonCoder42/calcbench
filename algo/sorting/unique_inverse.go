package sorting

import "cmp"

// UniqueInverse returns the distinct values of values in ascending order and, for each input, the
// position of its value among them: values[i] == uniques[inverse[i]] (up to -0/+0 and NaN). This is label
// encoding: inverse maps each value to a dense code 0…u-1 that preserves order. Values compare with
// cmp.Compare: all NaNs share one code, the smallest, and -0 and +0 share one. Empty input gives empty,
// non-nil slices. values is read, never modified.
//
// ID:      ALG-SORT-014
// Time:    O(n log n)   n = len(values)
// Space:   O(n) extra (the sorting permutation and merge buffer); O(n) for inverse, O(u) for uniques
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2 (sorting brings equal keys together) and §5.2.4 (stable merge sort)
// Tags:    sorting, unique, distinct, inverse, label-encoding, rank, codes
func UniqueInverse[T cmp.Ordered](values []T) (uniques []T, inverse []int) {
	uniques, _, inverse, _ = uniqueSorted(values, false, true, false)
	return uniques, inverse
}
