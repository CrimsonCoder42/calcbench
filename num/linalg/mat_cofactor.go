package linalg

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkSquare checks an n×n matrix (n ≥ 1) with finite entries.
func checkSquare(a []float64, n int) error {
	if n < 1 {
		return fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if err := checkMatrix("a", a, n, n); err != nil {
		return err
	}
	if i := firstNotFinite(a); i >= 0 {
		return notFiniteAt("a", i)
	}
	return nil
}

// submatrix writes a with row i and column j deleted into w ((n−1)×(n−1)).
func submatrix(w, a []float64, n, i, j int) {
	k := 0
	for r := range n {
		if r == i {
			continue
		}
		for c := range n {
			if c == j {
				continue
			}
			w[k] = a[r*n+c]
			k++
		}
	}
}

// MatMinor returns the minor M_ij of an n×n matrix (n ≥ 2): the determinant of the (n−1)×(n−1) matrix left
// after deleting row i and column j (0-based), by MatDet (LU with partial pivoting).
// Errors: ErrDimension when n < 2 or len(a) ≠ n²; ErrOutOfRange when i or j is outside 0…n−1; ErrNotFinite
// for a NaN or ±Inf entry; ErrOverflow when the minor is outside the float64 range.
//
// ID:      NUM-LINALG-127
// Time:    O(n³)   one LU of order n−1
// Space:   O((n−1)²) extra (the submatrix and its LU)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" (M_ij: remove row i and column j; C_ij = (−1)^(i+j) det M_ij)
// Tags:    matrix, n-dimensional, minor, cofactor, determinant, linear-algebra
// Units:   a in unit A; the result in A^(n−1)
// Example: [1 2 3; 4 5 6; 7 8 10], i = 1, j = 2 → det [1 2; 7 8] = −6 (hand-checked)
func MatMinor(a []float64, n, i, j int) (float64, error) {
	if n < 2 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 2", ErrDimension, n)
	}
	if err := checkSquare(a, n); err != nil {
		return 0, err
	}
	if i < 0 || i >= n || j < 0 || j >= n {
		return 0, fmt.Errorf("%w: (i, j) = (%d, %d), want 0…%d", ErrOutOfRange, i, j, n-1)
	}
	w := make([]float64, (n-1)*(n-1))
	submatrix(w, a, n, i, j)
	return MatDet(w, n-1)
}

// MatCofactorMatrix returns the matrix of cofactors C, C_ij = (−1)^(i+j)·M_ij, of an n×n matrix in dst. When
// A is nonsingular it uses C = det(A)·A⁻ᵀ from one LU (O(n³)); when A is singular (det A = 0) it computes
// every minor by its own LU (O(n⁵)), so singular matrices, whose cofactors may still be nonzero, get the
// exact definition. n = 1 gives C = [1]. dst must not overlap a.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOverflow when a
// cofactor is outside the float64 range.
//
// ID:      NUM-LINALG-128
// Time:    O(n³) when A is nonsingular; O(n⁵) when it is singular
// Space:   O(n²) extra; O(n²) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A, so C = det A·(A⁻¹)ᵀ)
// Tags:    matrix, n-dimensional, cofactor, cofactor-matrix, minor, linear-algebra
// Units:   a in unit A; the result in A^(n−1)
// Example: [1 2; 3 4] → [4 −3; −2 1] (hand-checked)
func MatCofactorMatrix(dst, a []float64, n int) ([]float64, error) {
	adj, err := MatAdjugate(nil, a, n)
	if err != nil {
		return nil, err
	}
	return MatTranspose(dst, adj, n, n)
}

