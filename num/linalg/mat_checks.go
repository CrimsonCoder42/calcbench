package linalg

import (
	"errors"
	"fmt"
)

// checkTol returns the default tolerance n·2⁻⁵² for tol < 0, else tol.
func checkTol(tol float64, n int) float64 {
	if tol < 0 {
		return float64(max(n, 1)) * 0x1p-52
	}
	return tol
}

// MatIsSymmetric reports whether an n×n matrix is symmetric: |aᵢⱼ − aⱼᵢ| ≤ tol·max|a| for every pair
// (tol < 0 selects n·2⁻⁵²; tol = 0 asks for exact symmetry).
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-136
// Time:    O(n²)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.4 "Symmetric Matrices" (Aᵀ = A)
// Tags:    matrix, n-dimensional, symmetric, check, predicate, linear-algebra
// Units:   a in any unit; the result is a truth value
// Example: [1 2; 2 3] → true; [1 2; 3 4] → false (hand-checked)
func MatIsSymmetric(a []float64, n int, tol float64) (bool, error) {
	if err := checkSquare(a, n); err != nil {
		return false, err
	}
	big, _ := maxAbs(a)
	eps := checkTol(tol, n) * big
	for i := range n {
		for j := i + 1; j < n; j++ {
			if abs(a[i*n+j]-a[j*n+i]) > eps {
				return false, nil
			}
		}
	}
	return true, nil
}

// matOffPattern reports whether every entry of a rows×cols matrix where off(i, j) is true has
// |aᵢⱼ| ≤ tol·max|a|.
func matOffPattern(a []float64, rows, cols int, tol float64, off func(i, j int) bool) (bool, error) {
	if rows < 1 || cols < 1 {
		return false, fmt.Errorf("%w: a is %d×%d", ErrDimension, rows, cols)
	}
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return false, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return false, notFiniteAt("a", i)
	}
	big, _ := maxAbs(a)
	eps := checkTol(tol, max(rows, cols)) * big
	for i := range rows {
		for j := range cols {
			if off(i, j) && abs(a[i*cols+j]) > eps {
				return false, nil
			}
		}
	}
	return true, nil
}

// MatIsDiagonal reports whether every off-diagonal entry of a rows×cols matrix has |aᵢⱼ| ≤ tol·max|a|
// (tol < 0 selects max(rows, cols)·2⁻⁵²; tol = 0 asks for exact zeros).
// Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-137
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.4 "Symmetric Matrices" (the diagonal matrix Λ) and §2.6 "Elimination = Factorization: A = LU" (diagonal D in LDU)
// Tags:    matrix, n-dimensional, diagonal, check, predicate, linear-algebra
// Units:   a in any unit; the result is a truth value
// Example: [2 0; 0 3] → true; [2 1e−20; 0 3] with tol 0 → false (hand-checked)
func MatIsDiagonal(a []float64, rows, cols int, tol float64) (bool, error) {
	return matOffPattern(a, rows, cols, tol, func(i, j int) bool { return i != j })
}

// MatIsUpperTriangular reports whether every entry below the diagonal of a rows×cols matrix has
// |aᵢⱼ| ≤ tol·max|a| (tol < 0 selects max(rows, cols)·2⁻⁵²).
// Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-138
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (upper triangular: uᵢⱼ = 0 for i > j)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (upper triangular U)
// Tags:    matrix, n-dimensional, upper-triangular, check, predicate, linear-algebra
// Units:   a in any unit; the result is a truth value
// Example: [1 2; 0 3] → true; [1 2; 4 3] → false (hand-checked)
func MatIsUpperTriangular(a []float64, rows, cols int, tol float64) (bool, error) {
	return matOffPattern(a, rows, cols, tol, func(i, j int) bool { return i > j })
}

// MatIsLowerTriangular reports whether every entry above the diagonal of a rows×cols matrix has
// |aᵢⱼ| ≤ tol·max|a| (tol < 0 selects max(rows, cols)·2⁻⁵²).
// Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-139
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" (lower triangular: lᵢⱼ = 0 for i < j)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU" (lower triangular L)
// Tags:    matrix, n-dimensional, lower-triangular, check, predicate, linear-algebra
// Units:   a in any unit; the result is a truth value
// Example: [1 0; 2 3] → true; [1 5; 2 3] → false (hand-checked)
func MatIsLowerTriangular(a []float64, rows, cols int, tol float64) (bool, error) {
	return matOffPattern(a, rows, cols, tol, func(i, j int) bool { return i < j })
}

