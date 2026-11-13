package linalg

import (
	"errors"
	"fmt"
	"math"
)

// ErrRankDeficient is returned when a matrix that must have full column rank has a column that is (to the
// stated tolerance) a combination of the others.
var ErrRankDeficient = errors.New("linalg: matrix is rank deficient")

// householderQR overwrites w (m×n, m ≥ n) with R in its upper triangle and returns the n Householder
// vectors (column k stored in v[k*m:(k+1)*m], zero above row k) and their β = 2/(vᵀv) (0 for an identity
// reflector).
func householderQR(w []float64, m, n int) (v, beta []float64) {
	v = make([]float64, n*m+n)
	beta = v[n*m:]
	v = v[:n*m]
	for k := range n {
		alpha := 0.0
		for i := k; i < m; i++ {
			alpha = math.Hypot(alpha, w[i*n+k])
		}
		if alpha == 0 {
			continue
		}
		if w[k*n+k] > 0 {
			alpha = -alpha // v = x − alpha·e₁ with alpha of opposite sign to x₁: no cancellation
		}
		vk := v[k*m : k*m+m]
		for i := k; i < m; i++ {
			vk[i] = w[i*n+k]
		}
		vk[k] -= alpha
		// Normalise v by its largest entry: the reflector I − βvvᵀ does not change, and vᵀv ≥ 1 can
		// neither underflow nor overflow.
		big := 0.0
		for i := k; i < m; i++ {
			big = max(big, abs(vk[i]))
		}
		for i := k; i < m; i++ {
			vk[i] /= big
		}
		vtv := 0.0
		for i := k; i < m; i++ {
			vtv += vk[i] * vk[i]
		}
		beta[k] = 2 / vtv
		for j := k; j < n; j++ { // w ← (I − β v vᵀ) w on columns k…n−1
			s := 0.0
			for i := k; i < m; i++ {
				s += vk[i] * w[i*n+j]
			}
			s *= beta[k]
			for i := k; i < m; i++ {
				w[i*n+j] -= s * vk[i]
			}
		}
		for i := k + 1; i < m; i++ {
			w[i*n+k] = 0
		}
	}
	return v, beta
}

// checkTall checks an m×n matrix with m ≥ n ≥ 1 and finite entries.
func checkTall(a []float64, m, n int) error {
	if n < 1 || m < n {
		return fmt.Errorf("%w: a is %d×%d, want m ≥ n ≥ 1", ErrDimension, m, n)
	}
	if err := checkMatrix("a", a, m, n); err != nil {
		return err
	}
	if i := firstNotFinite(a); i >= 0 {
		return notFiniteAt("a", i)
	}
	return nil
}

// QRHouseholder returns the thin QR factorisation A = Q·R of an m×n matrix (m ≥ n) by Householder
// reflections: Q is m×n with orthonormal columns (in qDst) and R is n×n upper triangular (in rDst). Each
// reflector is I − βvvᵀ with v chosen so that no cancellation occurs (Householder 1958); the factorisation is
// backward stable. Neither dst may overlap a.
// Errors: ErrDimension unless m ≥ n ≥ 1 and len(a) = m·n; ErrNotFinite for a NaN or ±Inf element;
// ErrOverflow when an element of Q or R leaves the float64 range.
//
// ID:      NUM-LINALG-154
// Time:    O(m·n²)   about 4mn² − 4n³/3 flops including forming Q
// Space:   O(m·n) extra (the reflectors and a working copy); O(m·n + n²) for the results
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §5.2 "The QR Factorization" (Householder QR; the thin QR factorization) and the Householder vector computation it uses
// Source:  Householder, "Unitary Triangularization of a Nonsymmetric Matrix", Journal of the ACM, 1958, 5(4):339–342, doi:10.1145/320941.320947
// Tags:    matrix, n-dimensional, qr, householder, factorization, orthogonal, linear-algebra
// Units:   a in unit A; Q dimensionless; R in A
// Example: [3 0; 4 5] → |R| = [5 4; 0 3], Q·R = A and QᵀQ = I (hand-checked; the signs follow the reflectors: here R = [−5 −4; 0 −3])
func QRHouseholder(qDst, rDst, a []float64, m, n int) ([]float64, []float64, error) {
	if err := checkTall(a, m, n); err != nil {
		return nil, nil, err
	}
	// Scale A by an exact power of two into [½, 1) in magnitude, so no column norm over- or underflows;
	// R scales back exactly and Q does not change.
	big, _ := maxAbs(a)
	_, e := pow2Scale(max(big, 0x1p-1000))
	w := make([]float64, m*n)
	for i, x := range a {
		w[i] = math.Ldexp(x, -e)
	}
	v, beta := householderQR(w, m, n)
	rDst = grow(rDst, n*n)
	for i := range n {
		for j := range n {
			if j >= i {
				rDst[i*n+j] = math.Ldexp(w[i*n+j], e)
			} else {
				rDst[i*n+j] = 0
			}
		}
	}
	// Q = H₀H₁…Hₙ₋₁ applied to the first n columns of I, right to left.
	qDst = grow(qDst, m*n)
	clear(qDst)
	for i := range n {
		qDst[i*n+i] = 1
	}
	for k := n - 1; k >= 0; k-- {
		if beta[k] == 0 {
			continue
		}
		vk := v[k*m : k*m+m]
		for j := range n {
			s := 0.0
			for i := k; i < m; i++ {
				s += vk[i] * qDst[i*n+j]
			}
			s *= beta[k]
			for i := k; i < m; i++ {
				qDst[i*n+j] -= s * vk[i]
			}
		}
	}
	if i := firstNotFinite(rDst); i >= 0 {
		return nil, nil, fmt.Errorf("%w: R element %d", ErrOverflow, i)
	}
	return qDst, rDst, nil
}

