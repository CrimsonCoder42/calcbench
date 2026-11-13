package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Mat2 is a 2×2 matrix stored row-major: [m00, m01, m10, m11]. It is a value: functions take and return
// copies and never allocate. Vectors are columns: Mat2MulVec computes M·v.
type Mat2 [4]float64

// sumProd returns p·q + r·s within 2 ulps (Kahan's 2×2 determinant algorithm on p·q − (−r)·s).
func sumProd(p, q, r, s float64) float64 { return det2(p, -r, s, q) }

// Mat2Identity returns the 2×2 identity matrix.
//
// ID:      NUM-LINALG-056
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I)
// Tags:    matrix, 2x2, identity, linear-algebra
// Units:   dimensionless
// Example: → [1 0; 0 1]
func Mat2Identity() Mat2 { return Mat2{1, 0, 0, 1} }

// Mat2Mul returns the matrix product a·b. Each element is a two-term dot product computed within
// 2 ulps (Kahan's algorithm).
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the product is outside
// the float64 range.
//
// ID:      NUM-LINALG-057
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column)
// Source:  Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8
// Tags:    matrix, 2x2, multiply, product, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: [1 2; 3 4]·[5 6; 7 8] → [19 22; 43 50] (hand-checked)
func Mat2Mul(a, b Mat2) (Mat2, error) {
	if !reduce.AllFinite(a[:]...) {
		return Mat2{}, errNotFinite("a")
	}
	if !reduce.AllFinite(b[:]...) {
		return Mat2{}, errNotFinite("b")
	}
	c := Mat2{
		sumProd(a[0], b[0], a[1], b[2]), sumProd(a[0], b[1], a[1], b[3]),
		sumProd(a[2], b[0], a[3], b[2]), sumProd(a[2], b[1], a[3], b[3]),
	}
	if !reduce.AllFinite(c[:]...) {
		return Mat2{}, errOverflow("a·b")
	}
	return c, nil
}

// Mat2MulVec returns the matrix-vector product m·v (v a column), each element within 2 ulps. This is
// the "2×2 matrix times 2×1 matrix" product.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element is outside the float64 range.
//
// ID:      NUM-LINALG-058
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector)
// Source:  Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8
// Tags:    matrix, 2x2, matrix-vector, multiply, linear-algebra
// Units:   m in unit A, v in unit B; the result in A·B
// Example: [1 2; 3 4]·(5, 6) → (17, 39) (hand-checked)
func Mat2MulVec(m Mat2, v Vec2) (Vec2, error) {
	if !reduce.AllFinite(m[:]...) {
		return Vec2{}, errNotFinite("m")
	}
	if !finite2(v) {
		return Vec2{}, errNotFinite("v")
	}
	r := Vec2{sumProd(m[0], v[0], m[1], v[1]), sumProd(m[2], v[0], m[3], v[1])}
	if !finite2(r) {
		return Vec2{}, errOverflow("m·v")
	}
	return r, nil
}

// Mat2Transpose returns mᵀ.
// Errors: ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-059
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations"
// Tags:    matrix, 2x2, transpose, linear-algebra
// Units:   the unit of m
// Example: [1 2; 3 4] → [1 3; 2 4]
func Mat2Transpose(m Mat2) (Mat2, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat2{}, errNotFinite("m")
	}
	return Mat2{m[0], m[2], m[1], m[3]}, nil
}

// Mat2Det returns det m = m₀₀m₁₁ − m₀₁m₁₀, within 2 ulps by Kahan's fused multiply-add algorithm, even
// when the two products nearly cancel.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the determinant is outside the
// float64 range.
//
// ID:      NUM-LINALG-060
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.1 "The Properties of Determinants"
// Source:  Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8
// Tags:    matrix, 2x2, determinant, linear-algebra
// Units:   m in unit A; the result in A²
// Example: [1 2; 3 4] → −2 (hand-checked)
func Mat2Det(m Mat2) (float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return 0, errNotFinite("m")
	}
	d := det2(m[0], m[1], m[2], m[3])
	if !reduce.AllFinite(d) {
		return 0, errOverflow("det m")
	}
	return d, nil
}

