package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Mat4 is a 4×4 matrix stored row-major: [m00, m01, m02, m03, m10, …, m33]. It is a value: functions take
// and return copies and never allocate. Vectors are columns: Mat4MulVec computes M·v.
type Mat4 [16]float64

// Vec4 is a 4-component vector, used for homogeneous coordinates (x, y, z, w).
type Vec4 [4]float64

// Mat4Identity returns the 4×4 identity matrix.
//
// ID:      NUM-LINALG-077
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I)
// Tags:    matrix, 4x4, identity, homogeneous, linear-algebra
// Units:   dimensionless
// Example: → the 4×4 identity
func Mat4Identity() Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

func dot4(a0, a1, a2, a3, b0, b1, b2, b3 float64) float64 {
	return math.FMA(a0, b0, math.FMA(a1, b1, math.FMA(a2, b2, a3*b3)))
}

// Mat4Mul returns the matrix product a·b, each element a four-term dot product accumulated with fused
// multiply-adds. Composing transforms: (a·b)·v applies b first, then a.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the product is outside
// the float64 range.
//
// ID:      NUM-LINALG-078
// Time:    O(1)   64 multiply-adds
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column)
// Tags:    matrix, 4x4, multiply, product, homogeneous, transform, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: translation by (1, 2, 3) times translation by (4, 5, 6) → translation by (5, 7, 9) (hand-checked)
func Mat4Mul(a, b Mat4) (Mat4, error) {
	if !reduce.AllFinite(a[:]...) {
		return Mat4{}, errNotFinite("a")
	}
	if !reduce.AllFinite(b[:]...) {
		return Mat4{}, errNotFinite("b")
	}
	var c Mat4
	for i := range 4 {
		r := a[4*i : 4*i+4]
		for j := range 4 {
			c[4*i+j] = dot4(r[0], r[1], r[2], r[3], b[j], b[4+j], b[8+j], b[12+j])
		}
	}
	if !reduce.AllFinite(c[:]...) {
		return Mat4{}, errOverflow("a·b")
	}
	return c, nil
}

// Mat4MulVec returns the matrix-vector product m·v (v a column).
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element is outside the float64 range.
//
// ID:      NUM-LINALG-079
// Time:    O(1)   16 multiply-adds
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector)
// Tags:    matrix, 4x4, matrix-vector, multiply, homogeneous, linear-algebra
// Units:   m in unit A, v in unit B; the result in A·B
// Example: translation by (1, 2, 3) times (0, 0, 0, 1) → (1, 2, 3, 1) (hand-checked)
func Mat4MulVec(m Mat4, v Vec4) (Vec4, error) {
	if !reduce.AllFinite(m[:]...) {
		return Vec4{}, errNotFinite("m")
	}
	if !reduce.AllFinite(v[:]...) {
		return Vec4{}, errNotFinite("v")
	}
	var r Vec4
	for i := range 4 {
		r[i] = dot4(m[4*i], m[4*i+1], m[4*i+2], m[4*i+3], v[0], v[1], v[2], v[3])
	}
	if !reduce.AllFinite(r[:]...) {
		return Vec4{}, errOverflow("m·v")
	}
	return r, nil
}

// Mat4Transpose returns mᵀ.
// Errors: ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-080
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations"
// Tags:    matrix, 4x4, transpose, linear-algebra
// Units:   the unit of m
// Example: the transpose of a translation matrix moves the offset from the last column to the last row
func Mat4Transpose(m Mat4) (Mat4, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat4{}, errNotFinite("m")
	}
	var t Mat4
	for i := range 4 {
		for j := range 4 {
			t[4*j+i] = m[4*i+j]
		}
	}
	return t, nil
}

// minors4 returns the six 2×2 minors of rows 0–1 (s) and of rows 2–3 (c) used by the Laplace expansion
// of a 4×4 determinant along its first two rows, each within 2 ulps (Kahan's algorithm).
func minors4(a *Mat4) (s, c [6]float64) {
	s[0] = det2(a[0], a[1], a[4], a[5])
	s[1] = det2(a[0], a[2], a[4], a[6])
	s[2] = det2(a[0], a[3], a[4], a[7])
	s[3] = det2(a[1], a[2], a[5], a[6])
	s[4] = det2(a[1], a[3], a[5], a[7])
	s[5] = det2(a[2], a[3], a[6], a[7])
	c[0] = det2(a[8], a[9], a[12], a[13])
	c[1] = det2(a[8], a[10], a[12], a[14])
	c[2] = det2(a[8], a[11], a[12], a[15])
	c[3] = det2(a[9], a[10], a[13], a[14])
	c[4] = det2(a[9], a[11], a[13], a[15])
	c[5] = det2(a[10], a[11], a[14], a[15])
	return s, c
}

// det4 is the Laplace expansion Σ ± sᵢ·c₅₋ᵢ (complementary 2×2 minors).
func det4(s, c *[6]float64) float64 {
	var d dot2 // Dot2, so that exact cancellation (a singular matrix) gives 0
	d.add(s[0], c[5])
	d.add(-s[1], c[4])
	d.add(s[2], c[3])
	d.add(s[3], c[2])
	d.add(-s[4], c[1])
	d.add(s[5], c[0])
	return d.result()
}

