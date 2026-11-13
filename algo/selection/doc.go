// Package selection finds order statistics (the k-th smallest element, the median) in linear time
// without sorting.
//
// ID short form: SELECT (IDs are ALG-SELECT-NNN).
//
// Ordering: values are compared with cmp.Compare, a total order in which a floating-point NaN sorts
// before every other value and -0 equals +0 (doc/DECISIONS.md D-002). No function modifies its input:
// each copies it once and reorders the copy. The algorithm is introselect (D-001): O(n) in the worst
// case, with no randomness.
//
// KthSmallestInto copies into the caller's dst instead of a new slice; dst may be the input itself, which
// is then reordered in place.
//
// Functions:
//
//	ALG-SELECT-001  KthSmallest  the k-th smallest element, k = 1 being the minimum
//	ALG-SELECT-002  Median       the lower and upper middle elements
//	ALG-SELECT-003  KthLargest       the k-th largest element, k = 1 being the maximum
//	ALG-SELECT-004  KthSmallestInto  the k-th smallest, selected in a caller-supplied buffer (no allocation)
//	ALG-SELECT-005  KthSmallestIndex index of the k-th smallest and the index partition around it (argpartition)
package selection
