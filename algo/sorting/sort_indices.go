package sorting

import "fmt"

// SortIndices returns the stable sorting permutation of values (an argsort): the indexes 0…n-1 ordered
// so that values[idx[0]], values[idx[1]], … is ascending by compare, with equal values in ascending
// index order. values is read, never modified: only indexes move, so large elements are never copied.
// The result is written into dst (reused when cap(dst) >= len(values), else one new slice).
// compare follows the slices.SortFunc contract.
// Errors: ErrNilCompare when compare is nil (dst is returned unwritten).
//
// ID:      ALG-SORT-009
// Time:    O(n log n)   n = len(values)
// Space:   O(n) extra (one allocation, the merge buffer; none when n <= 12); O(n) for the result
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.5 (pointer sorting: sort references, not the data) and §2.2 (mergesort, stable)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.4 (sorting by merging)
// Tags:    sorting, argsort, indirect-sort, permutation, stable, merge-sort
func SortIndices[T any](dst []int, values []T, compare func(a, b T) int) ([]int, error) {
	if compare == nil {
		return dst, fmt.Errorf("%w: compare", ErrNilCompare)
	}
	n := len(values)
	idx := resize(dst, n)
	for i := range idx {
		idx[i] = i
	}
	var buf []int
	if n > smallRun {
		buf = make([]int, n)
	}
	// Stability of the merge keeps equal values in ascending index order, since idx starts ascending.
	mergeSortFunc(idx, buf, func(i, j int) int { return compare(values[i], values[j]) })
	return idx, nil
}
