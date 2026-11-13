package linalg

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrNotPositiveDefinite is returned when a matrix that must be symmetric positive definite has a
// non-positive pivot in its Cholesky factorisation.
var ErrNotPositiveDefinite = errors.New("linalg: matrix is not positive definite")

// Cholesky returns the lower-triangular L with positive diagonal such that A = L·Lᵀ, for a symmetric
// positive-definite n×n matrix A, in dst (which must not overlap a). Only the lower triangle of a is read;
// symmetry is not checked (MatIsSymmetric does that). Entries of dst above the diagonal are set to 0.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry;
// ErrNotPositiveDefinite when a pivot is ≤ 0 (A is not positive definite, to working precision).
//
// ID:      NUM-LINALG-152
// Time:    O(n³)   n³/3 flops
// Space:   O(1) extra; O(n²) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §4.2 "Positive Definite Systems" (the Cholesky factorization A = GGᵀ, existence for symmetric positive definite A, the gaxpy Cholesky algorithm)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (A = LDLᵀ with positive pivots; A = (L√D)(L√D)ᵀ is Cholesky)
// Tags:    matrix, n-dimensional, cholesky, factorization, positive-definite, linear-algebra
// Units:   a in unit A; L in √A
// Example: A = [4 2; 2 3] → L = [2 0; 1 √2] (hand-checked: L·Lᵀ = [4 2; 2 3])
func Cholesky(dst, a []float64, n int) ([]float64, error) {
	if n < 1 {
		return nil, fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if err := checkMatrix("a", a, n, n); err != nil {
		return nil, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, notFiniteAt("a", i)
	}
	dst = grow(dst, n*n)
	clear(dst)
	for j := range n {
		s := a[j*n+j]
		for k := range j {
			s -= dst[j*n+k] * dst[j*n+k]
		}
		if !(s > 0) {
			return nil, fmt.Errorf("%w: pivot %d = %g", ErrNotPositiveDefinite, j, s)
		}
		d := math.Sqrt(s)
		dst[j*n+j] = d
		for i := j + 1; i < n; i++ {
			t := a[i*n+j]
			for k := range j {
				t -= dst[i*n+k] * dst[j*n+k]
			}
			dst[i*n+j] = t / d
		}
	}
	if i := firstNotFinite(dst); i >= 0 {
		return nil, fmt.Errorf("%w: L element %d", ErrOverflow, i)
	}
	return dst, nil
}

// CholeskySolve solves Ax = b given the Cholesky factor L of A (A = L·Lᵀ, as returned by Cholesky) by a
// forward solve Ly = b and a back solve Lᵀx = y, in dst (dst may be b).
// Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrSingular when a
// diagonal entry of L is zero; ErrOverflow when a component leaves the float64 range.
//
// ID:      NUM-LINALG-153
// Time:    O(n²)   2n² flops
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §4.2 "Positive Definite Systems" (solving Ax = b via Gy = b, Gᵀx = y) and §3.1 "Triangular Systems"
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (Cholesky) and §2.6 "Elimination = Factorization: A = LU" (two triangular solves)
// Tags:    matrix, n-dimensional, cholesky, solve, positive-definite, linear-algebra
// Units:   L in √A, b in unit B; x in B/A
// Example: L = [2 0; 1 √2] (A = [4 2; 2 3]), b = (8, 7) → x = (1.25, 1.5) (hand-checked: 4·1.25 + 2·1.5 = 8, 2·1.25 + 3·1.5 = 7)
func CholeskySolve(dst, l []float64, n int, b []float64) ([]float64, error) {
	y, err := ForwardSubstitution(dst, l, n, b)
	if err != nil {
		return nil, err
	}
	for i := n - 1; i >= 0; i-- {
		s := y[i]
		for j := i + 1; j < n; j++ {
			s -= l[j*n+i] * y[j]
		}
		y[i] = s / l[i*n+i]
		if !reduce.AllFinite(y[i]) {
			return nil, fmt.Errorf("%w: x[%d]", ErrOverflow, i)
		}
	}
	return y, nil
}
