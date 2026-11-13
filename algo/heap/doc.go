// Package heap holds binary-heap algorithms: the k largest elements of a sequence.
// It also holds the reusable generic binary heap (Heap, ordered by a less function given at
// construction), the k smallest pair sums, selection in a sorted matrix, and a running median.
//
// ID short form: HEAP (IDs are ALG-HEAP-NNN).
//
// Ordering: values are compared with cmp.Compare, a total order in which a floating-point NaN sorts
// before every other value and -0 equals +0 (doc/DECISIONS.md D-002). No function modifies its input.
// The package name is the topic; it does not wrap or replace container/heap.
// Heap orders by the caller's less function instead; TopK keeps its own unexported sift (S1), and Heap
// is the type other packages reuse. Methods on a nil or zero Heap never panic (Push reports false).
//
// Functions:
//
//	ALG-HEAP-001  TopK  the k largest elements, in descending order
//	ALG-HEAP-002  NewHeap  an empty Heap ordered by a less function
//	ALG-HEAP-003  Heapify  a Heap built from a copy of a slice in O(n) (Floyd)
//	ALG-HEAP-004  Heap.Push  add an element, O(log n)
//	ALG-HEAP-005  Heap.Pop  remove and return the first element, ok false when empty
//	ALG-HEAP-006  Heap.Peek  the first element without removing it
//	ALG-HEAP-007  Heap.Len  the number of elements
//	ALG-HEAP-008  KSmallestPairs  the k index pairs with the smallest sums from two ascending slices
//	ALG-HEAP-009  KthSmallestInSortedMatrix  the k-th smallest element of a row- and column-ascending matrix
//	ALG-HEAP-010  NewMedianStream  a presized running-median stream (the zero MedianStream is ready too)
//	ALG-HEAP-011  MedianStream.Add  add a value, O(log n)
//	ALG-HEAP-012  MedianStream.Median  the lower and upper medians so far, O(1)
//	ALG-HEAP-013  MedianStream.Len  the number of values added
package heap