// Mat2Inverse returns m⁻¹ = [m₁₁ −m₀₁; −m₁₀ m₀₀]/det m. Each row and then each column is first scaled by an exact
// power of two (equilibration), so the determinant neither overflows nor underflows for a badly scaled but
// well-conditioned matrix. An ill-conditioned matrix has an inaccurate inverse; check with
// MatConditionEstimate1.
// Errors: ErrSingular when det m = 0; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element
// of the inverse is outside the float64 range.
//
// ID:      NUM-LINALG-061
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the 2 by 2 inverse) and §5.3 "Cramer's Rule, Inverses, and Volumes"
// Tags:    matrix, 2x2, inverse, linear-algebra
// Units:   m in unit A; the result in 1/A
// Example: [4 7; 2 6] → [0.6 −0.7; −0.2 0.4] (hand-checked: det = 10)
func Mat2Inverse(m Mat2) (Mat2, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat2{}, errNotFinite("m")
	}
	var w Mat2
	var e [2]int
	var f [2]int
	if !rowEquilibrate(m[:], w[:], e[:], 2) {
		return Mat2{}, fmt.Errorf("%w: a row is zero", ErrSingular)
	}
	colEquilibrate(w[:], f[:], 2)
	d := det2(w[0], w[1], w[2], w[3])
	if d == 0 {
		return Mat2{}, fmt.Errorf("%w: det m = 0", ErrSingular)
	}
	iw := Mat2{w[3] / d, -w[1] / d, -w[2] / d, w[0] / d}
	var r Mat2
	if err := unscaleInverse(r[:], iw[:], f[:], e[:], 2); err != nil {
		return Mat2{}, err
	}
	return r, nil
}

// rowEquilibrate writes into w the n×n matrix m with each row scaled by the exact power of two 2^(−eᵢ)
// that brings its largest |element| into [1, 2) (row equilibration), and the exponents into e. For
// w = S·m with S = diag(2^(−eᵢ)), m⁻¹ = w⁻¹·S. It reports false when a row is all zeros (m singular).
func rowEquilibrate(m, w []float64, e []int, n int) bool {
	for i := range n {
		mx := 0.0
		for j := range n {
			mx = max(mx, math.Abs(m[i*n+j]))
		}
		if mx == 0 {
			return false
		}
		_, e[i] = pow2Scale(mx)
		for j := range n {
			w[i*n+j] = math.Ldexp(m[i*n+j], -e[i])
		}
	}
	return true
}

// colEquilibrate scales each column of the row-equilibrated w by the exact power of two 2^(−fⱼ) that
// brings its largest |element| into [1, 2), so that w₂ = S·m·C. Then m⁻¹ = C·w₂⁻¹·S.
func colEquilibrate(w []float64, f []int, n int) {
	for j := range n {
		mx := 0.0
		for i := range n {
			mx = max(mx, math.Abs(w[i*n+j]))
		}
		f[j] = 0
		if mx == 0 {
			continue
		}
		_, f[j] = pow2Scale(mx)
		for i := range n {
			w[i*n+j] = math.Ldexp(w[i*n+j], -f[j])
		}
	}
}

// unscaleInverse writes m⁻¹ = C·w₂⁻¹·S into dst: element (i, j) of w₂⁻¹ (in iw) times 2^(−fᵢ−eⱼ).
func unscaleInverse(dst, iw []float64, f, e []int, n int) error {
	for i := range n {
		for j := range n {
			x := math.Ldexp(iw[i*n+j], -f[i]-e[j])
			if !reduce.AllFinite(x) {
				return errOverflow("m⁻¹")
			}
			dst[i*n+j] = x
		}
	}
	return nil
}

// Mat2Trace returns m₀₀ + m₁₁, the sum of the eigenvalues.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum is outside the float64 range.
//
// ID:      NUM-LINALG-062
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues)
// Tags:    matrix, 2x2, trace, linear-algebra
// Units:   the unit of m
// Example: [1 2; 3 4] → 5
func Mat2Trace(m Mat2) (float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return 0, errNotFinite("m")
	}
	t := m[0] + m[3]
	if !reduce.AllFinite(t) {
		return 0, errOverflow("trace")
	}
	return t, nil
}

// Mat2Rotation returns the matrix that turns a column vector counterclockwise by angleRad:
// [cos −sin; sin cos].
// Errors: ErrNotFinite when angleRad is NaN or ±Inf.
//
// ID:      NUM-LINALG-063
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (16) (the upper-left 2×2 block of R₃, transposed: Diebel's matrices are frame rotations)
// Tags:    matrix, 2x2, rotation, linear-algebra
// Units:   angleRad in radians; the result is dimensionless
// Example: π/2 → [0 −1; 1 0] (to within cos(π/2) ≈ 6.1e−17)
func Mat2Rotation(angleRad float64) (Mat2, error) {
	if !reduce.AllFinite(angleRad) {
		return Mat2{}, errNotFinite("angleRad")
	}
	s, c := math.Sincos(angleRad)
	return Mat2{c, -s, s, c}, nil
}
