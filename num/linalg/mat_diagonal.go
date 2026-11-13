package linalg

// MatDiagonal returns the offset-th diagonal of a rows×cols matrix: entries a[i][i+offset] (offset > 0
// above the main diagonal, < 0 below), in dst (reused when its capacity suffices). A diagonal outside the
// matrix is empty (length 0), as numpy's diagonal.
// Errors: ErrDimension when len(a) ≠ rows·cols.
//
// ID:      NUM-LINALG-186
// Time:    O(min(rows, cols))
// Space:   O(1) extra; O(length) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §1.2.1 "Band Matrices" (the diagonals of a matrix, indexed by offset)
// Tags:    matrix, n-dimensional, diagonal, superdiagonal, subdiagonal, linear-algebra
// Units:   the unit of a
// Example: [1 2 3; 4 5 6] → offset 0 (1, 5), offset 1 (2, 6), offset −1 (4) (hand-checked)
func MatDiagonal(dst, a []float64, rows, cols, offset int) ([]float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return nil, err
	}
	r0, c0 := 0, offset
	if offset < 0 {
		r0, c0 = -offset, 0
	}
	length := 0
	if r0 < rows && c0 < cols && r0 >= 0 && c0 >= 0 {
		length = min(rows-r0, cols-c0)
	}
	if offset <= -rows || offset >= cols {
		length = 0
	}
	dst = grow(dst, length)
	for i := range length {
		dst[i] = a[(r0+i)*cols+c0+i]
	}
	return dst, nil
}
