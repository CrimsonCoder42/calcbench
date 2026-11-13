package sorting

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrLength: two slices that must have the same length do not.
var ErrLength = errors.New("sorting: slices differ in length")

// LexSortIndices returns the stable permutation that sorts records lexicographically by several keys:
// record i is (keys[0][i], keys[1][i], …), compared on keys[0] first, then keys[1] to break ties, and so
// on; records equal on every key keep ascending index order. keys[0] is the primary key (numpy's lexsort
// takes the primary key last; pass the keys reversed to match it). Values compare with cmp.Compare (NaN
// first, -0 equals +0). The result is written into dst (reused when cap(dst) >= n, else one new slice);
// with no keys the result is empty. keys are read, never modified.
// Errors: ErrLength when the keys differ in length (dst is returned unwritten).
//
// ID:      ALG-SORT-011
// Time:    O(k·n log n)   n = records, k = len(keys) (a comparison reads at most k keys)
// Space:   O(n) extra (one allocation, the merge buffer; none when n <= 12); O(n) for the result
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2 (sorting on multiple keys by lexicographic order of the key tuple) and §5.2.4 (sorting by merging)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.5 (multikey sorting; pointer sorting)
// Tags:    sorting, lexsort, multikey, lexicographic, argsort, indirect-sort, stable, merge-sort
func LexSortIndices[T cmp.Ordered](dst []int, keys [][]T) ([]int, error) {
	if len(keys) == 0 {
		return dst[:0], nil
	}
	n := len(keys[0])
	for j, key := range keys {
		if len(key) != n {
			return dst, fmt.Errorf("%w: len(keys[%d]) = %d, len(keys[0]) = %d", ErrLength, j, len(key), n)
		}
	}
	idx := resize(dst, n)
	for i := range idx {
		idx[i] = i
	}
	var buf []int
	if n > smallRun {
		buf = make([]int, n)
	}
	mergeSortFunc(idx, buf, func(a, b int) int {
		for _, key := range keys {
			if c := cmp.Compare(key[a], key[b]); c != 0 {
				return c
			}
		}
		return 0 // stability keeps a before b when a < b.
	})
	return idx, nil
}