// MatAdjugate returns the adjugate adj A = Cᵀ (the transposed cofactor matrix) of an n×n matrix in dst, so
// A·adj A = adj A·A = det(A)·I and det(adj A) = det(A)^(n−1). Nonsingular A: adj A = det(A)·A⁻¹ from one LU;
// singular A: every cofactor by its own LU (O(n⁵)). n = 1 gives [1]. dst must not overlap a.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOverflow when an
// entry is outside the float64 range.
//
// ID:      NUM-LINALG-129
// Time:    O(n³) when A is nonsingular; O(n⁵) when it is singular (n minors of order n−1 each; an O(n³) route through the SVD exists and is LATER)
// Space:   O(n²) extra; O(n²) for the result, one allocation only if dst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (A·Cᵀ = (det A)·I)
// Tags:    matrix, n-dimensional, adjugate, classical-adjoint, cofactor, linear-algebra
// Units:   a in unit A; the result in A^(n−1)
// Example: [1 2; 2 4] (singular) → [4 −2; −2 1] (hand-checked: A·adj A = 0 = det A·I)
func MatAdjugate(dst, a []float64, n int) ([]float64, error) {
	if err := checkSquare(a, n); err != nil {
		return nil, err
	}
	dst = grow(dst, n*n)
	if n == 1 {
		dst[0] = 1
		return dst, nil
	}
	det, err := MatDet(a, n)
	if err != nil {
		return nil, err
	}
	if det != 0 {
		inv, err := MatInverse(nil, a, n)
		if err == nil {
			for i := range inv {
				dst[i] = det * inv[i]
			}
			if i := firstNotFinite(dst); i >= 0 {
				return nil, fmt.Errorf("%w: entry %d", ErrOverflow, i)
			}
			return dst, nil
		}
		if !errors.Is(err, ErrSingular) {
			return nil, err
		}
	}
	w := make([]float64, (n-1)*(n-1))
	for i := range n {
		for j := range n {
			submatrix(w, a, n, i, j)
			m, err := MatDet(w, n-1)
			if err != nil {
				return nil, err
			}
			if (i+j)%2 == 1 {
				m = -m
			}
			dst[j*n+i] = m
		}
	}
	return dst, nil
}

// MatCharPolyFaddeevLeVerrier returns the characteristic polynomial det(λI − A) of an n×n matrix as n + 1
// ascending coefficients (D-006; the λⁿ coefficient is 1) by the Faddeev–LeVerrier recursion
// M₀ = 0, c_n = 1; Mₖ = A·Mₖ₋₁ + c_{n−k+1}·I, c_{n−k} = −tr(A·Mₖ)/k. It needs only matrix products and traces
// and is exact in exact arithmetic, but in floating point the coefficients can lose all accuracy for large or
// ill-scaled n (the recursion amplifies rounding error roughly like ‖A‖ⁿ/n!·u), so it is for small, well-
// scaled matrices; roots of the result are poorly conditioned functions of the coefficients, so eigenvalues
// should come from MatSymmetricEigenJacobi or Mat2Eigen, not from these roots. dst must not overlap a.
// Errors: ErrDimension when n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf entry; ErrOverflow when a
// coefficient leaves the float64 range.
//
// ID:      NUM-LINALG-148
// Time:    O(n⁴)   n matrix products of order n
// Space:   O(n²) extra (two work matrices); O(n) for the result, one allocation only if dst lacks capacity
// Source:  Hou, "Classroom Note: A Simple Proof of the Leverrier–Faddeev Characteristic Polynomial Algorithm", SIAM Review, 1998, 40(3):706–709, doi:10.1137/S003614459732076X (the recursive relations of the Leverrier–Faddeev method and their proof)
// Source:  Faddeev, Faddeeva, Computational Methods of Linear Algebra, W. H. Freeman, 1963 (Leverrier's method in Faddeev's modification; attribution, no page cited)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (det(A − λI); the sum of the eigenvalues is the trace)
// Tags:    matrix, n-dimensional, characteristic-polynomial, faddeev-leverrier, eigenvalues, linear-algebra
// Units:   a in unit A; coefficient k in A^(n−k)
// Example: [2 1; 1 2] → λ² − 4λ + 3, i.e. [3, −4, 1] (hand-checked)
func MatCharPolyFaddeevLeVerrier(dst, a []float64, n int) ([]float64, error) {
	if err := checkSquare(a, n); err != nil {
		return nil, err
	}
	dst = grow(dst, n+1)
	dst[n] = 1
	m := make([]float64, 2*n*n)
	mk, am := m[:n*n], m[n*n:] // Mₖ and A·Mₖ
	for k := 1; k <= n; k++ {
		// Mₖ = A·Mₖ₋₁ + c_{n−k+1}·I (with A·M₀ = 0)
		if k == 1 {
			clear(mk)
		} else {
			copy(mk, am)
		}
		for i := range n {
			mk[i*n+i] += dst[n-k+1]
		}
		if _, err := MatMul(am, a, n, n, mk, n); err != nil {
			return nil, err
		}
		tr, err := MatTrace(am, n)
		if err != nil {
			return nil, err
		}
		dst[n-k] = -tr / float64(k)
		if !reduce.AllFinite(dst[n-k]) {
			return nil, fmt.Errorf("%w: coefficient %d", ErrOverflow, n-k)
		}
	}
	return dst, nil
}
