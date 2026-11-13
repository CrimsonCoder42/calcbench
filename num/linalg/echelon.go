package linalg

import "fmt"

// echelon reduces w (rows×cols, in place) by Gaussian elimination with partial pivoting, treating
// |pivot| ≤ tol·max|a| as zero; when reduced is true each pivot is scaled to 1 and the entries above it
// are eliminated too (RREF). It returns the pivot columns in order.
func echelon(w []float64, rows, cols int, tol float64, reduced bool) []int {
	big, _ := maxAbs(w)
	eps := tol * big
	pivots := make([]int, 0, min(rows, cols))
	r := 0
	for c := 0; c < cols && r < rows; c++ {
		p, best := r, abs(w[r*cols+c])
		for i := r + 1; i < rows; i++ {
			if x := abs(w[i*cols+c]); x > best {
				p, best = i, x
			}
		}
		if best <= eps || best == 0 {
			for i := r; i < rows; i++ {
				w[i*cols+c] = 0 // below the tolerance: an exact zero, so the shape is clean
			}
			continue
		}
		if p != r {
			for j := range cols {
				w[r*cols+j], w[p*cols+j] = w[p*cols+j], w[r*cols+j]
			}
		}
		piv := w[r*cols+c]
		if reduced {
			for j := c; j < cols; j++ {
				w[r*cols+j] /= piv
			}
			w[r*cols+c] = 1
			piv = 1
		}
		start := r + 1
		if reduced {
			start = 0
		}
		for i := start; i < rows; i++ {
			if i == r {
				continue
			}
			f := w[i*cols+c] / piv
			if f == 0 {
				continue
			}
			for j := c; j < cols; j++ {
				w[i*cols+j] -= f * w[r*cols+j]
			}
			w[i*cols+c] = 0
		}
		pivots = append(pivots, c)
		r++
	}
	return pivots
}

// checkEchelonInput checks a rows×cols matrix with finite entries and returns the default tolerance.
func checkEchelonInput(a []float64, rows, cols int, tol float64) (float64, error) {
	if rows < 1 || cols < 1 {
		return 0, fmt.Errorf("%w: a is %d×%d", ErrDimension, rows, cols)
	}
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return 0, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return 0, notFiniteAt("a", i)
	}
	if tol < 0 {
		tol = float64(max(rows, cols)) * 0x1p-52
	}
	return tol, nil
}

// MatRowEchelon returns a row echelon form of a rows×cols matrix by Gaussian elimination with partial
// pivoting (the largest |entry| in the column becomes the pivot), in dst, with the pivot columns. An entry
// with |x| ≤ tol·max|aᵢⱼ| is treated as zero (tol < 0 selects max(rows, cols)·2⁻⁵²) and set to exactly 0, so
// the staircase shape is exact. The echelon form is not unique (it depends on the pivoting); the pivot
// columns and their number (the rank) are. dst must not overlap a.
// Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-130
// Time:    O(r·c·min(r, c))   r = rows, c = cols
// Space:   O(min(r, c)) extra for the pivot list; O(r·c) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.2 "The Idea of Elimination" and §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (echelon form U, pivot columns)
// Tags:    matrix, n-dimensional, row-echelon, elimination, pivoting, linear-algebra
// Units:   a in unit A; the result in A
// Example: [1 2 1; 2 4 0; 3 6 3] → pivots in columns 0 and 2 (rank 2) (hand-checked)
func MatRowEchelon(dst, a []float64, rows, cols int, tol float64) ([]float64, []int, error) {
	tol, err := checkEchelonInput(a, rows, cols, tol)
	if err != nil {
		return nil, nil, err
	}
	dst = grow(dst, rows*cols)
	copy(dst, a)
	return dst, echelon(dst, rows, cols, tol, false), nil
}

