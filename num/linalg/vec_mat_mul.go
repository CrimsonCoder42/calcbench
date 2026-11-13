package linalg

import "fmt"

// VecMatMul returns the row vector xᵀ·A (numpy's vecmat) for a rows×cols matrix and len(x) = rows, length
// cols, in dst (reused when its capacity suffices; it must not overlap x or a). A is read row by row
// (contiguous), so no transpose is formed.
// Errors: ErrDimension when len(x) ≠ rows or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry;
// ErrOverflow when an entry of the result leaves the float64 range.
//
// ID:      NUM-LINALG-185
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1) extra; O(c) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §1.1 "Basic Algorithms and Notation" (vector-matrix products, row-oriented saxpy form)
// Tags:    vector, matrix, n-dimensional, vecmat, row-vector, product, linear-algebra
// Units:   x in X, a in A; the result in X·A
// Example: (1, 2)·[1 2 3; 4 5 6] → (9, 12, 15) (hand-checked)
func VecMatMul(dst, x, a []float64, rows, cols int) ([]float64, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return nil, err
	}
	if len(x) != rows {
		return nil, lengthError("x", len(x), "rows", rows)
	}
	if i := firstNotFinite(x); i >= 0 {
		return nil, notFiniteAt("x", i)
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, notFiniteAt("a", i)
	}
	dst = grow(dst, cols)
	clear(dst)
	for i, xi := range x {
		row := a[i*cols : i*cols+cols]
		for j, v := range row {
			dst[j] += xi * v
		}
	}
	if i := firstNotFinite(dst); i >= 0 {
		return nil, fmt.Errorf("%w: element %d", ErrOverflow, i)
	}
	return dst, nil
}
