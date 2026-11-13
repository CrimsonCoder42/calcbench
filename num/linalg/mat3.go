package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Mat3 is a 3×3 matrix stored row-major: [m00, m01, m02, m10, …, m22]. It is a value: functions take
// and return copies and never allocate. Vectors are columns: Mat3MulVec computes M·v.
type Mat3 [9]float64

func row3(m *Mat3, i int) Vec3 { return Vec3{m[3*i], m[3*i+1], m[3*i+2]} }

// Mat3Identity returns the 3×3 identity matrix.
//
// ID:      NUM-LINALG-065
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.5 "Inverse Matrices" (the identity I)
// Tags:    matrix, 3x3, identity, linear-algebra
// Units:   dimensionless
// Example: → [1 0 0; 0 1 0; 0 0 1]
func Mat3Identity() Mat3 { return Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1} }

// Mat3Mul returns the matrix product a·b, each element a three-term dot product accumulated with fused
// multiply-adds.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element of the product is outside
// the float64 range.
//
// ID:      NUM-LINALG-066
// Time:    O(1)   27 multiply-adds
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (row times column)
// Tags:    matrix, 3x3, multiply, product, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: [1 2 3; 4 5 6; 7 8 9]·I → the same matrix; [1 2 0; 0 1 0; 0 0 2]·[1 0 0; 3 1 0; 0 0 1] → [7 2 0; 3 1 0; 0 0 2] (hand-checked)
func Mat3Mul(a, b Mat3) (Mat3, error) {
	if !reduce.AllFinite(a[:]...) {
		return Mat3{}, errNotFinite("a")
	}
	if !reduce.AllFinite(b[:]...) {
		return Mat3{}, errNotFinite("b")
	}
	var c Mat3
	for i := range 3 {
		for j := range 3 {
			c[3*i+j] = math.FMA(a[3*i], b[j], math.FMA(a[3*i+1], b[3+j], a[3*i+2]*b[6+j]))
		}
	}
	if !reduce.AllFinite(c[:]...) {
		return Mat3{}, errOverflow("a·b")
	}
	return c, nil
}

// Mat3MulVec returns the matrix-vector product m·v (v a column).
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when an element is outside the float64 range.
//
// ID:      NUM-LINALG-067
// Time:    O(1)   9 multiply-adds
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.4 "Rules for Matrix Operations" (matrix times vector)
// Tags:    matrix, 3x3, matrix-vector, multiply, linear-algebra
// Units:   m in unit A, v in unit B; the result in A·B
// Example: [1 2 3; 4 5 6; 7 8 9]·(1, 0, −1) → (−2, −2, −2) (hand-checked)
func Mat3MulVec(m Mat3, v Vec3) (Vec3, error) {
	if !reduce.AllFinite(m[:]...) {
		return Vec3{}, errNotFinite("m")
	}
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	r := mulVec3(&m, v)
	if !finite3(r) {
		return Vec3{}, errOverflow("m·v")
	}
	return r, nil
}

func mulVec3(m *Mat3, v Vec3) Vec3 {
	return Vec3{dot3(row3(m, 0), v), dot3(row3(m, 1), v), dot3(row3(m, 2), v)}
}

// Mat3Transpose returns mᵀ.
// Errors: ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-068
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §2.7 "Transposes and Permutations"
// Tags:    matrix, 3x3, transpose, linear-algebra
// Units:   the unit of m
// Example: [1 2 3; 4 5 6; 7 8 9] → [1 4 7; 2 5 8; 3 6 9]
func Mat3Transpose(m Mat3) (Mat3, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat3{}, errNotFinite("m")
	}
	return transpose3(m), nil
}

func transpose3(m Mat3) Mat3 {
	return Mat3{m[0], m[3], m[6], m[1], m[4], m[7], m[2], m[5], m[8]}
}

