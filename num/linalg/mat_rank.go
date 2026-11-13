package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MatRank returns the numerical rank of the rows×cols matrix a: the number of pivots larger than tol in
// Gaussian elimination with complete pivoting (at each step the largest |entry| of the remaining block
// is brought to the pivot position). Elimination stops when no remaining entry exceeds tol.
// tol < 0 selects the default max(rows, cols)·ε·max|aᵢⱼ| (ε = 2⁻⁵², machine epsilon): elimination's
// backward error is of that order, so pivots below it cannot be told from zero. tol = 0 counts every
// nonzero pivot. The input is not modified.
// The work copy is scaled so its largest |element| is in [1, 2), so elimination cannot overflow.
// Errors: ErrDimension when len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf element or tol.
//
// ID:      NUM-LINALG-099
// Time:    O(rows·cols·min(rows, cols))
// Space:   O(rows·cols) work (one allocation)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.4 "Pivoting" (complete pivoting)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (rank = number of pivots)
// Tags:    matrix, n-dimensional, rank, gaussian-elimination, complete-pivoting, linear-algebra
// Units:   tol in the unit of a
// Example: [1 2 3; 2 4 6; 1 0 1] → 2 (the second row is twice the first; hand-checked)
func MatRank(a []float64, rows, cols int, tol float64) (int, error) {
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return 0, err
	}
	if !reduce.AllFinite(tol) {
		return 0, errNotFinite("tol")
	}
	if i := firstNotFinite(a); i >= 0 {
		return 0, fmt.Errorf("%w: a[%d]", ErrNotFinite, i)
	}
	if rows == 0 || cols == 0 {
		return 0, nil
	}
	mx := 0.0
	for _, x := range a {
		mx = max(mx, math.Abs(x))
	}
	if mx == 0 {
		return 0, nil
	}
	if tol < 0 {
		tol = float64(max(rows, cols)) * 0x1p-52 * mx
	}
	// Work on a·2^(−e) (largest |element| in [1, 2)) so elimination cannot overflow near the float64
	// maximum; the rank does not change, and tol is scaled the same way.
	_, e := pow2Scale(mx)
	w := make([]float64, len(a))
	for i, x := range a {
		w[i] = math.Ldexp(x, -e)
	}
	tol = math.Ldexp(tol, -e)
	rank := 0
	for k := range min(rows, cols) {
		p, q, best := k, k, -1.0
		for i := k; i < rows; i++ {
			for j := k; j < cols; j++ {
				if v := math.Abs(w[i*cols+j]); v > best {
					p, q, best = i, j, v
				}
			}
		}
		if best <= tol {
			break
		}
		if p != k {
			for j := range cols {
				w[k*cols+j], w[p*cols+j] = w[p*cols+j], w[k*cols+j]
			}
		}
		if q != k {
			for i := range rows {
				w[i*cols+k], w[i*cols+q] = w[i*cols+q], w[i*cols+k]
			}
		}
		rank++
		pivot := w[k*cols+k]
		for i := k + 1; i < rows; i++ {
			l := w[i*cols+k] / pivot
			if l == 0 {
				continue
			}
			for j := k + 1; j < cols; j++ {
				w[i*cols+j] = math.FMA(-l, w[k*cols+j], w[i*cols+j])
			}
			w[i*cols+k] = 0
		}
	}
	return rank, nil
}
