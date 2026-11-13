package linalg

import "fmt"

// LeastSquaresMinNorm solves min ‖A·x − b‖₂ for any m×n matrix, rank-deficient, underdetermined or
// overdetermined, returning the minimum-norm solution x = V·Σ⁺·Uᵀ·b from the SVD (SVDJacobi,
// NUM-LINALG-160), with singular values σᵢ ≤ tol·σ₁ treated as zero (tol < 0 selects max(m, n)·2⁻⁵²), the
// residual norm ‖A·x − b‖₂ and the numerical rank. It equals MatPseudoInverse(A)·b without forming A⁺.
// xDst receives x (length n; it must not overlap a or b).
// Errors: ErrDimension for m or n < 1, len(a) ≠ m·n or len(b) ≠ m; ErrNotFinite for a NaN or ±Inf entry;
// ErrNoConvergence from the SVD; ErrOverflow when a coefficient leaves the float64 range.
//
// ID:      NUM-LINALG-173
// Time:    O(m·n·min(m, n)) per SVD sweep, plus O(m·n) for the products
// Space:   O(m·n) for the SVD; O(n) for the result
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §5.5 "Rank-Deficient Least Squares" (§5.5.1, the minimum-norm solution x_LS = Σ (uᵢᵀb/σᵢ)·vᵢ over σᵢ > 0) and §5.4
// Tags:    matrix, n-dimensional, least-squares, minimum-norm, rank-deficient, underdetermined, svd, lstsq, linear-algebra
// Units:   a in unit A, b in unit B; x in B/A; the residual norm in B
// Example: a = [1 1] (1×2), b = (2) → x = (1, 1), residual 0, rank 1 (the shortest x with x₁ + x₂ = 2; hand-checked)
func LeastSquaresMinNorm(xDst, a []float64, m, n int, b []float64, tol float64) (LeastSquaresResult, int, error) {
	if len(b) != m {
		return LeastSquaresResult{}, 0, lengthError("b", len(b), "m", m)
	}
	if i := firstNotFinite(b); i >= 0 {
		return LeastSquaresResult{}, 0, notFiniteAt("b", i)
	}
	svd, err := SVDJacobi(a, m, n, -1, true)
	if err != nil {
		return LeastSquaresResult{}, 0, err
	}
	if tol < 0 {
		tol = float64(max(m, n)) * 0x1p-52
	}
	k := min(m, n)
	cut := tol * svd.S[0]
	xDst = grow(xDst, n)
	clear(xDst)
	rank := 0
	for s := range k {
		sigma := svd.S[s]
		if sigma <= cut || sigma == 0 {
			continue
		}
		rank++
		ub := 0.0
		for i := range m {
			ub += svd.U[i*k+s] * b[i]
		}
		c := ub / sigma
		for i := range n {
			xDst[i] += c * svd.V[i*k+s]
		}
	}
	if i := firstNotFinite(xDst); i >= 0 {
		return LeastSquaresResult{}, 0, fmt.Errorf("%w: x[%d]", ErrOverflow, i)
	}
	r := make([]float64, m)
	for i := range m {
		s := -b[i]
		for j := range n {
			s += a[i*n+j] * xDst[j]
		}
		r[i] = s
	}
	rn, err := VecNorm(r)
	if err != nil {
		return LeastSquaresResult{}, 0, fmt.Errorf("%w: residual", ErrOverflow)
	}
	return LeastSquaresResult{X: xDst, ResidualNorm: rn}, rank, nil
}
