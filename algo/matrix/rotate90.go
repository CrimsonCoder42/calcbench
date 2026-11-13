package matrix

import (
	"errors"
	"fmt"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrNotSquare: a square grid was required.
	ErrNotSquare = errors.New("matrix: grid is not square")
	// ErrRagged: the rows of a grid have different lengths.
	ErrRagged = errors.New("matrix: grid rows have different lengths")
	// ErrEmpty: the grid has no cells.
	ErrEmpty = errors.New("matrix: grid is empty")
	// ErrShape: the operands' dimensions do not fit together.
	ErrShape = errors.New("matrix: dimensions do not match")
	// ErrInvalidCSR: a CSR matrix is malformed.
	ErrInvalidCSR = errors.New("matrix: malformed CSR matrix")
	// ErrOverflow: an integer sum or product does not fit its type.
	ErrOverflow = errors.New("matrix: integer overflow")
	// ErrNotFound: no candidate satisfies the bound.
	ErrNotFound = errors.New("matrix: no sub-rectangle satisfies the bound")
)

// Rotate90InPlace rotates the square grid a quarter turn clockwise, in place: the cell at (r, c) moves
// to (c, n−1−r), so the first column read bottom-up becomes the first row. An empty grid is left as is.
// grid is modified. Errors: ErrNotSquare when some row's length differs from the number of rows (grid is
// then unchanged).
//
// ID:      ALG-MATRIX-001
// Time:    O(n²)   n = number of rows; each cell is moved once
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.3.3 (a permutation applied in place cycle by cycle); the rotation's cycles have length 4, or 1 at the centre
// Tags:    matrix, rotate, in-place, permutation-cycles, grid
func Rotate90InPlace[T any](grid [][]T) error {
	n := len(grid)
	for r, row := range grid {
		if len(row) != n {
			return fmt.Errorf("%w: %d rows, row %d has %d cells", ErrNotSquare, n, r, len(row))
		}
	}
	for r := range n / 2 {
		for c := r; c < n-1-r; c++ {
			// One 4-cycle: (r,c) ← (n−1−c,r) ← (n−1−r,n−1−c) ← (c,n−1−r) ← (r,c).
			tmp := grid[r][c]
			grid[r][c] = grid[n-1-c][r]
			grid[n-1-c][r] = grid[n-1-r][n-1-c]
			grid[n-1-r][n-1-c] = grid[c][n-1-r]
			grid[c][n-1-r] = tmp
		}
	}
	return nil
}