// GramSchmidtModified returns the thin QR factorisation A = Q·R of an m×n matrix (m ≥ n) by the modified
// Gram–Schmidt process, which orthogonalises each remaining column against qₖ as soon as qₖ is formed. Its
// loss of orthogonality is ‖QᵀQ − I‖ ≈ u·κ(A), far better than classical Gram–Schmidt but worse than
// Householder (use QRHouseholder when κ(A) is large). Q is m×n in qDst, R n×n in rDst; neither may overlap a.
// Errors: ErrDimension unless m ≥ n ≥ 1 and len(a) = m·n; ErrNotFinite for a NaN or ±Inf element;
// ErrRankDeficient when a column's remaining norm is ≤ tol times its original norm (tol < 0 selects
// n·2⁻⁵²), i.e. it depends on earlier columns.
//
// ID:      NUM-LINALG-155
// Time:    O(m·n²)   2mn² flops
// Space:   O(n) extra (the original column norms); O(m·n + n²) for the results, one allocation each only if a dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §5.2 "The QR Factorization" (classical and modified Gram–Schmidt; the MGS orthogonality bound)
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.4 "Orthonormal Bases and Gram-Schmidt" (subtract projections one at a time; A = QR)
// Tags:    matrix, n-dimensional, gram-schmidt, qr, orthonormal, factorization, linear-algebra
// Units:   a in unit A; Q dimensionless; R in A
// Example: [1 1; 1 0; 0 1] → q₁ = (1, 1, 0)/√2, q₂ = (1, −1, 2)/√6, R = [√2 1/√2; 0 √(3/2)] (hand-checked)
func GramSchmidtModified(qDst, rDst, a []float64, m, n int, tol float64) ([]float64, []float64, error) {
	if err := checkTall(a, m, n); err != nil {
		return nil, nil, err
	}
	if tol < 0 {
		tol = float64(n) * 0x1p-52
	}
	qDst = grow(qDst, m*n)
	copy(qDst, a)
	rDst = grow(rDst, n*n)
	clear(rDst)
	orig := make([]float64, n)
	for j := range n {
		for i := range m {
			orig[j] = math.Hypot(orig[j], qDst[i*n+j])
		}
	}
	for k := range n {
		nk := 0.0
		for i := range m {
			nk = math.Hypot(nk, qDst[i*n+k])
		}
		if nk <= tol*orig[k] || nk == 0 {
			return nil, nil, fmt.Errorf("%w: column %d", ErrRankDeficient, k)
		}
		rDst[k*n+k] = nk
		for i := range m {
			qDst[i*n+k] /= nk
		}
		for j := k + 1; j < n; j++ {
			s := 0.0
			for i := range m {
				s += qDst[i*n+k] * qDst[i*n+j]
			}
			rDst[k*n+j] = s
			for i := range m {
				qDst[i*n+j] -= s * qDst[i*n+k]
			}
		}
	}
	return qDst, rDst, nil
}

// LeastSquaresResult is the solution of min ‖A·x − b‖₂: the coefficients X and the residual norm.
type LeastSquaresResult struct {
	X            []float64
	ResidualNorm float64
}

