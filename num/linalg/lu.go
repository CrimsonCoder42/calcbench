package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LU is the factorisation PA = LU of an n×n matrix by Gaussian elimination with partial pivoting: L is
// unit lower triangular, U upper triangular, P a row permutation. Build it with LUDecompose and reuse it
// for many right-hand sides with LUSolve. Its fields are unexported; it is safe to copy (the copies
// share the read-only factor storage).
type LU struct {
	n        int
	lu       []float64 // L below the diagonal (unit diagonal implied), U on and above it, row-major
	piv      []int     // row i of PA is row piv[i] of A
	sign     float64   // det P = ±1
	singular bool      // a pivot was exactly zero
}

// LUDecompose returns the factorisation PA = LU of the n×n matrix a by Gaussian elimination with
// partial pivoting: at step k the row with the largest |entry| in column k is swapped up, so every
// multiplier has |l| ≤ 1. The input is not modified. A singular matrix still factors (some pivot is
// exactly 0); LUSolve then reports ErrSingular and MatDet returns 0.
// Errors: ErrDimension when n < 1 or len(a) ≠ n·n; ErrNotFinite for a NaN or ±Inf element; ErrOverflow
// when an element grows beyond the float64 range during elimination.
//
// ID:      NUM-LINALG-094
// Time:    O(n³)   n = order of the matrix; about n³/3 multiply-adds (2n³/3 flops)
// Space:   O(n²) for the factor (one allocation of n² floats and one of n ints)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" and §3.4 "Pivoting" (partial pivoting)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU"
// Tags:    matrix, n-dimensional, lu, factorization, gaussian-elimination, partial-pivoting, linear-algebra
// Units:   the unit of a
// Example: [0 1; 2 3] → rows swapped (det P = −1), L = [1 0; 0 1], U = [2 3; 0 1] (hand-checked)
func LUDecompose(a []float64, n int) (LU, error) {
	if n < 1 {
		return LU{}, fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if err := checkMatrix("a", a, n, n); err != nil {
		return LU{}, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return LU{}, fmt.Errorf("%w: a[%d]", ErrNotFinite, i)
	}
	f := LU{n: n, lu: make([]float64, n*n), piv: make([]int, n), sign: 1}
	copy(f.lu, a)
	for i := range f.piv {
		f.piv[i] = i
	}
	w := f.lu
	for k := range n {
		p, best := k, math.Abs(w[k*n+k])
		for i := k + 1; i < n; i++ {
			if v := math.Abs(w[i*n+k]); v > best {
				p, best = i, v
			}
		}
		if best == 0 {
			f.singular = true
			continue
		}
		if p != k {
			rk, rp := w[k*n:k*n+n], w[p*n:p*n+n]
			for j := range rk {
				rk[j], rp[j] = rp[j], rk[j]
			}
			f.piv[k], f.piv[p] = f.piv[p], f.piv[k]
			f.sign = -f.sign
		}
		pivot := w[k*n+k]
		rk := w[k*n : k*n+n]
		for i := k + 1; i < n; i++ {
			ri := w[i*n : i*n+n]
			l := ri[k] / pivot
			ri[k] = l
			if l == 0 {
				continue
			}
			for j := k + 1; j < n; j++ {
				ri[j] = math.FMA(-l, rk[j], ri[j])
			}
		}
	}
	if i := firstNotFinite(w); i >= 0 {
		return LU{}, fmt.Errorf("%w: element growth during elimination", ErrOverflow)
	}
	return f, nil
}

// LUSolve returns the solution x of Ax = b from the factorisation f = LUDecompose(A): x = U⁻¹L⁻¹Pb by
// forward and back substitution, written into dst (which must not overlap b).
// Errors: ErrSingular when A is singular; ErrDimension when f is empty or len(b) ≠ n; ErrNotFinite for a
// NaN or ±Inf element of b; ErrOverflow when an element of x is outside the float64 range.
//
// ID:      NUM-LINALG-095
// Time:    O(n²)   n = order of the matrix
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.1 "Triangular Systems" and §3.4 "Pivoting"
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.6 "Elimination = Factorization: A = LU"
// Tags:    matrix, n-dimensional, lu, solve, linear-system, forward-substitution, back-substitution, linear-algebra
// Units:   b in unit B with A in unit A; x in B/A
// Example: A = [2 1; 1 3], b = (3, 5) → x = (0.8, 1.4) (hand-checked)
func LUSolve(dst []float64, f LU, b []float64) ([]float64, error) {
	if f.n == 0 {
		return nil, fmt.Errorf("%w: empty factorisation", ErrDimension)
	}
	if len(b) != f.n {
		return nil, fmt.Errorf("%w: len(b) = %d, want %d", ErrDimension, len(b), f.n)
	}
	if f.singular {
		return nil, fmt.Errorf("%w: zero pivot", ErrSingular)
	}
	if i := firstNotFinite(b); i >= 0 {
		return nil, fmt.Errorf("%w: b[%d]", ErrNotFinite, i)
	}
	dst = grow(dst, f.n)
	for i := range f.n {
		dst[i] = b[f.piv[i]]
	}
	f.solveInPlace(dst)
	if i := firstNotFinite(dst); i >= 0 {
		return nil, fmt.Errorf("%w: x[%d]", ErrOverflow, i)
	}
	return dst, nil
}

// solveInPlace overwrites y (already permuted, y = Pb) with U⁻¹L⁻¹y.
func (f *LU) solveInPlace(y []float64) {
	n, w := f.n, f.lu
	for i := 1; i < n; i++ {
		s := y[i]
		for j := range i {
			s = math.FMA(-w[i*n+j], y[j], s)
		}
		y[i] = s
	}
	for i := n - 1; i >= 0; i-- {
		s := y[i]
		for j := i + 1; j < n; j++ {
			s = math.FMA(-w[i*n+j], y[j], s)
		}
		y[i] = s / w[i*n+i]
	}
}

// solveTransposeInPlace overwrites y with z solving Aᵀz = y: Aᵀ = UᵀLᵀP, so Uᵀu = y (forward), Lᵀv = u
// (back), and z = Pᵀv.
func (f *LU) solveTransposeInPlace(y, tmp []float64) {
	n, w := f.n, f.lu
	for i := range n {
		s := y[i]
		for j := range i {
			s = math.FMA(-w[j*n+i], y[j], s)
		}
		y[i] = s / w[i*n+i]
	}
	for i := n - 2; i >= 0; i-- {
		s := y[i]
		for j := i + 1; j < n; j++ {
			s = math.FMA(-w[j*n+i], y[j], s)
		}
		y[i] = s
	}
	for i := range n {
		tmp[f.piv[i]] = y[i]
	}
	copy(y, tmp)
}

// MatSolve returns the solution x of Ax = b for the n×n matrix a, by LU with partial pivoting, written
// into dst (which must not overlap b). To solve for several right-hand sides, factor once with
// LUDecompose and call LUSolve. For an ill-conditioned A the answer is inaccurate; check with
// MatConditionEstimate1.
// Errors: as LUDecompose and LUSolve.
//
// ID:      NUM-LINALG-096
// Time:    O(n³)   n = order of the matrix
// Space:   O(n²) for the factor; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" and §3.4 "Pivoting"
// Tags:    matrix, n-dimensional, solve, linear-system, lu, gaussian-elimination, linear-algebra
// Units:   b in unit B with A in unit A; x in B/A
// Example: A = [2 1; 1 3], b = (3, 5) → x = (0.8, 1.4) (hand-checked)
func MatSolve(dst, a []float64, n int, b []float64) ([]float64, error) {
	f, err := LUDecompose(a, n)
	if err != nil {
		return nil, err
	}
	return LUSolve(dst, f, b)
}

// MatDet returns the determinant of the n×n matrix a as det P · Π uᵢᵢ from LU with partial pivoting. The
// product is kept as a mantissa and a separate binary exponent, so it does not overflow or underflow
// part-way. A singular matrix has determinant 0 (not an error).
// Errors: ErrDimension when n < 1 or len(a) ≠ n·n; ErrNotFinite for a NaN or ±Inf element; ErrOverflow
// when |det| exceeds the float64 range.
//
// ID:      NUM-LINALG-097
// Time:    O(n³)   n = order of the matrix
// Space:   O(n²) for the factor
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" (det A = det P · u₁₁⋯uₙₙ)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.1 "The Properties of Determinants" (the determinant is ± the product of the pivots)
// Tags:    matrix, n-dimensional, determinant, lu, linear-algebra
// Units:   a in unit A; the result in Aⁿ
// Example: [1 2 3; 0 1 4; 5 6 0] → 1 (hand-checked)
func MatDet(a []float64, n int) (float64, error) {
	f, err := LUDecompose(a, n)
	if err != nil {
		return 0, err
	}
	if f.singular {
		return 0, nil
	}
	frac, exp := f.sign, 0
	for i := range n {
		m, e := math.Frexp(f.lu[i*n+i])
		frac *= m
		exp += e
		frac, e = math.Frexp(frac)
		exp += e
	}
	d := math.Ldexp(frac, exp)
	if !reduce.AllFinite(d) {
		return 0, errOverflow("det a")
	}
	return d, nil
}

// MatInverse returns the inverse of the n×n matrix a, column by column from LU with partial pivoting,
// in dst (which must not overlap a). Forming an inverse is rarely needed to solve Ax = b (use MatSolve);
// for an ill-conditioned A the inverse is inaccurate, so check with MatConditionEstimate1.
// Errors: ErrSingular when a is singular; otherwise as LUDecompose; ErrOverflow when an element of the
// inverse is outside the float64 range.
//
// ID:      NUM-LINALG-098
// Time:    O(n³)   n = order of the matrix
// Space:   O(n²) for the factor and O(n) work; O(n²) for the result, one allocation only if dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §3.2 "The LU Factorization" and §3.4 "Pivoting"
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (A⁻¹ column by column: Axⱼ = eⱼ)
// Tags:    matrix, n-dimensional, inverse, lu, linear-algebra
// Units:   a in unit A; the result in 1/A
// Example: [1 2 3; 0 1 4; 5 6 0] → [−24 18 5; 20 −15 −4; −5 4 1] (hand-checked)
func MatInverse(dst, a []float64, n int) ([]float64, error) {
	f, err := LUDecompose(a, n)
	if err != nil {
		return nil, err
	}
	if f.singular {
		return nil, fmt.Errorf("%w: zero pivot", ErrSingular)
	}
	col := make([]float64, n)
	dst = grow(dst, n*n)
	for j := range n {
		for i := range n {
			col[i] = 0
			if f.piv[i] == j {
				col[i] = 1
			}
		}
		f.solveInPlace(col)
		for i, x := range col {
			if !reduce.AllFinite(x) {
				return nil, fmt.Errorf("%w: element (%d, %d)", ErrOverflow, i, j)
			}
			dst[i*n+j] = x
		}
	}
	return dst, nil
}
