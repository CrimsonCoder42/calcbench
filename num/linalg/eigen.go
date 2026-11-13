package linalg

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrNoConvergence is returned when an iterative method does not meet its tolerance within its iteration
// limit.
var ErrNoConvergence = errors.New("linalg: iteration did not converge")

// Mat2EigenResult holds the eigenvalues of a real 2×2 matrix, sorted by real part then imaginary part, and
// a unit eigenvector (2-norm 1) for each. Vectors[k] belongs to Values[k]. Defective is true when the two
// eigenvalues are equal but the matrix is not a multiple of I (one independent eigenvector only; both
// entries of Vectors then hold it).
type Mat2EigenResult struct {
	Values    [2]complex128
	Vectors   [2][2]complex128
	Defective bool
}

// Mat2Eigen returns the eigenvalues and eigenvectors of a general real 2×2 matrix m = [a b; c d]: the roots
// λ = (a + d)/2 ± √q of λ² − (a + d)λ + (ad − bc) = 0, with the discriminant q = ((a − d)/2)² + bc formed
// directly (not as trace² − 4·det, which loses √u of accuracy for nearly equal eigenvalues), a real pair or a
// complex-conjugate pair; the root of smaller modulus comes from Vieta's det/λ, so neither cancels. For each λ
// the eigenvector is (b, λ − a) or (λ − d, c), whichever is larger, normalised to unit length. m is scaled by an
// exact power of two first, so no intermediate over- or underflows.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an eigenvalue is outside the float64 range.
//
// ID:      NUM-LINALG-149
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (λ² − (trace)λ + det = 0; eigenvectors from the null space of A − λI; complex eigenvalues of a real matrix come in conjugate pairs)
// Tags:    matrix, 2x2, eigenvalues, eigenvectors, complex, linear-algebra
// Units:   m in unit A; eigenvalues in A; eigenvectors dimensionless
// Example: [4 1; 2 3] → λ = 2, 5 with vectors ∝ (1, −2), (1, 1); [0 −1; 1 0] → λ = ±i (hand-checked)
func Mat2Eigen(m Mat2) (Mat2EigenResult, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat2EigenResult{}, errNotFinite("m")
	}
	// Scale by an exact power of two so that no square or product below overflows or underflows; the
	// eigenvalues scale back exactly and the eigenvectors do not change.
	big, _ := maxAbs(m[:])
	_, e := pow2Scale(max(big, 0x1p-1000))
	a, b, c, d := math.Ldexp(m[0], -e), math.Ldexp(m[1], -e), math.Ldexp(m[2], -e), math.Ldexp(m[3], -e)
	mid, half := a/2+d/2, a/2-d/2
	var disc dot2 // ((a − d)/2)² + bc, exact enough to separate nearly equal eigenvalues
	disc.add(half, half)
	disc.add(b, c)
	q := disc.result()
	var l1, l2 complex128
	if q >= 0 {
		r := math.Sqrt(q)
		big := mid + math.Copysign(r, mid) // no cancellation
		small := mid - math.Copysign(r, mid)
		if big != 0 {
			small = det2(a, b, c, d) / big // Vieta: λ₁λ₂ = det
		}
		l1, l2 = complex(min(big, small), 0), complex(max(big, small), 0)
	} else {
		im := math.Sqrt(-q)
		l1, l2 = complex(mid, -im), complex(mid, im)
	}
	res := Mat2EigenResult{Values: [2]complex128{l1, l2}}
	scalar := b == 0 && c == 0 && a == d
	for k, lam := range res.Values {
		if scalar {
			res.Vectors[k][k] = 1
			continue
		}
		v1 := [2]complex128{complex(b, 0), lam - complex(a, 0)}
		v2 := [2]complex128{lam - complex(d, 0), complex(c, 0)}
		v := v1
		if cnorm2(v2) > cnorm2(v1) {
			v = v2
		}
		n := cnorm2(v)
		res.Vectors[k] = [2]complex128{v[0] / complex(n, 0), v[1] / complex(n, 0)}
	}
	res.Defective = !scalar && l1 == l2
	for k := range res.Values {
		v := res.Values[k]
		res.Values[k] = complex(math.Ldexp(real(v), e), math.Ldexp(imag(v), e))
		if cmplx.IsInf(res.Values[k]) {
			return Mat2EigenResult{}, errOverflow("eigenvalue")
		}
	}
	return res, nil
}

// cnorm2 is the 2-norm of a complex 2-vector without overflow of the squares.
func cnorm2(v [2]complex128) float64 { return math.Hypot(cmplx.Abs(v[0]), cmplx.Abs(v[1])) }

// maxJacobiSweeps bounds the cyclic Jacobi method; convergence is quadratic, so real matrices need far
// fewer than 10 sweeps (Rutishauser 1966).
const maxJacobiSweeps = 60