// MatRREF returns the reduced row echelon form R of a rows×cols matrix (each pivot 1, zeros above and below
// it) by Gauss–Jordan elimination with partial pivoting, in dst, with the pivot columns. |x| ≤ tol·max|aᵢⱼ|
// is treated as zero (tol < 0 selects max(rows, cols)·2⁻⁵²). R is unique in exact arithmetic; in floating
// point the decision "zero or pivot" is the tolerance's. The RREF of an invertible matrix is I. dst must not
// overlap a.
// Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-131
// Time:    O(r·c·min(r, c))   r = rows, c = cols
// Space:   O(min(r, c)) extra; O(r·c) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (the reduced row echelon form R = rref(A)) and §3.3 "The Complete Solution to Ax = b"
// Tags:    matrix, n-dimensional, rref, reduced-row-echelon, gauss-jordan, pivoting, linear-algebra
// Units:   a in unit A; the result dimensionless (pivot rows are scaled to 1)
// Example: [1 2 2 4; 3 8 6 16] → R = [1 0 2 0; 0 1 0 2], pivots 0 and 1 (Strang §3.2 example; hand-checked)
func MatRREF(dst, a []float64, rows, cols int, tol float64) ([]float64, []int, error) {
	tol, err := checkEchelonInput(a, rows, cols, tol)
	if err != nil {
		return nil, nil, err
	}
	dst = grow(dst, rows*cols)
	copy(dst, a)
	return dst, echelon(dst, rows, cols, tol, true), nil
}

// MatNullSpace returns a basis of the null space {x : A·x = 0} of a rows×cols matrix as the columns of a
// cols×k matrix in dst (k = cols − rank), with k: from R = rref(A), one special solution per free column f
// (x_f = 1, x at each pivot column p = −R[row of p, f], other free entries 0). The basis is not orthonormal;
// orthonormalise it with QRHouseholder if needed. tol as in MatRREF. dst must not overlap a; k = 0 returns
// a zero-length dst.
// Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-132
// Time:    O(r·c·min(r, c))   the RREF dominates
// Space:   O(r·c) extra (the RREF); O(c·k) for the result
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.2 "The Nullspace of A: Solving Ax = 0 and Rx = 0" (special solutions: one per free variable; the nullspace matrix N)
// Tags:    matrix, n-dimensional, null-space, kernel, rref, basis, linear-algebra
// Units:   a in unit A; the basis dimensionless
// Example: [1 2 2 4; 3 8 6 16] → basis columns (−2, 0, 1, 0) and (0, −2, 0, 1) (Strang §3.2 example; hand-checked)
func MatNullSpace(dst, a []float64, rows, cols int, tol float64) ([]float64, int, error) {
	r, pivots, err := MatRREF(nil, a, rows, cols, tol)
	if err != nil {
		return nil, 0, err
	}
	isPivot := make([]int, cols) // row index + 1 of a pivot column, 0 for free
	for i, p := range pivots {
		isPivot[p] = i + 1
	}
	k := cols - len(pivots)
	dst = grow(dst, cols*k)
	clear(dst)
	col := 0
	for f := range cols {
		if isPivot[f] != 0 {
			continue
		}
		dst[f*k+col] = 1
		for i, p := range pivots {
			dst[p*k+col] = -r[i*cols+f]
		}
		col++
	}
	return dst, k, nil
}

// MatColumnSpace returns a basis of the column space of a rows×cols matrix: the pivot columns of A itself
// (found by MatRREF's pivoting), as the columns of a rows×rank matrix in dst, with the rank. tol as in
// MatRREF. dst must not overlap a.
// Errors: ErrDimension for rows or cols < 1 or len(a) ≠ rows·cols; ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-133
// Time:    O(r·c·min(r, c))   the RREF dominates
// Space:   O(r·c) extra (the RREF); O(r·rank) for the result
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §3.5 "Dimensions of the Four Subspaces" (the pivot columns of A are a basis for its column space C(A)) and §3.4 "Independence, Basis and Dimension"
// Tags:    matrix, n-dimensional, column-space, range, basis, rref, linear-algebra
// Units:   a in unit A; the basis in A
// Example: [1 2 2 4; 3 8 6 16] → columns (1, 3) and (2, 8), rank 2 (Strang §3.2 example; hand-checked)
func MatColumnSpace(dst, a []float64, rows, cols int, tol float64) ([]float64, int, error) {
	_, pivots, err := MatRREF(nil, a, rows, cols, tol)
	if err != nil {
		return nil, 0, err
	}
	k := len(pivots)
	dst = grow(dst, rows*k)
	for i := range rows {
		for c, p := range pivots {
			dst[i*k+c] = a[i*cols+p]
		}
	}
	return dst, k, nil
}
