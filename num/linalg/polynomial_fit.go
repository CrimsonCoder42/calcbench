package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialFit returns the ascending coefficients c[0…degree] of the polynomial that minimises the
// weighted squared error Σ wᵢ²·(p(xᵢ) − yᵢ)² (numpy's polyfit weighting: w multiplies the residual), in
// coeffDst. weights may be nil (all 1). The Vandermonde matrix (NUM-LINALG-170) is scaled row by row by
// wᵢ and column by column by 1/‖column‖₂ (exact powers of two, to tame its conditioning) and solved by
// Householder QR (LeastSquaresQR, NUM-LINALG-156); normal equations are never formed. tol is passed to
// LeastSquaresQR (tol < 0 selects its default).
// Errors: ErrDimension when degree < 0, len(y) ≠ len(x), len(weights) ∉ {0, len(x)}, or len(x) < degree+1;
// ErrNotFinite for a NaN or ±Inf input; ErrOutOfRange for a negative weight; ErrRankDeficient when fewer
// than degree+1 distinct x carry nonzero weight; ErrOverflow when a power or coefficient leaves the
// float64 range.
//
// ID:      NUM-LINALG-171
// Time:    O(m·d²)   m = len(x), d = degree + 1
// Space:   O(m·d) extra (the scaled Vandermonde matrix and the QR working copy); O(d) for the result
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §5.3 "The Full-Rank Least Squares Problem" (QR solution) and §4.6 (Vandermonde matrices)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.3 "Least Squares Approximations" (fitting a parabola)
// Tags:    polynomial, fit, regression, least-squares, vandermonde, qr, weighted, curve-fitting, linear-algebra
// Units:   x in X, y in Y; c[j] in Y/Xʲ; weights in 1/Y (any consistent scale)
// Example: degree 1 through (0, 6), (1, 0), (2, 0) → c = (5, −3) (Strang §4.3 worked example)
func PolynomialFit(coeffDst, x, y, weights []float64, degree int, tol float64) ([]float64, error) {
	if degree < 0 {
		return nil, fmt.Errorf("%w: degree = %d", ErrDimension, degree)
	}
	m, d := len(x), degree+1
	if len(y) != m {
		return nil, lengthError("x", m, "y", len(y))
	}
	if len(weights) != 0 && len(weights) != m {
		return nil, lengthError("x", m, "weights", len(weights))
	}
	if m < d {
		return nil, fmt.Errorf("%w: %d points for degree %d", ErrDimension, m, degree)
	}
	if i := firstNotFinite(y); i >= 0 {
		return nil, notFiniteAt("y", i)
	}
	for i, w := range weights {
		if !reduce.AllFinite(w) {
			return nil, notFiniteAt("weights", i)
		}
		if w < 0 {
			return nil, fmt.Errorf("%w: weights[%d] = %g", ErrOutOfRange, i, w)
		}
	}
	buf := make([]float64, m*d+m+d)
	v, rhs, colScale := buf[:m*d], buf[m*d:m*d+m], buf[m*d+m:]
	if _, err := Vandermonde(v, x, d); err != nil {
		return nil, err
	}
	for i := range m {
		w := 1.0
		if len(weights) != 0 {
			w = weights[i]
		}
		rhs[i] = w * y[i]
		if !reduce.AllFinite(rhs[i]) {
			return nil, fmt.Errorf("%w: weights[%d]·y[%d]", ErrOverflow, i, i)
		}
		for j := range d {
			v[i*d+j] *= w
		}
	}
	for j := range d {
		big := 0.0
		for i := range m {
			big = max(big, math.Abs(v[i*d+j]))
		}
		colScale[j] = 1
		if big > 0 {
			s, _ := pow2Scale(big)
			colScale[j] = s
			for i := range m {
				v[i*d+j] *= s
			}
		}
		if !reduce.AllFinite(colScale[j]) {
			return nil, errOverflow("column scale")
		}
	}
	res, err := LeastSquaresQR(coeffDst, v, m, d, rhs, tol)
	if err != nil {
		return nil, err
	}
	c := res.X
	for j := range d {
		c[j] *= colScale[j]
		if !reduce.AllFinite(c[j]) {
			return nil, fmt.Errorf("%w: coefficient %d", ErrOverflow, j)
		}
	}
	return c, nil
}
