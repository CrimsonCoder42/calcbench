package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MatNormNuclear returns the nuclear (trace, Schatten-1) norm ‖A‖_* = Σ σᵢ of an m×n matrix, the sum of its
// singular values (SVDJacobi), added with Neumaier compensation.
// Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element;
// ErrNoConvergence from the SVD; ErrOverflow when the sum leaves the float64 range.
//
// ID:      NUM-LINALG-175
// Time:    O(m·n·min(m, n)) per SVD sweep
// Space:   O(m·n), the SVD's working copy
// Source:  Horn, Johnson, Matrix Analysis, 2nd ed., Cambridge University Press, 2013, §5.6 "Matrix norms" (Example 5.6.0.0 ff., the Schatten p-norms; the trace norm Σσᵢ) and §7.4.7
// Tags:    matrix, n-dimensional, norm, nuclear-norm, trace-norm, schatten, svd, linear-algebra
// Units:   a in unit A; the result in A
// Example: [3 0; 4 5] → 3√5 + √5 = 4√5 ≈ 8.944272 (hand-checked from σ = (3√5, √5))
func MatNormNuclear(a []float64, m, n int) (float64, error) {
	svd, err := SVDJacobi(a, m, n, -1, false)
	if err != nil {
		return 0, err
	}
	s, c := 0.0, 0.0
	for _, x := range svd.S {
		s, c = reduce.NeumaierAdd(s, c, x)
	}
	if !reduce.AllFinite(s + c) {
		return 0, errOverflow("nuclear norm")
	}
	return s + c, nil
}

// MatNormNeg1 returns minⱼ Σᵢ |aᵢⱼ|, the smallest absolute column sum of a rows×cols matrix (numpy's
// norm(ord=−1); not a norm). Each column sum is compensated (Neumaier).
// Errors: ErrDimension when rows or cols < 1 or len(a) disagrees; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when a column sum leaves the float64 range.
//
// ID:      NUM-LINALG-176
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (column sums of |A|, here the minimum)
// Tags:    matrix, n-dimensional, column-sum, minimum, norm-variant, linear-algebra
// Units:   a in unit A; the result in A
// Example: [1 −2; 3 4] → min(4, 6) = 4 (hand-checked)
func MatNormNeg1(a []float64, rows, cols int) (float64, error) {
	if rows < 1 || cols < 1 {
		return 0, fmt.Errorf("%w: a is %d×%d", ErrDimension, rows, cols)
	}
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return 0, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return 0, notFiniteAt("a", i)
	}
	best := 0.0
	for j := range cols {
		s, c := 0.0, 0.0
		for i := range rows {
			s, c = reduce.NeumaierAdd(s, c, abs(a[i*cols+j]))
		}
		t := s + c
		if !reduce.AllFinite(t) {
			return 0, fmt.Errorf("%w: column %d sum", ErrOverflow, j)
		}
		if j == 0 || t < best {
			best = t
		}
	}
	return best, nil
}

// MatNormNegInf returns minᵢ Σⱼ |aᵢⱼ|, the smallest absolute row sum of a rows×cols matrix (numpy's
// norm(ord=−inf); not a norm). Each row sum is compensated (Neumaier).
// Errors: as MatNormNeg1.
//
// ID:      NUM-LINALG-177
// Time:    O(r·c)   r = rows, c = cols
// Space:   O(1), no allocation
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (row sums of |A|, here the minimum)
// Tags:    matrix, n-dimensional, row-sum, minimum, norm-variant, linear-algebra
// Units:   a in unit A; the result in A
// Example: [1 −2; 3 4] → min(3, 7) = 3 (hand-checked)
func MatNormNegInf(a []float64, rows, cols int) (float64, error) {
	if rows < 1 || cols < 1 {
		return 0, fmt.Errorf("%w: a is %d×%d", ErrDimension, rows, cols)
	}
	if err := checkMatrix("a", a, rows, cols); err != nil {
		return 0, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return 0, notFiniteAt("a", i)
	}
	best := 0.0
	for i := range rows {
		s, c := 0.0, 0.0
		for _, x := range a[i*cols : i*cols+cols] {
			s, c = reduce.NeumaierAdd(s, c, abs(x))
		}
		t := s + c
		if !reduce.AllFinite(t) {
			return 0, fmt.Errorf("%w: row %d sum", ErrOverflow, i)
		}
		if i == 0 || t < best {
			best = t
		}
	}
	return best, nil
}

// MatNormNeg2 returns σ_min, the smallest of the min(m, n) singular values of an m×n matrix (numpy's
// norm(ord=−2); not a norm).
// Errors: as MatNorm2.
//
// ID:      NUM-LINALG-178
// Time:    O(m·n·min(m, n)) per SVD sweep
// Space:   O(m·n), the SVD's working copy
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.4 "The Singular Value Decomposition" (σ_min = min ‖Ax‖₂ over ‖x‖₂ = 1, Corollary 2.4.3)
// Tags:    matrix, n-dimensional, smallest-singular-value, svd, norm-variant, linear-algebra
// Units:   a in unit A; the result in A
// Example: [3 0; 4 5] → √5 ≈ 2.236068 (hand-checked)
func MatNormNeg2(a []float64, m, n int) (float64, error) {
	svd, err := SVDJacobi(a, m, n, -1, false)
	if err != nil {
		return 0, err
	}
	return svd.S[len(svd.S)-1], nil
}

// VecNorm0 returns the number of nonzero entries of v (numpy's vector norm(ord=0), the "ℓ₀ norm"; not a
// norm).
// Errors: ErrNotFinite for a NaN or ±Inf entry.
//
// ID:      NUM-LINALG-179
// Time:    O(n)   n = len(v)
// Space:   O(1)
// Source:  Horn, Johnson, Matrix Analysis, 2nd ed., Cambridge University Press, 2013, §5.2 (vector norms; ℓ₀ counts the nonzero entries and fails homogeneity, so it is not one)
// Tags:    vector, n-dimensional, count-nonzero, sparsity, l0, linear-algebra
// Units:   v in any unit; the result is a count
// Example: (0, 3, 0, −1) → 2 (hand-checked)
func VecNorm0(v []float64) (int, error) {
	n := 0
	for i, x := range v {
		if !reduce.AllFinite(x) {
			return 0, notFiniteAt("v", i)
		}
		if x != 0 {
			n++
		}
	}
	return n, nil
}
