package search

import "cmp"

// SearchSortedMatrixSaddleback returns a position of target in matrix, whose rows are each ascending
// and whose columns are each ascending (top to bottom). found is false and row, col are -1 when target
// is absent. Saddleback (staircase) search: start at the top-right corner; a value above target rules
// out its column (step left), a value below target rules out its row (step down). Rows must all have
// the length of the first; a shorter row is skipped without being read, so nothing panics. matrix is
// read, never modified.
//
// ID:      ALG-SEARCH-010
// Time:    O(m + n)   m = len(matrix), n = len(matrix[0]); at most m + n - 1 comparisons
// Space:   O(1) extra, no allocation
// Source:  Gries, The Science of Programming, Springer, 1981, p. 215 (saddleback search; page from a search summary, not checked against a copy)
// Source:  Bird, "Improving saddleback search: a lesson in algorithm design", Mathematics of Program Construction (MPC 2006), LNCS 4014, Springer, 2006, pp. 82–89, doi:10.1007/11783596_8
// Tags:    search, matrix, saddleback, staircase, sorted-rows-and-columns
func SearchSortedMatrixSaddleback[T cmp.Ordered](matrix [][]T, target T) (row, col int, found bool) {
	if len(matrix) == 0 {
		return -1, -1, false
	}
	r, c := 0, len(matrix[0])-1
	for r < len(matrix) && c >= 0 {
		cells := matrix[r]
		if c >= len(cells) {
			r++ // a short row breaks the precondition; skip it rather than index past its end.
			continue
		}
		switch cmp.Compare(cells[c], target) {
		case 0:
			return r, c, true
		case 1:
			c--
		default:
			r++
		}
	}
	return -1, -1, false
}
