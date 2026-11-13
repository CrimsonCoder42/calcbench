package prefix

import "fmt"

// SummedAreaTable holds the two-dimensional prefix sums of a rows × cols grid: entry (r, c) is the sum
// of grid[0:r][0:c]. Build one with PrefixSums2D and query it with RectangleSum. The zero value is the
// table of an empty grid.
type SummedAreaTable[T integer] struct {
	rows, cols int
	sums       []T // (rows+1) × (cols+1), row-major; row 0 and column 0 are zero
}

// Rows returns the number of rows of the grid the table was built from.
//
// ID:      ALG-PREFIX-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Crow, "Summed-area tables for texture mapping", ACM SIGGRAPH Computer Graphics (Proc. SIGGRAPH '84), 1984, 18(3):207–212, doi:10.1145/964965.808600
// Tags:    summed-area-table, dimensions
func (t SummedAreaTable[T]) Rows() int { return t.rows }

// Cols returns the number of columns of the grid the table was built from.
//
// ID:      ALG-PREFIX-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Crow, "Summed-area tables for texture mapping", ACM SIGGRAPH Computer Graphics (Proc. SIGGRAPH '84), 1984, 18(3):207–212, doi:10.1145/964965.808600
// Tags:    summed-area-table, dimensions
func (t SummedAreaTable[T]) Cols() int { return t.cols }

// PrefixSums2D returns the summed-area table of grid (its two-dimensional prefix sums), from which
// RectangleSum gives the sum of any axis-aligned rectangle in O(1). grid is rows × cols, row-major;
// every row must have the same length. An empty grid (no rows, or rows of length 0) gives an empty table.
// grid is read, never modified.
// Errors: ErrRagged when a row's length differs from row 0's; ErrOverflow when the sum of a rectangle
// grid[0:r][0:c], or of a row prefix grid[r][0:c], does not fit T.
//
// ID:      ALG-PREFIX-003
// Time:    O(r·c)   r = rows, c = columns; 2·r·c checked additions
// Space:   O((r+1)·(c+1)) for the result (one allocation); O(1) extra
// Source:  Crow, "Summed-area tables for texture mapping", ACM SIGGRAPH Computer Graphics (Proc. SIGGRAPH '84), 1984, 18(3):207–212, doi:10.1145/964965.808600
// Source:  Viola, Jones, "Robust Real-Time Face Detection", International Journal of Computer Vision, 2004, 57(2):137–154, doi:10.1023/B:VISI.0000013087.49260.fb (§2.1, the integral image)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.1 "Static array queries", two-dimensional prefix sums, pp. 84–85, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    prefix-sum, summed-area-table, integral-image, 2d, range-query, integer, overflow-checked
func PrefixSums2D[T integer](grid [][]T) (SummedAreaTable[T], error) {
	rows := len(grid)
	cols := 0
	if rows > 0 {
		cols = len(grid[0])
	}
	for r, row := range grid {
		if len(row) != cols {
			return SummedAreaTable[T]{}, fmt.Errorf("%w: len(grid[%d]) = %d, len(grid[0]) = %d", ErrRagged, r, len(row), cols)
		}
	}
	stride := cols + 1
	sums := make([]T, (rows+1)*stride)
	for r, row := range grid {
		above := sums[r*stride : (r+1)*stride]
		here := sums[(r+1)*stride : (r+2)*stride]
		var rowSum T
		for c, v := range row {
			// rowSum is grid[r][0:c+1]; here[c+1] is grid[0:r+1][0:c+1]. Both are true sums, so an
			// overflow here means a real rectangle's sum does not fit.
			s, ok := checkedAdd(rowSum, v)
			if !ok {
				return SummedAreaTable[T]{}, fmt.Errorf("%w: row sum through grid[%d][%d]", ErrOverflow, r, c)
			}
			rowSum = s
			s, ok = checkedAdd(above[c+1], rowSum)
			if !ok {
				return SummedAreaTable[T]{}, fmt.Errorf("%w: rectangle sum through grid[%d][%d]", ErrOverflow, r, c)
			}
			here[c+1] = s
		}
	}
	return SummedAreaTable[T]{rows: rows, cols: cols, sums: sums}, nil
}

// RectangleSum returns the sum of grid[row0:row1][col0:col1] (half-open in both axes) from the
// table: S(row1, col1) − S(row0, col1) − S(row1, col0) + S(row0, col0). An empty rectangle gives 0.
// The four terms are combined in 128-bit arithmetic, so the result is exact whenever it fits T.
// Errors: ErrRange unless 0 <= row0 <= row1 <= Rows() and 0 <= col0 <= col1 <= Cols(); ErrOverflow
// when the rectangle's sum does not fit T.
//
// ID:      ALG-PREFIX-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Crow, "Summed-area tables for texture mapping", ACM SIGGRAPH Computer Graphics (Proc. SIGGRAPH '84), 1984, 18(3):207–212, doi:10.1145/964965.808600
// Source:  Viola, Jones, "Robust Real-Time Face Detection", International Journal of Computer Vision, 2004, 57(2):137–154, doi:10.1023/B:VISI.0000013087.49260.fb (§2.1, a rectangle from four array references)
// Tags:    prefix-sum, summed-area-table, integral-image, 2d, range-sum, range-query, integer, overflow-checked
func (t SummedAreaTable[T]) RectangleSum(row0, col0, row1, col1 int) (T, error) {
	if row0 < 0 || row1 < row0 || row1 > t.rows || col0 < 0 || col1 < col0 || col1 > t.cols {
		return 0, fmt.Errorf("%w: rows [%d, %d), cols [%d, %d) in a %d × %d table", ErrRange, row0, row1, col0, col1, t.rows, t.cols)
	}
	if row0 == row1 || col0 == col1 {
		return 0, nil
	}
	stride := t.cols + 1
	at := func(r, c int) int128 { return widen(t.sums[r*stride+c]) }
	sum := at(row1, col1).sub(at(row0, col1)).sub(at(row1, col0)).add(at(row0, col0))
	v, ok := narrow[T](sum)
	if !ok {
		return 0, fmt.Errorf("%w: rows [%d, %d), cols [%d, %d)", ErrOverflow, row0, row1, col0, col1)
	}
	return v, nil
}
