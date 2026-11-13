package linalg

import (
	"fmt"
	"math"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// SVDResult is the thin singular value decomposition A = U·diag(S)·Vᵀ of an m×n matrix, k = min(m, n):
// S holds the k singular values in descending order; U (m×k) and V (n×k), row-major, are nil unless
// requested. Columns belonging to nonzero singular values are orthonormal; for a zero singular value the
// corresponding column of U (of V when m < n, where the roles swap) is zero.
type SVDResult struct {
	S    []float64
	U    []float64
	V    []float64
	Rows int // m
	Cols int // n
}

// SVDJacobi returns the thin SVD of an m×n matrix by one-sided Jacobi (Hestenes): plane rotations applied
// to pairs of columns of a working copy of A (of Aᵀ when m < n) until every pair is orthogonal to within
// tol (|uᵢ·uⱼ| ≤ tol·‖uᵢ‖‖uⱼ‖; tol < 0 selects max(m, n)·2⁻⁵²); the column norms are the singular values.
// One-sided Jacobi computes the singular values to high relative accuracy (Demmel–Veselić 1992). U and V are
// formed only when wantVectors is true.
// Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element;
// ErrNoConvergence after 60 sweeps.
//
// ID:      NUM-LINALG-160
// Time:    O(m·n·min(m, n)) per sweep   typically 6–10 sweeps
// Space:   O(m·n + min(m, n)²) for the working copy, the rotation accumulator and the result; for m ≫ n a QR preconditioning step would cut the per-sweep cost to O(n³) (LATER)
// Source:  Hestenes, "Inversion of Matrices by Biorthogonalization and Related Results", Journal of the Society for Industrial and Applied Mathematics, 1958, 6(1):51–90, doi:10.1137/0106005
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.6 "Computing the SVD" (Jacobi SVD procedures)
// Source:  Demmel, Veselić, "Jacobi's Method is More Accurate than QR", SIAM Journal on Matrix Analysis and Applications, 1992, 13(4):1204–1245, doi:10.1137/0613074
// Tags:    matrix, n-dimensional, svd, singular-values, jacobi, hestenes, factorization, linear-algebra
// Units:   a in unit A; singular values in A; U and V dimensionless
// Example: [3 0; 4 5] → S = (3√5, √5) ≈ (6.708204, 2.236068) (hand-checked: σ₁σ₂ = |det| = 15, σ₁² + σ₂² = ‖A‖_F² = 50)
func SVDJacobi(a []float64, m, n int, tol float64, wantVectors bool) (SVDResult, error) {
	if m < 1 || n < 1 {
		return SVDResult{}, fmt.Errorf("%w: a is %d×%d", ErrDimension, m, n)
	}
	if err := checkMatrix("a", a, m, n); err != nil {
		return SVDResult{}, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return SVDResult{}, notFiniteAt("a", i)
	}
	if tol < 0 {
		tol = float64(max(m, n)) * 0x1p-52
	}
	// Work on B = A (p = m rows, q = n columns) or B = Aᵀ when m < n, so that p ≥ q.
	p, q, tr := m, n, false
	if m < n {
		p, q, tr = n, m, true
	}
	big, _ := maxAbs(a)
	_, e := pow2Scale(max(big, 0x1p-1000))
	b := make([]float64, p*q)
	for i := range m {
		for j := range n {
			x := math.Ldexp(a[i*n+j], -e)
			if tr {
				b[j*q+i] = x
			} else {
				b[i*q+j] = x
			}
		}
	}
	v, _ := MatIdentity(nil, q)
	for sweep := 0; ; sweep++ {
		rotated := false
		for i := 0; i < q-1; i++ {
			for j := i + 1; j < q; j++ {
				alpha, beta, gamma := 0.0, 0.0, 0.0
				for r := range p {
					bi, bj := b[r*q+i], b[r*q+j]
					alpha += bi * bi
					beta += bj * bj
					gamma += bi * bj
				}
				if gamma == 0 || abs(gamma) <= tol*math.Sqrt(alpha)*math.Sqrt(beta) {
					continue
				}
				rotated = true
				zeta := (beta - alpha) / (2 * gamma)
				t := 1 / (abs(zeta) + math.Sqrt(1+zeta*zeta))
				if zeta < 0 {
					t = -t
				}
				c := 1 / math.Sqrt(1+t*t)
				s := c * t
				for r := range p {
					bi, bj := b[r*q+i], b[r*q+j]
					b[r*q+i] = c*bi - s*bj
					b[r*q+j] = s*bi + c*bj
				}
				for r := range q {
					vi, vj := v[r*q+i], v[r*q+j]
					v[r*q+i] = c*vi - s*vj
					v[r*q+j] = s*vi + c*vj
				}
			}
		}
		if !rotated {
			break
		}
		if sweep == maxJacobiSweeps {
			return SVDResult{}, fmt.Errorf("%w: %d sweeps", ErrNoConvergence, sweep)
		}
	}
	sig := make([]float64, q)
	for j := range q {
		for r := range p {
			sig[j] = math.Hypot(sig[j], b[r*q+j])
		}
	}
	idx := make([]int, q)
	for i := range q {
		idx[i] = i
	}
	slices.SortStableFunc(idx, func(x, y int) int { return cmpFloat(sig[y], sig[x]) })
	res := SVDResult{S: make([]float64, q), Rows: m, Cols: n}
	for k, j := range idx {
		res.S[k] = math.Ldexp(sig[j], e)
	}
	if !wantVectors {
		return res, nil
	}
	// B's left vectors are b's normalised columns (p×q), its right vectors are v (q×q).
	left := make([]float64, p*q)
	right := make([]float64, q*q)
	for k, j := range idx {
		for r := range p {
			if sig[j] > 0 {
				left[r*q+k] = b[r*q+j] / sig[j]
			}
		}
		for r := range q {
			right[r*q+k] = v[r*q+j]
		}
	}
	if tr { // A = Bᵀ = (V_B) S (U_B)ᵀ
		res.U, res.V = right, left
	} else {
		res.U, res.V = left, right
	}
	return res, nil
}

// MatPseudoInverse returns the Moore–Penrose pseudo-inverse A⁺ = V·diag(1/σᵢ)·Uᵀ of an m×n matrix, an n×m
// matrix, from SVDJacobi; singular values σᵢ ≤ tol·σ₁ are treated as zero (tol < 0 selects
// max(m, n)·2⁻⁵²), which gives the minimum-norm least-squares solution x = A⁺b and A·A⁺·A = A. dst must not
// overlap a.
// Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element;
// ErrNoConvergence from the SVD; ErrOverflow when an element leaves the float64 range.
//
// ID:      NUM-LINALG-161
// Time:    O(m·n·min(m, n)) per SVD sweep, plus O(m·n·min(m, n)) to form A⁺
// Space:   O(m·n) for the SVD; O(m·n) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, ch. 7 "The Singular Value Decomposition (SVD)" (the pseudoinverse A⁺ = VΣ⁺Uᵀ)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.6 "Computing the SVD" (the SVD as the tool for rank-deficient problems and the pseudo-inverse)
// Tags:    matrix, n-dimensional, pseudo-inverse, moore-penrose, svd, least-squares, linear-algebra
// Units:   a in unit A; the result in 1/A
// Example: [1 2; 2 4] (rank 1) → A⁺ = A/25 = [0.04 0.08; 0.08 0.16] (hand-checked: A·A⁺·A = A)
func MatPseudoInverse(dst, a []float64, m, n int, tol float64) ([]float64, error) {
	svd, err := SVDJacobi(a, m, n, -1, true)
	if err != nil {
		return nil, err
	}
	if tol < 0 {
		tol = float64(max(m, n)) * 0x1p-52
	}
	k := min(m, n)
	cut := tol * svd.S[0]
	dst = grow(dst, n*m)
	clear(dst)
	for s := range k {
		sigma := svd.S[s]
		if sigma <= cut || sigma == 0 {
			continue
		}
		for i := range n {
			vi := svd.V[i*k+s] / sigma
			for j := range m {
				dst[i*m+j] += vi * svd.U[j*k+s]
			}
		}
	}
	if i := firstNotFinite(dst); i >= 0 {
		return nil, fmt.Errorf("%w: element %d", ErrOverflow, i)
	}
	return dst, nil
}

// MatCondition2 returns the 2-norm condition number κ₂(A) = σ_max/σ_min of an m×n matrix from its singular
// values (SVDJacobi), k = min(m, n) of them.
// Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element; ErrSingular
// when σ_min = 0; ErrNoConvergence from the SVD; ErrOverflow when the ratio is outside the float64 range.
//
// ID:      NUM-LINALG-162
// Time:    O(m·n·min(m, n)) per SVD sweep
// Space:   O(m·n), the SVD's working copy
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (‖A‖ = σ_max; c = σ_max/σ_min)
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖₂) and §8.6 "Computing the SVD"
// Tags:    matrix, n-dimensional, condition-number, svd, spectral, conditioning, linear-algebra
// Units:   a in unit A; the result dimensionless
// Example: [3 0; 4 5] → 3√5/√5 = 3 (hand-checked)
func MatCondition2(a []float64, m, n int) (float64, error) {
	svd, err := SVDJacobi(a, m, n, -1, false)
	if err != nil {
		return 0, err
	}
	lo := svd.S[len(svd.S)-1]
	if lo == 0 {
		return 0, fmt.Errorf("%w: smallest singular value is 0", ErrSingular)
	}
	c := svd.S[0] / lo
	if !reduce.AllFinite(c) {
		return 0, errOverflow("condition number")
	}
	return c, nil
}

// MatNorm2 returns the spectral norm ‖A‖₂ = σ_max, the largest singular value of an m×n matrix
// (SVDJacobi).
// Errors: ErrDimension for m or n < 1 or len(a) ≠ m·n; ErrNotFinite for a NaN or ±Inf element;
// ErrNoConvergence from the SVD.
//
// ID:      NUM-LINALG-166
// Time:    O(m·n·min(m, n)) per SVD sweep
// Space:   O(m·n), the SVD's working copy
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §2.3 "Matrix Norms" (‖A‖₂ is the square root of the largest eigenvalue of AᵀA)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §11.2 "Norms and Condition Numbers" (‖A‖ = σ_max)
// Tags:    matrix, n-dimensional, norm, spectral-norm, two-norm, svd, linear-algebra
// Units:   a in unit A; the result in A
// Example: [3 0; 4 5] → 3√5 ≈ 6.708204 (hand-checked)
func MatNorm2(a []float64, m, n int) (float64, error) {
	svd, err := SVDJacobi(a, m, n, -1, false)
	if err != nil {
		return 0, err
	}
	return svd.S[0], nil
}