// MatSymmetricEigenJacobi returns the eigenvalues of a symmetric n×n matrix in ascending order and, when
// wantVectors is true, an n×n matrix whose column k is a unit eigenvector for value k, by the cyclic Jacobi
// method: sweeps of 2×2 symmetric Schur rotations over every off-diagonal pair (p, q) in row order until the
// off-diagonal Frobenius norm is at most tol·‖A‖_F (tol < 0 selects n·2⁻⁵²). Jacobi computes small
// eigenvalues to high relative accuracy (Demmel–Veselić 1992). Results go to valuesDst (length n) and
// vectorsDst (length n², untouched when wantVectors is false); neither may overlap a.
// Errors: ErrDimension for n < 1 or len(a) ≠ n²; ErrNotFinite for a NaN or ±Inf element; ErrNotSymmetric when
// |aᵢⱼ − aⱼᵢ| > tol'·max|a| for some pair (tol' = max(tol, default)); ErrNoConvergence after 60 sweeps.
//
// ID:      NUM-LINALG-150
// Time:    O(n³) per sweep   n = order; typically 5–10 sweeps
// Space:   O(n²) extra (the working copy and, with vectors, one n² buffer to reorder them; O(n) for the sort index); O(n) and O(n²) for the results, one allocation each only if a dst lacks capacity
// Source:  Golub, Van Loan, Matrix Computations, 4th ed., Johns Hopkins University Press, 2013, §8.5 "Jacobi Methods" (the 2-by-2 symmetric Schur decomposition, the cyclic-by-row Jacobi algorithm, off(A))
// Source:  Rutishauser, "The Jacobi method for real symmetric matrices", Numerische Mathematik, 1966, 9(1):1–10, doi:10.1007/BF02165223
// Source:  Demmel, Veselić, "Jacobi's Method is More Accurate than QR", SIAM Journal on Matrix Analysis and Applications, 1992, 13(4):1204–1245, doi:10.1137/0613074
// Tags:    matrix, n-dimensional, eigenvalues, eigenvectors, symmetric, jacobi, linear-algebra
// Units:   a in unit A; eigenvalues in A; eigenvectors dimensionless
// Example: [2 1; 1 2] → (1, 3) with vectors ∝ (1, −1), (1, 1); [2 −1 0; −1 2 −1; 0 −1 2] → 2 − √2, 2, 2 + √2 (hand-checked)
func MatSymmetricEigenJacobi(valuesDst, vectorsDst, a []float64, n int, tol float64, wantVectors bool) ([]float64, []float64, error) {
	if n < 1 {
		return nil, nil, fmt.Errorf("%w: n = %d", ErrDimension, n)
	}
	if err := checkMatrix("a", a, n, n); err != nil {
		return nil, nil, err
	}
	if i := firstNotFinite(a); i >= 0 {
		return nil, nil, notFiniteAt("a", i)
	}
	def := float64(n) * 0x1p-52
	if tol < 0 {
		tol = def
	}
	big, _ := maxAbs(a)
	symTol := max(tol, def) * big
	for i := range n {
		for j := i + 1; j < n; j++ {
			if abs(a[i*n+j]-a[j*n+i]) > symTol {
				return nil, nil, fmt.Errorf("%w: a[%d,%d] = %g, a[%d,%d] = %g", ErrNotSymmetric, i, j, a[i*n+j], j, i, a[j*n+i])
			}
		}
	}
	// Work on a copy scaled by a power of two so that squares cannot overflow.
	w := make([]float64, n*n)
	_, e := pow2Scale(max(big, 0x1p-1000))
	for i := range a {
		w[i] = math.Ldexp(a[i], -e)
	}
	var v []float64
	if wantVectors {
		vectorsDst, _ = MatIdentity(vectorsDst, n)
		v = vectorsDst
	}
	frob, _ := VecNorm(w)
	for sweep := 0; ; sweep++ {
		off := 0.0
		for i := range n {
			for j := range n {
				if i != j {
					off = math.Hypot(off, w[i*n+j])
				}
			}
		}
		if off <= tol*frob {
			break
		}
		if sweep == maxJacobiSweeps {
			return nil, nil, fmt.Errorf("%w: off(A)/‖A‖ = %g after %d sweeps", ErrNoConvergence, off/frob, sweep)
		}
		for p := 0; p < n-1; p++ {
			for q := p + 1; q < n; q++ {
				apq := w[p*n+q]
				if apq == 0 {
					continue
				}
				theta := (w[q*n+q] - w[p*n+p]) / (2 * apq)
				t := 1 / (abs(theta) + math.Sqrt(1+theta*theta))
				if theta < 0 {
					t = -t
				}
				if math.IsInf(theta, 0) {
					t = 0
				}
				c := 1 / math.Sqrt(1+t*t)
				s := t * c
				for k := range n { // columns p and q: W ← W·J
					wkp, wkq := w[k*n+p], w[k*n+q]
					w[k*n+p] = c*wkp - s*wkq
					w[k*n+q] = s*wkp + c*wkq
				}
				for k := range n { // rows p and q: W ← Jᵀ·W
					wpk, wqk := w[p*n+k], w[q*n+k]
					w[p*n+k] = c*wpk - s*wqk
					w[q*n+k] = s*wpk + c*wqk
				}
				w[p*n+q], w[q*n+p] = 0, 0
				if wantVectors {
					for k := range n {
						vkp, vkq := v[k*n+p], v[k*n+q]
						v[k*n+p] = c*vkp - s*vkq
						v[k*n+q] = s*vkp + c*vkq
					}
				}
			}
		}
	}
	valuesDst = grow(valuesDst, n)
	idx := make([]int, n)
	for i := range n {
		idx[i] = i
	}
	slices.SortStableFunc(idx, func(x, y int) int { return cmpFloat(w[x*n+x], w[y*n+y]) })
	for k, i := range idx {
		valuesDst[k] = math.Ldexp(w[i*n+i], e)
	}
	if wantVectors {
		tmp := make([]float64, n*n)
		for k, i := range idx {
			for r := range n {
				tmp[r*n+k] = v[r*n+i]
			}
		}
		copy(v, tmp)
		return valuesDst, v, nil
	}
	return valuesDst, nil, nil
}

