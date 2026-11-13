package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialFitCovariance returns the coefficients of the weighted least-squares polynomial fit
// (exactly as PolynomialFit, NUM-LINALG-171) together with their d×d covariance matrix
// Σ = σ²·(VᵀW²V)⁻¹, row-major in covDst, where V is the Vandermonde matrix, W = diag(weights) and
// σ² = Σ wᵢ²·(p(xᵢ) − yᵢ)² / (m − d) is the residual variance (numpy's polyfit(..., cov=True) scaling).
// (VᵀW²V)⁻¹ is formed from the triangular factor of the Householder QR of the column-scaled matrix
// WV·S, as S·R⁻¹R⁻ᵀ·S; the normal-equations matrix is never formed or inverted (its condition number is
// κ(WV)²). coeffDst and covDst may be nil.
// Errors: as PolynomialFit, and ErrDimension when m ≤ d (no degrees of freedom for σ²); ErrOverflow
// when σ² or a covariance entry leaves the float64 range.
//
// ID:      NUM-LINALG-190
// Time:    O(m·d²)   m = len(x), d = degree + 1 (two QR factorisations of the m×d matrix and one d×d triangular inverse)
// Space:   O(m·d) extra; O(d + d²) for the results
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §5.3 "The Full-Rank Least Squares Problem" (QR solution; (AᵀA)⁻¹ = R⁻¹R⁻ᵀ)
// Source:  Björck, Numerical Methods for Least Squares Problems, SIAM, 1996, §1.2 "Statistical properties" (covariance σ²(AᵀA)⁻¹ and the unbiased σ² = ‖r‖²/(m − n))
// Tags:    polynomial, fit, regression, least-squares, covariance, standard-error, uncertainty, qr, weighted, curve-fitting, linear-algebra
// Units:   x in X, y in Y; c[j] in Y/Xʲ; cov[j][k] in Y²/X^(j+k); weights in 1/Y (any consistent scale)
// Example: degree 1 through (0, 6), (1, 0), (2, 0) → c = (5, −3), σ² = 6, cov = [5 −3; −3 3] (hand-checked: (AᵀA)⁻¹ = [5/6 −1/2; −1/2 1/2])
func PolynomialFitCovariance(coeffDst, covDst, x, y, weights []float64, degree int, tol float64) (coeffs, cov []float64, err error) {
	c, err := PolynomialFit(coeffDst, x, y, weights, degree, tol)
	if err != nil {
		return nil, nil, err
	}
	m, d := len(x), degree+1
	if m <= d {
		return nil, nil, fmt.Errorf("%w: %d points for degree %d leave no degrees of freedom", ErrDimension, m, degree)
	}
	buf := make([]float64, m*d+d+d*d)
	v, colScale, rinv := buf[:m*d], buf[m*d:m*d+d], buf[m*d+d:]
	if _, err := Vandermonde(v, x, d); err != nil {
		return nil, nil, err
	}
	// Weighted residual sum of squares, scaled by a running power of two to avoid overflow.
	ssr, scale := 0.0, 0.0 // ssr·scale² = Σ rᵢ² (LAPACK dnrm2-style)
	for i := range m {
		w := 1.0
		if len(weights) != 0 {
			w = weights[i]
		}
		p := 0.0
		for j := d - 1; j >= 0; j-- {
			p = p*x[i] + c[j]
		}
		r := w * (p - y[i])
		if !reduce.AllFinite(r) {
			return nil, nil, fmt.Errorf("%w: residual %d", ErrOverflow, i)
		}
		if r != 0 {
			ar := math.Abs(r)
			if scale < ar {
				ssr = 1 + ssr*(scale/ar)*(scale/ar)
				scale = ar
			} else {
				ssr += (ar / scale) * (ar / scale)
			}
		}
		for j := range d {
			v[i*d+j] *= w
		}
	}
	sigma2 := scale * scale * ssr / float64(m-d)
	if !reduce.AllFinite(sigma2) {
		return nil, nil, errOverflow("residual variance")
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
	}
	householderQR(v, m, d) // R in the upper d×d triangle of v (row stride d)
	// R⁻¹ (upper triangular) by column-wise back substitution.
	for k := range d {
		rinv[k*d+k] = 1 / v[k*d+k]
		for i := k - 1; i >= 0; i-- {
			s := 0.0
			for j := i + 1; j <= k; j++ {
				s += v[i*d+j] * rinv[j*d+k]
			}
			rinv[i*d+k] = -s / v[i*d+i]
		}
	}
	if cap(covDst) >= d*d {
		covDst = covDst[:d*d]
	} else {
		covDst = make([]float64, d*d)
	}
	// Σ = σ²·S·(R⁻¹R⁻ᵀ)·S, symmetric: fill both halves from the upper one.
	for i := range d {
		for k := i; k < d; k++ {
			s := 0.0
			for j := k; j < d; j++ { // R⁻¹ upper: row i nonzero from i, row k from k
				s += rinv[i*d+j] * rinv[k*d+j]
			}
			e := sigma2 * colScale[i] * colScale[k] * s
			if !reduce.AllFinite(e) {
				return nil, nil, fmt.Errorf("%w: cov[%d][%d]", ErrOverflow, i, k)
			}
			covDst[i*d+k], covDst[k*d+i] = e, e
		}
	}
	return c, covDst, nil
}
