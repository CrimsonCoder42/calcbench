package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkSquareSystem checks an n×n matrix and an n-vector with finite entries.
func checkSquareSystem(a []float64, n int, b []float64) error {
	if n < 1 {
		return fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if err := checkMatrix("a", a, n, n); err != nil {
		return err
	}
	if len(b) != n {
		return lengthError("b", len(b), "n", n)
	}
	if i := firstNotFinite(a); i >= 0 {
		return notFiniteAt("a", i)
	}
	if i := firstNotFinite(b); i >= 0 {
		return notFiniteAt("b", i)
	}
	return nil
}

// ForwardSubstitution solves Lx = b for an n×n lower-triangular L (entries above the diagonal are ignored)
// by row-oriented forward substitution, in dst (dst may be b).
// Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrSingular when a
// diagonal entry is zero; ErrOverflow when a component leaves the float64 range.
//
// ID:      NUM-LINALG-157
// Time:    O(n²)   n² flops
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (row-oriented forward substitution)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (forward elimination Lc = b)
// Tags:    matrix, n-dimensional, triangular, forward-substitution, solve, linear-algebra
// Units:   l in unit L, b in unit B; x in B/L
// Example: L = [2 0; 3 4], b = (2, 11) → x = (1, 2) (hand-checked)
func ForwardSubstitution(dst, l []float64, n int, b []float64) ([]float64, error) {
	if err := checkSquareSystem(l, n, b); err != nil {
		return nil, err
	}
	dst = grow(dst, n)
	copy(dst, b)
	for i := range n {
		s := dst[i]
		for j := range i {
			s -= l[i*n+j] * dst[j]
		}
		d := l[i*n+i]
		if d == 0 {
			return nil, fmt.Errorf("%w: l[%d,%d] = 0", ErrSingular, i, i)
		}
		dst[i] = s / d
		if !reduce.AllFinite(dst[i]) {
			return nil, fmt.Errorf("%w: x[%d]", ErrOverflow, i)
		}
	}
	return dst, nil
}

// BackSubstitution solves Ux = b for an n×n upper-triangular U (entries below the diagonal are ignored)
// by row-oriented back substitution, in dst (dst may be b).
// Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrSingular when a
// diagonal entry is zero; ErrOverflow when a component leaves the float64 range.
//
// ID:      NUM-LINALG-158
// Time:    O(n²)   n² flops
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (row-oriented back substitution)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.2 "The Idea of Elimination" (back substitution on the upper triangular system)
// Tags:    matrix, n-dimensional, triangular, back-substitution, solve, linear-algebra
// Units:   u in unit U, b in unit B; x in B/U
// Example: U = [2 1; 0 4], b = (4, 8) → x = (1, 2) (hand-checked)
func BackSubstitution(dst, u []float64, n int, b []float64) ([]float64, error) {
	if err := checkSquareSystem(u, n, b); err != nil {
		return nil, err
	}
	dst = grow(dst, n)
	copy(dst, b)
	for i := n - 1; i >= 0; i-- {
		s := dst[i]
		for j := i + 1; j < n; j++ {
			s -= u[i*n+j] * dst[j]
		}
		d := u[i*n+i]
		if d == 0 {
			return nil, fmt.Errorf("%w: u[%d,%d] = 0", ErrSingular, i, i)
		}
		dst[i] = s / d
		if !reduce.AllFinite(dst[i]) {
			return nil, fmt.Errorf("%w: x[%d]", ErrOverflow, i)
		}
	}
	return dst, nil
}

// TridiagonalSolve solves Tx = d for an n×n tridiagonal T given by its sub-diagonal sub (len n−1), diagonal
// diag (len n) and super-diagonal sup (len n−1), by the Thomas algorithm (Gaussian elimination without
// pivoting specialised to three bands), in dst (which must not overlap an input). Without pivoting it is
// stable for diagonally dominant or symmetric positive-definite T; for other T use MatSolve.
// Errors: ErrDimension for mismatched lengths or n < 1; ErrNotFinite for a NaN or ±Inf entry; ErrSingular
// when a pivot is zero; ErrOverflow when a component leaves the float64 range.
//
// ID:      NUM-LINALG-159
// Time:    O(n)   8n flops
// Space:   O(n) extra (the modified super-diagonal, one allocation); O(n) for the result
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (a tridiagonal matrix has bidiagonal L and U, O(n) work)
// Tags:    matrix, n-dimensional, tridiagonal, thomas-algorithm, banded, solve, linear-algebra
// Units:   T in unit T, d in unit D; x in D/T
// Example: sub = (1, 1), diag = (4, 4, 4), sup = (1, 1), d = (5, 6, 5) → x = (1, 1, 1) (hand-checked)
func TridiagonalSolve(dst, sub, diag, sup, d []float64) ([]float64, error) {
	n := len(diag)
	if n < 1 {
		return nil, fmt.Errorf("%w: len(diag) = 0", ErrDimension)
	}
	if len(sub) != n-1 || len(sup) != n-1 || len(d) != n {
		return nil, fmt.Errorf("%w: len(sub) = %d, len(sup) = %d, len(d) = %d, want %d, %d, %d", ErrDimension, len(sub), len(sup), len(d), n-1, n-1, n)
	}
	for _, in := range [...]struct {
		name string
		v    []float64
	}{{"sub", sub}, {"diag", diag}, {"sup", sup}, {"d", d}} {
		if i := firstNotFinite(in.v); i >= 0 {
			return nil, notFiniteAt(in.name, i)
		}
	}
	c := make([]float64, n)
	dst = grow(dst, n)
	if diag[0] == 0 {
		return nil, fmt.Errorf("%w: pivot 0 is zero", ErrSingular)
	}
	if n > 1 {
		c[0] = sup[0] / diag[0]
	}
	dst[0] = d[0] / diag[0]
	for i := 1; i < n; i++ {
		piv := diag[i] - sub[i-1]*c[i-1]
		if piv == 0 {
			return nil, fmt.Errorf("%w: pivot %d is zero", ErrSingular, i)
		}
		if i < n-1 {
			c[i] = sup[i] / piv
		}
		dst[i] = (d[i] - sub[i-1]*dst[i-1]) / piv
	}
	for i := n - 2; i >= 0; i-- {
		dst[i] -= c[i] * dst[i+1]
	}
	if i := firstNotFinite(dst); i >= 0 {
		return nil, fmt.Errorf("%w: x[%d]", ErrOverflow, i)
	}
	return dst, nil
}