func cmpFloat(x, y float64) int {
	switch {
	case x < y:
		return -1
	case x > y:
		return 1
	}
	return 0
}

// PowerIterationResult is the dominant eigenpair found by MatPowerIteration: the eigenvalue (the Rayleigh
// quotient of the final vector), a unit eigenvector, and the number of iterations used.
type PowerIterationResult struct {
	Value      float64
	Vector     []float64
	Iterations int
}

// MatPowerIteration returns the dominant eigenpair of an n×n matrix by the power method: x ← A·x/‖A·x‖,
// λ = xᵀA·x (Rayleigh quotient), stopping when ‖A·x − λx‖₂ ≤ tol·|λ| (tol < 0 selects 1e−12) or after
// maxIter steps. It converges when A has a single eigenvalue of largest modulus, real, with |λ₁| > |λ₂|, and
// x0 has a component along its eigenvector; the error then falls like |λ₂/λ₁|^k. vecDst receives the
// vector (length n; it must not overlap a or x0).
// Errors: ErrDimension for mismatched sizes; ErrNotFinite for a NaN or ±Inf element; ErrZeroVector when x0
// or an iterate A·x is zero; ErrOutOfRange when maxIter < 1; ErrNoConvergence after maxIter steps (the
// condition above fails: equal-modulus or complex dominant eigenvalues).
//
// ID:      NUM-LINALG-151
// Time:    O(k·n²)   k = iterations used ≤ maxIter
// Space:   O(n) extra (one work vector); O(n) for the result, one allocation only if vecDst lacks capacity
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (uₖ = Aᵏu₀ = Σ cᵢλᵢᵏxᵢ is dominated by the eigenvalue of largest |λ|) and §6.4 "Symmetric Matrices" (the Rayleigh quotient xᵀAx/xᵀx)
// Tags:    matrix, n-dimensional, eigenvalues, power-method, power-iteration, dominant, linear-algebra
// Units:   a in unit A; the eigenvalue in A; the vector dimensionless
// Example: [2 1; 1 2], x0 = (1, 0) → λ = 3, x ∝ (1, 1) (hand-checked)
func MatPowerIteration(vecDst, a []float64, n int, x0 []float64, tol float64, maxIter int) (PowerIterationResult, error) {
	if err := checkSquareSystem(a, n, x0); err != nil {
		return PowerIterationResult{}, err
	}
	if maxIter < 1 {
		return PowerIterationResult{}, fmt.Errorf("%w: maxIter = %d", ErrOutOfRange, maxIter)
	}
	if tol < 0 {
		tol = 1e-12
	}
	x, err := VecNormalize(vecDst, x0)
	if err != nil {
		return PowerIterationResult{}, fmt.Errorf("x0: %w", err)
	}
	y := make([]float64, n)
	for k := 1; k <= maxIter; k++ {
		if _, err := MatVecMul(y, a, n, n, x); err != nil {
			return PowerIterationResult{}, err
		}
		if big, _ := maxAbs(y); big == 0 {
			// x lies in the null space of A: it has no component along the dominant eigenvector.
			return PowerIterationResult{}, fmt.Errorf("%w: A·x = 0 at step %d", ErrZeroVector, k)
		}
		lam, err := VecDot(x, y)
		if err != nil {
			return PowerIterationResult{}, err
		}
		res := 0.0
		for i := range n {
			res = math.Hypot(res, y[i]-lam*x[i])
		}
		if res <= tol*abs(lam) {
			return PowerIterationResult{Value: lam, Vector: x, Iterations: k}, nil
		}
		if _, err := VecNormalize(x, y); err != nil {
			return PowerIterationResult{}, fmt.Errorf("A·x: %w", err)
		}
	}
	return PowerIterationResult{}, fmt.Errorf("%w: %d iterations", ErrNoConvergence, maxIter)
}