// Mat3Det returns det m as the triple product r₀·(r₁ × r₂) of its rows: the cross product's components
// each within 2 ulps (Kahan's algorithm), the final dot product as if in twice the working precision
// (Dot2), so an exactly singular matrix gives 0. When an intermediate overflows, the rows are rescaled by
// exact powers of two and the determinant recomputed.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the determinant is outside the
// float64 range.
//
// ID:      NUM-LINALG-069
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (the triple product is the determinant)
// Tags:    matrix, 3x3, determinant, triple-product, linear-algebra
// Units:   m in unit A; the result in A³
// Example: [1 2 3; 0 1 4; 5 6 0] → 1 (hand-checked: 1·(−24) − 2·(−20) + 3·(−5))
func Mat3Det(m Mat3) (float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return 0, errNotFinite("m")
	}
	d := tripleDet(row3(&m, 0), row3(&m, 1), row3(&m, 2))
	if reduce.AllFinite(d) {
		return d, nil
	}
	// An intermediate overflowed: det m = 2^(e₀+e₁+e₂)·det w with each row of w scaled into [1, 2).
	var w Mat3
	var e [3]int
	if !rowEquilibrate(m[:], w[:], e[:], 3) {
		return 0, nil
	}
	d = math.Ldexp(tripleDet(row3(&w, 0), row3(&w, 1), row3(&w, 2)), e[0]+e[1]+e[2])
	if !reduce.AllFinite(d) {
		return 0, errOverflow("det m")
	}
	return d, nil
}

// tripleDet returns r₀·(r₁ × r₂) with the final sum by Dot2.
func tripleDet(r0, r1, r2 Vec3) float64 {
	c := cross3(r1, r2)
	var d dot2
	d.add(r0[0], c[0])
	d.add(r0[1], c[1])
	d.add(r0[2], c[2])
	return d.result()
}

// Mat3Inverse returns m⁻¹ = Cᵀ/det m (the cofactor formula): column j of the inverse is the cross
// product of the two rows other than j, in cyclic order, divided by det m. Each row and then each column is first
// scaled by an exact power of two (equilibration), so the cofactors and the determinant neither overflow nor
// underflow for a badly scaled but well-conditioned matrix. An ill-conditioned matrix has an inaccurate
// inverse; check with MatConditionEstimate1.
// Errors: ErrSingular when det m = 0; ErrNotFinite for a NaN or ±Inf element; ErrOverflow when an element
// of the inverse is outside the float64 range.
//
// ID:      NUM-LINALG-070
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A)
// Tags:    matrix, 3x3, inverse, cofactor, adjugate, linear-algebra
// Units:   m in unit A; the result in 1/A
// Example: [1 2 3; 0 1 4; 5 6 0] → [−24 18 5; 20 −15 −4; −5 4 1] (hand-checked: det = 1, so the inverse is Cᵀ)
func Mat3Inverse(m Mat3) (Mat3, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat3{}, errNotFinite("m")
	}
	var w Mat3
	var e [3]int
	var f [3]int
	if !rowEquilibrate(m[:], w[:], e[:], 3) {
		return Mat3{}, fmt.Errorf("%w: a row is zero", ErrSingular)
	}
	colEquilibrate(w[:], f[:], 3)
	r0, r1, r2 := row3(&w, 0), row3(&w, 1), row3(&w, 2)
	c0, c1, c2 := cross3(r1, r2), cross3(r2, r0), cross3(r0, r1)
	d := dot3(r0, c0)
	if d == 0 {
		return Mat3{}, fmt.Errorf("%w: det m = 0", ErrSingular)
	}
	// The columns of w⁻¹ are c₀/d, c₁/d, c₂/d.
	iw := Mat3{
		c0[0] / d, c1[0] / d, c2[0] / d,
		c0[1] / d, c1[1] / d, c2[1] / d,
		c0[2] / d, c1[2] / d, c2[2] / d,
	}
	var r Mat3
	if err := unscaleInverse(r[:], iw[:], f[:], e[:], 3); err != nil {
		return Mat3{}, err
	}
	return r, nil
}

// Mat3Trace returns m₀₀ + m₁₁ + m₂₂, the sum of the eigenvalues.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when the sum is outside the float64 range.
//
// ID:      NUM-LINALG-071
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (the trace equals the sum of the eigenvalues)
// Tags:    matrix, 3x3, trace, linear-algebra
// Units:   the unit of m
// Example: [1 2 3; 4 5 6; 7 8 9] → 15
func Mat3Trace(m Mat3) (float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return 0, errNotFinite("m")
	}
	t := m[0] + m[4] + m[8]
	if !reduce.AllFinite(t) {
		return 0, errOverflow("trace")
	}
	return t, nil
}
