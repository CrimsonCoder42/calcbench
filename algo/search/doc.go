// Package search finds positions in sorted slices by binary search.
//
// ID short form: SEARCH (IDs are ALG-SEARCH-NNN).
//
// Ordering: values are compared with cmp.Compare, a total order in which a floating-point NaN sorts
// before every other value and -0 equals +0. A slice sorted by slices.Sort is sorted in this order.
// No function modifies its input. Sortedness is a precondition and is not checked: on unsorted input
// the result is still an index in [0, len(sorted)], and nothing panics.
//
// Functions that report a position return -1 (and found = false) when the target is absent; ALG-SEARCH-001…004
// return the insertion point instead. SearchPredicate searches a range of ints rather
// than a slice and returns ErrBounds or ErrNilPredicate on invalid arguments.
//
// Functions:
//
//	ALG-SEARCH-001  LowerBound         first index whose element is not less than the target
//	ALG-SEARCH-002  UpperBound         first index whose element is greater than the target
//	ALG-SEARCH-003  BinarySearchFirst  index of the first element equal to the target
//	ALG-SEARCH-004  BinarySearchLast   index of the last element equal to the target
//	ALG-SEARCH-005  SearchRotated                 index of the target in a rotated sorted slice of distinct values
//	ALG-SEARCH-006  SearchRotatedDup              the same with duplicates allowed (O(n) worst case)
//	ALG-SEARCH-007  PeakElement                   index of an element not less than its neighbours
//	ALG-SEARCH-008  SearchPredicate               smallest x in [lo, hi) where a monotone predicate is true
//	ALG-SEARCH-009  SearchSortedMatrix            position of the target in a row-major sorted matrix
//	ALG-SEARCH-010  SearchSortedMatrixSaddleback  position of the target in a row- and column-sorted matrix
package search
