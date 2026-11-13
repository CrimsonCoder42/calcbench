// Package sorting sorts slices and answers the questions sorting answers: is it sorted, which elements
// are distinct, which k are smallest, and in what order the indexes of a slice would sort.
// The name is sorting, not sort, so that callers can import the standard sort package beside it.
//
// ID short form: SORT (IDs are ALG-SORT-NNN).
//
// Ordering: functions over cmp.Ordered compare with cmp.Compare, a total order in which a floating-point
// NaN sorts before every other value and -0 equals +0 (doc/DECISIONS.md D-002). Functions that take a
// compare func follow the slices.SortFunc contract (negative, zero or positive); a nil compare returns
// ErrNilCompare.
//
// Output: every function that produces a slice writes it into dst. dst is reused when its capacity is
// large enough, otherwise one new slice is allocated, so a caller that keeps dst between calls stops
// allocating. dst may be the input itself (the same first element) to work in place; any other overlap
// between dst and the input is not supported. The input is never modified unless it is dst.
// The Unique functions return several slices, so they take no dst and allocate their results.
//
// Small ranges are finished by insertion sort. The cutoff, 12 elements, is inside the 5 to 15 that
// Sedgewick and Wayne, Algorithms, 4th ed., §2.3 give as working well; it was not measured here.
//
// Functions:
//
//	ALG-SORT-001  MergeSortStable    stable sort by a compare func (bottom-up merge sort)
//	ALG-SORT-002  InsertionSort      sort for small or nearly sorted input
//	ALG-SORT-003  CountingSort       sort of ints in [0, maxKey] in O(n + maxKey)
//	ALG-SORT-004  RadixSortUint64    LSD radix sort of uint64 keys, 8-bit digits
//	ALG-SORT-005  BucketSort         sort of floats in [0, 1), expected linear time
//	ALG-SORT-006  IsSorted           whether a slice is ascending, and where it first is not
//	ALG-SORT-007  DedupeSorted       the distinct values of a sorted slice
//	ALG-SORT-008  PartialSort        the k smallest values, ascending
//	ALG-SORT-009  SortIndices        the stable sorting permutation (argsort)
//	ALG-SORT-010  ThreeWayPartition  less / equal / greater blocks around a pivot (Dutch national flag)
//	ALG-SORT-011  LexSortIndices     the stable permutation sorting records by several keys (lexsort)
//	ALG-SORT-012  UniqueAll          distinct values with first index, inverse and counts
//	ALG-SORT-013  UniqueCounts       distinct values and how often each occurs
//	ALG-SORT-014  UniqueInverse      distinct values and each input's code among them (label encoding)
package sorting
