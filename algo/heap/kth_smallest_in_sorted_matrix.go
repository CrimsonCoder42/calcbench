package heap

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrRank: k is outside 1…(number of elements).
var ErrRank = errors.New("heap: rank k is out of range")

// ErrRagged: the rows of a matrix do not all have the same length.
var ErrRagged = errors.New("heap: matrix rows differ in length")

// KthSmallestInSortedMatrix returns the k-th smallest element (k is 1-based; duplicates count
// separately) of a matrix whose every row and every column is ascending. The order is not checked: a
// matrix whose rows or columns are not ascending gives an unspecified element, never a panic.
// matrix is read, never modified.
// Errors: ErrRagged when rows differ in length; ErrRank when k < 1 or k > rows·columns (an empty
// matrix has no rank).
//
// ID:      ALG-HEAP-009
// Time:    O(r + k log m)   r = len(matrix), m = min(k, r)
// Space:   O(m) extra (the heap)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5, Exercise 6.5-9 (k-way merge of the rows with a min-heap, stopped at the k-th output)
// Source:  Frederickson, Johnson, "Generalized selection and ranking: sorted matrices", SIAM Journal on Computing, 1984, 13(1):14–30, doi:10.1137/0213002 (the problem; O(r) method not used, see catalog heap.md)
// Tags:    heap, selection, k-th-smallest, order-statistic, sorted-matrix, young-tableau
func KthSmallestInSortedMatrix[T cmp.Ordered](matrix [][]T, k int) (T, error) {
	var zero T
	cols := 0
	if len(matrix) > 0 {
		cols = len(matrix[0])
	}
	for i, row := range matrix {
		if len(row) != cols {
			return zero, fmt.Errorf("%w: len(matrix[%d]) = %d, len(matrix[0]) = %d", ErrRagged, i, len(row), cols)
		}
	}
	if k < 1 || cols == 0 || (k-1)/cols >= len(matrix) { // k > rows*cols without forming the product
		return zero, fmt.Errorf("%w: k = %d, matrix %d×%d", ErrRank, k, len(matrix), cols)
	}
	// With ascending columns, element (i, j) has the i+1 elements (0…i, j) at or below it, so the k-th
	// smallest lies in the first min(k, r) rows. Merge those rows with a min-heap of their heads.
	type cell struct{ i, j int }
	m := min(k, len(matrix))
	heads := make([]cell, m)
	for i := range heads {
		heads[i] = cell{i, 0}
	}
	h := &Heap[cell]{items: heads, less: func(p, q cell) bool {
		return cmp.Less(matrix[p.i][p.j], matrix[q.i][q.j])
	}}
	for i := m/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	for range k - 1 {
		c, _ := h.Pop()
		if c.j+1 < cols {
			h.Push(cell{c.i, c.j + 1})
		}
	}
	c, _ := h.Peek()
	return matrix[c.i][c.j], nil
}
