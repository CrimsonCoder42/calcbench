package search

import "cmp"

// SearchSortedMatrix returns the position of target in matrix, whose rows read one after another form
// one ascending list: each row is ascending and each row's first element is greater than the previous
// row's last. found is false and row, col are -1 when target is absent. It binary searches for the row
// by first element, then binary searches that row. Rows may differ in length but must not be empty; on
// input that breaks the precondition the result is unspecified and nothing panics. matrix is read,
// never modified. For matrices whose rows and columns are sorted separately use
// SearchSortedMatrixSaddleback.
//
// ID:      ALG-SEARCH-009
// Time:    O(log m + log n) = O(log(m·n))   m = len(matrix), n = length of the row searched
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search), over the rows and then within one row
// Tags:    search, binary-search, matrix, sorted, two-level
func SearchSortedMatrix[T cmp.Ordered](matrix [][]T, target T) (row, col int, found bool) {
	// The first row whose first element is greater than target; target can only be in the row before it.
	next, _ := SearchPredicate(0, len(matrix), func(r int) bool {
		return len(matrix[r]) > 0 && cmp.Less(target, matrix[r][0])
	})
	if next == 0 {
		return -1, -1, false
	}
	row = next - 1
	cells := matrix[row]
	col = LowerBound(cells, target)
	if col < len(cells) && cmp.Compare(cells[col], target) == 0 {
		return row, col, true
	}
	return -1, -1, false
}