// LeastSquaresQR solves the full-rank least-squares problem min ‖A·x − b‖₂ for an m×n matrix (m ≥ n) by
// Householder QR: Qᵀb is formed by applying the reflectors, R·x = (Qᵀb)₀…ₙ₋₁ is solved by back substitution,
// and the residual norm is ‖(Qᵀb)ₙ…ₘ₋₁‖₂. Normal equations are never formed (they square κ(A)). xDst
// receives x (length n; it must not overlap a or b).
// Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf entry; ErrRankDeficient when
// |rₖₖ| ≤ tol·max|rᵢᵢ| for some k (tol < 0 selects max(m, n)·2⁻⁵²); ErrOverflow when a coefficient leaves
// the float64 range.
//
// ID:      NUM-LINALG-156
// Time:    O(m·n²)   2mn² − 2n³/3 flops
// Space:   O(m·n) extra (a working copy and the reflectors); O(n) for the result
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.3 "Least Squares Approximations" and §4.4 "Orthonormal Bases and Gram-Schmidt" (A = QR turns AᵀAx̂ = Aᵀb into Rx̂ = Qᵀb)
// Tags:    matrix, n-dimensional, least-squares, qr, regression, overdetermined, linear-algebra
// Units:   a in unit A, b in unit B; x in B/A; the residual norm in B
// Example: fit y = c + d·t to (0, 6), (1, 0), (2, 0): a = [1 0; 1 1; 1 2], b = (6, 0, 0) → x = (5, −3), residual √6 (Strang §4.3 worked example)
func LeastSquaresQR(xDst, a []float64, m, n int, b []float64, tol float64) (LeastSquaresResult, error) {
	if err := checkTall(a, m, n); err != nil {
		return LeastSquaresResult{}, err
	}
	if len(b) != m {
		return LeastSquaresResult{}, lengthError("b", len(b), "m", m)
	}
	if i := firstNotFinite(b); i >= 0 {
		return LeastSquaresResult{}, notFiniteAt("b", i)
	}
	if tol < 0 {
		tol = float64(max(m, n)) * 0x1p-52
	}
	// Scale A and b by exact powers of two (2^ea, 2^eb) so nothing over- or underflows; then
	// x = 2^(eb−ea)·x' and the residual norm is 2^eb·r'.
	bigA, _ := maxAbs(a)
	_, ea := pow2Scale(max(bigA, 0x1p-1000))
	bigB, _ := maxAbs(b)
	_, eb := pow2Scale(max(bigB, 0x1p-1000))
	w := make([]float64, m*n+m)
	qtb := w[m*n:]
	w = w[:m*n]
	for i, x := range a {
		w[i] = math.Ldexp(x, -ea)
	}
	for i, x := range b {
		qtb[i] = math.Ldexp(x, -eb)
	}
	v, beta := householderQR(w, m, n)
	for k := range n {
		if beta[k] == 0 {
			continue
		}
		vk := v[k*m : k*m+m]
		s := 0.0
		for i := k; i < m; i++ {
			s += vk[i] * qtb[i]
		}
		s *= beta[k]
		for i := k; i < m; i++ {
			qtb[i] -= s * vk[i]
		}
	}
	rmax := 0.0
	for k := range n {
		rmax = max(rmax, abs(w[k*n+k]))
	}
	for k := range n {
		if abs(w[k*n+k]) <= tol*rmax || rmax == 0 {
			return LeastSquaresResult{}, fmt.Errorf("%w: |r[%d,%d]| = %g", ErrRankDeficient, k, k, abs(w[k*n+k]))
		}
	}
	xDst = grow(xDst, n)
	for i := n - 1; i >= 0; i-- {
		s := qtb[i]
		for j := i + 1; j < n; j++ {
			s -= w[i*n+j] * xDst[j]
		}
		xDst[i] = s / w[i*n+i]
	}
	for i := range xDst {
		xDst[i] = math.Ldexp(xDst[i], eb-ea)
	}
	if i := firstNotFinite(xDst); i >= 0 {
		return LeastSquaresResult{}, fmt.Errorf("%w: x[%d]", ErrOverflow, i)
	}
	res := 0.0
	for i := n; i < m; i++ {
		res = math.Hypot(res, qtb[i])
	}
	return LeastSquaresResult{X: xDst, ResidualNorm: math.Ldexp(res, eb)}, nil
}