// Mat4Det returns det m by the Laplace expansion along the first two rows: the sum of the six products
// of complementary 2×2 minors, with signs. Each minor is within 2 ulps (Kahan's algorithm) and the sum is
// taken by Dot2. When an intermediate overflows, the rows are rescaled by exact powers of two and the
// determinant recomputed.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the determinant is outside the
// float64 range.
//
// ID:      NUM-LINALG-081
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors"
// Tags:    matrix, 4x4, determinant, laplace-expansion, linear-algebra
// Units:   m in unit A; the result in A⁴
// Example: diag(1, 2, 3, 4) → 24; a homogeneous translation → 1 (hand-checked)
func Mat4Det(m Mat4) (float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return 0, errNotFinite("m")
	}
	s, c := minors4(&m)
	d := det4(&s, &c)
	if reduce.AllFinite(d) {
		return d, nil
	}
	// An intermediate overflowed: det m = 2^(Σeᵢ)·det w with each row of w scaled into [1, 2).
	var w Mat4
	var e [4]int
	if !rowEquilibrate(m[:], w[:], e[:], 4) {
		return 0, nil
	}
	s, c = minors4(&w)
	d = math.Ldexp(det4(&s, &c), e[0]+e[1]+e[2]+e[3])
	if !reduce.AllFinite(d) {
		return 0, errOverflow("det m")
	}
	return d, nil
}

// Mat4Inverse returns m⁻¹ = Cᵀ/det m (the cofactor formula), building every cofactor from the twelve
// 2×2 minors of the Laplace expansion. Each row and then each column is first scaled by an exact power of two
// (equilibration), so the minors and the determinant neither overflow nor underflow for a badly scaled but
// well-conditioned matrix. A tiny determinant relative to the elements means an ill-conditioned matrix; check
// with MatConditionEstimate1, or use MatInverse (LU with partial pivoting).
// Errors: ErrSingular when det m = 0; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element
// of the inverse is outside the float64 range.
//
// ID:      NUM-LINALG-082
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A)
// Tags:    matrix, 4x4, inverse, cofactor, adjugate, homogeneous, linear-algebra
// Units:   m in unit A; the result in 1/A
// Example: translation by (1, 2, 3) → translation by (−1, −2, −3) (hand-checked)
func Mat4Inverse(m Mat4) (Mat4, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat4{}, errNotFinite("m")
	}
	var a Mat4
	var ex [4]int
	var f [4]int
	if !rowEquilibrate(m[:], a[:], ex[:], 4) {
		return Mat4{}, fmt.Errorf("%w: a row is zero", ErrSingular)
	}
	colEquilibrate(a[:], f[:], 4)
	s, c := minors4(&a)
	d := det4(&s, &c)
	if d == 0 {
		return Mat4{}, fmt.Errorf("%w: det m = 0", ErrSingular)
	}
	inv := Mat4{
		a[5]*c[5] - a[6]*c[4] + a[7]*c[3],
		-a[1]*c[5] + a[2]*c[4] - a[3]*c[3],
		a[13]*s[5] - a[14]*s[4] + a[15]*s[3],
		-a[9]*s[5] + a[10]*s[4] - a[11]*s[3],

		-a[4]*c[5] + a[6]*c[2] - a[7]*c[1],
		a[0]*c[5] - a[2]*c[2] + a[3]*c[1],
		-a[12]*s[5] + a[14]*s[2] - a[15]*s[1],
		a[8]*s[5] - a[10]*s[2] + a[11]*s[1],

		a[4]*c[4] - a[5]*c[2] + a[7]*c[0],
		-a[0]*c[4] + a[1]*c[2] - a[3]*c[0],
		a[12]*s[4] - a[13]*s[2] + a[15]*s[0],
		-a[8]*s[4] + a[9]*s[2] - a[11]*s[0],

		-a[4]*c[3] + a[5]*c[1] - a[6]*c[0],
		a[0]*c[3] - a[1]*c[1] + a[2]*c[0],
		-a[12]*s[3] + a[13]*s[1] - a[14]*s[0],
		a[8]*s[3] - a[9]*s[1] + a[10]*s[0],
	}
	for i := range inv {
		inv[i] /= d
	}
	var r Mat4
	if err := unscaleInverse(r[:], inv[:], f[:], ex[:], 4); err != nil {
		return Mat4{}, err
	}
	return r, nil
}

// Mat4Trace returns the sum of the diagonal elements.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum is outside the float64 range.
//
// ID:      NUM-LINALG-083
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues)
// Tags:    matrix, 4x4, trace, linear-algebra
// Units:   the unit of m
// Example: diag(1, 2, 3, 4) → 10
func Mat4Trace(m Mat4) (float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return 0, errNotFinite("m")
	}
	t := m[0] + m[5] + m[10] + m[15]
	if !reduce.AllFinite(t) {
		return 0, errOverflow("trace")
	}
	return t, nil
}