// MatIsOrthogonal reports whether an n×n matrix Q is orthogonal: ‖QᵀQ − I‖_max ≤ tol (absolute, since QᵀQ
// is dimensionless; tol < 0 selects n·2⁻⁵²·8 to allow a few roundings per entry).
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-140
// Time:    O(n³)   the product QᵀQ
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.4 "Orthonormal Bases and Gram-Schmidt" (QᵀQ = I; a square Q is orthogonal)
// Tags:    matrix, n-dimensional, orthogonal, check, predicate, linear-algebra
// Units:   a dimensionless; the result is a truth value
// Example: the rotation [0 −1; 1 0] → true; [1 1; 0 1] → false (hand-checked)
func MatIsOrthogonal(a []float64, n int, tol float64) (bool, error) {
	if err := checkSquare(a, n); err != nil {
		return false, err
	}
	if tol < 0 {
		tol = 8 * float64(n) * 0x1p-52
	}
	for i := range n {
		for j := i; j < n; j++ {
			var d dot2
			for k := range n {
				d.add(a[k*n+i], a[k*n+j])
			}
			want := 0.0
			if i == j {
				want = 1
			}
			if abs(d.result()-want) > tol {
				return false, nil
			}
		}
	}
	return true, nil
}

// MatIsPositiveDefinite reports whether a symmetric n×n matrix is positive definite, by attempting its
// Cholesky factorisation (it succeeds exactly when every pivot is > 0); a matrix that is not symmetric to
// tol (MatIsSymmetric) is reported false.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-141
// Time:    O(n³)   n³/3 flops (Cholesky)
// Space:   O(n²) extra (the factor)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (positive pivots ⇔ positive eigenvalues ⇔ xᵀAx > 0)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §4.2 "Positive Definite Systems" (Cholesky exists exactly for symmetric positive definite A)
// Tags:    matrix, n-dimensional, positive-definite, cholesky, check, predicate, linear-algebra
// Units:   a in any unit; the result is a truth value
// Example: [2 −1; −1 2] → true; [1 2; 2 1] → false (eigenvalues 3, −1) (hand-checked)
func MatIsPositiveDefinite(a []float64, n int, tol float64) (bool, error) {
	sym, err := MatIsSymmetric(a, n, tol)
	if err != nil || !sym {
		return false, err
	}
	_, err = Cholesky(nil, a, n)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, ErrNotPositiveDefinite):
		return false, nil
	}
	return false, err
}

// MatIsPositiveSemidefinite reports whether a symmetric n×n matrix is positive semidefinite: its smallest
// eigenvalue (MatSymmetricEigenJacobi, run with the same tol for its symmetry test and its stopping test) is
// ≥ −tol·max|λ| (tol < 0 selects n·2⁻⁵²·8 here and Jacobi's own default there). Cholesky cannot decide
// this, since a semidefinite matrix has a zero pivot. A matrix not symmetric to tol is reported false.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrNoConvergence
// from the eigenvalue sweep.
//
// ID:      NUM-LINALG-142
// Time:    O(n³) per Jacobi sweep
// Space:   O(n²) extra (the working copy)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.5 "Positive Definite Matrices" (semidefinite: all eigenvalues ≥ 0)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.5 "Jacobi Methods"
// Tags:    matrix, n-dimensional, positive-semidefinite, eigenvalues, check, predicate, linear-algebra
// Units:   a in any unit; the result is a truth value
// Example: [1 1; 1 1] → true (eigenvalues 0, 2); [1 2; 2 1] → false (hand-checked)
func MatIsPositiveSemidefinite(a []float64, n int, tol float64) (bool, error) {
	sym, err := MatIsSymmetric(a, n, tol)
	if err != nil || !sym {
		return false, err
	}
	vals, _, err := MatSymmetricEigenJacobi(nil, nil, a, n, tol, false)
	if err != nil {
		return false, err
	}
	if tol < 0 {
		tol = 8 * float64(n) * 0x1p-52
	}
	scale := max(abs(vals[0]), abs(vals[n-1]))
	return vals[0] >= -tol*scale, nil
}

// MatIsSingular reports whether an n×n matrix is singular to working precision: its numerical rank
// (MatRank with tol; tol < 0 selects MatRank's default) is below n.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-143
// Time:    O(n³)   complete-pivoting elimination
// Space:   O(n²) extra (MatRank's working copy)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.4 "Independence, Basis and Dimension" (singular: dependent columns, rank < n)
// Tags:    matrix, n-dimensional, singular, rank, check, predicate, linear-algebra
// Units:   a in any unit; the result is a truth value
// Example: [1 2; 2 4] → true; [1 2; 3 4] → false (hand-checked)
func MatIsSingular(a []float64, n int, tol float64) (bool, error) {
	if err := checkSquare(a, n); err != nil {
		return false, err
	}
	r, err := MatRank(a, n, n, tol)
	if err != nil {
		return false, err
	}
	return r < n, nil
}
