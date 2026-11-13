package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// minor3 returns the 2×2 minor of m with row i and column j deleted (0-based), within 2 ulps (det2).
func minor3(m *Mat3, i, j int) float64 {
	var r [2]int
	var c [2]int
	for k, n := 0, 0; k < 3; k++ {
		if k != i {
			r[n] = k
			n++
		}
	}
	for k, n := 0, 0; k < 3; k++ {
		if k != j {
			c[n] = k
			n++
		}
	}
	return det2(m[3*r[0]+c[0]], m[3*r[0]+c[1]], m[3*r[1]+c[0]], m[3*r[1]+c[1]])
}

// Mat3Minor returns the minor M_ij of a 3×3 matrix: the determinant of the 2×2 matrix left after deleting
// row i and column j (0-based). The cofactor is (−1)^(i+j)·M_ij.
// Errors: ErrOutOfRange when i or j is outside 0–2; ErrNotFinite for a NaN or ±Inf element; ErrOverflow
// when the minor is outside the float64 range.
//
// ID:      NUM-LINALG-123
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" (the cofactor C_ij = (−1)^(i+j) det M_ij, M_ij with row i and column j removed)
// Tags:    matrix, 3x3, minor, cofactor, determinant, linear-algebra
// Units:   m in unit A; the result in A²
// Example: m = [1 2 3; 4 5 6; 7 8 10], i = 0, j = 0 → 5·10 − 6·8 = 2 (hand-checked)
func Mat3Minor(m Mat3, i, j int) (float64, error) {
	if i < 0 || i > 2 || j < 0 || j > 2 {
		return 0, fmt.Errorf("%w: (i, j) = (%d, %d), want 0–2", ErrOutOfRange, i, j)
	}
	if !reduce.AllFinite(m[:]...) {
		return 0, errNotFinite("m")
	}
	d := minor3(&m, i, j)
	if !reduce.AllFinite(d) {
		return 0, errOverflow("minor")
	}
	return d, nil
}

// Mat2Adjugate returns the adjugate (classical adjoint) adj m = Cᵀ = [d −b; −c a] of m = [a b; c d], the
// matrix with m·adj m = adj m·m = det(m)·I. It exists for singular m too.
// Errors: ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-124
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A; the 2×2 case [d −b; −c a])
// Tags:    matrix, 2x2, adjugate, cofactor, classical-adjoint, linear-algebra
// Units:   m in unit A; the result in A
// Example: [1 2; 3 4] → [4 −2; −3 1] (hand-checked)
func Mat2Adjugate(m Mat2) (Mat2, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat2{}, errNotFinite("m")
	}
	return Mat2{m[3], -m[1], -m[2], m[0]}, nil
}

// Mat3Adjugate returns the adjugate adj m = Cᵀ of a 3×3 matrix, each cofactor a 2×2 determinant within 2
// ulps (Kahan's algorithm); m·adj m = det(m)·I, and adj m exists for singular m.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a cofactor is outside the float64 range.
//
// ID:      NUM-LINALG-125
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A)
// Tags:    matrix, 3x3, adjugate, cofactor, classical-adjoint, linear-algebra
// Units:   m in unit A; the result in A²
// Example: [1 2 3; 0 1 4; 5 6 0] → [−24 18 5; 20 −15 −4; −5 4 1] (hand-checked; det = 1, so this is the inverse)
func Mat3Adjugate(m Mat3) (Mat3, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat3{}, errNotFinite("m")
	}
	var adj Mat3
	for i := range 3 {
		for j := range 3 {
			c := minor3(&m, i, j)
			if (i+j)%2 == 1 {
				c = -c
			}
			adj[3*j+i] = c // transpose
		}
	}
	if !reduce.AllFinite(adj[:]...) {
		return Mat3{}, errOverflow("cofactor")
	}
	return adj, nil
}

// Mat4Adjugate returns the adjugate adj m = Cᵀ of a 4×4 matrix, each cofactor a 3×3 determinant (a scalar
// triple product summed by Dot2); m·adj m = det(m)·I, and adj m exists for singular m.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a cofactor is outside the float64 range.
//
// ID:      NUM-LINALG-126
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.2 "Permutations and Cofactors" and §5.3 "Cramer's Rule, Inverses, and Volumes" (A⁻¹ = Cᵀ/det A)
// Tags:    matrix, 4x4, adjugate, cofactor, classical-adjoint, linear-algebra
// Units:   m in unit A; the result in A³
// Example: diag(1, 2, 3, 4) → diag(24, 12, 8, 6) (hand-checked)
func Mat4Adjugate(m Mat4) (Mat4, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat4{}, errNotFinite("m")
	}
	var adj Mat4
	for i := range 4 {
		for j := range 4 {
			var sub Mat3
			k := 0
			for r := range 4 {
				if r == i {
					continue
				}
				for c := range 4 {
					if c == j {
						continue
					}
					sub[k] = m[4*r+c]
					k++
				}
			}
			d := tripleDet(row3(&sub, 0), row3(&sub, 1), row3(&sub, 2))
			if (i+j)%2 == 1 {
				d = -d
			}
			adj[4*j+i] = d
		}
	}
	if !reduce.AllFinite(adj[:]...) {
		return Mat4{}, errOverflow("cofactor")
	}
	return adj, nil
}

// Mat2SolveCramer solves the 2×2 system m·x = b by Cramer's rule, xⱼ = det(Mⱼ)/det(m) with Mⱼ = m whose
// column j is replaced by b, each determinant within 2 ulps (Kahan's algorithm). It is the closed form for
// a tiny system and needs no allocation; for an ill-conditioned or larger system MatSolve (LU with partial
// pivoting, NUM-LINALG-096) is the stable general choice.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrSingular when det m = 0; ErrOverflow when a component is
// outside the float64 range.
//
// ID:      NUM-LINALG-144
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (xⱼ = det Bⱼ / det A)
// Tags:    matrix, 2x2, cramer, solve, linear-system, linear-algebra
// Units:   m in unit A, b in unit B; x in B/A
// Example: 2x + y = 5, x + 3y = 10 → x = (1, 3) (hand-checked)
func Mat2SolveCramer(m Mat2, b Vec2) (Vec2, error) {
	if !reduce.AllFinite(m[:]...) || !finite2(b) {
		return Vec2{}, errNotFinite("m or b")
	}
	d := det2(m[0], m[1], m[2], m[3])
	if d == 0 {
		return Vec2{}, fmt.Errorf("%w: det m = 0", ErrSingular)
	}
	x := Vec2{det2(b[0], m[1], b[1], m[3]) / d, det2(m[0], b[0], m[2], b[1]) / d}
	if !finite2(x) {
		return Vec2{}, errOverflow("x")
	}
	return x, nil
}

// Mat3SolveCramer solves the 3×3 system m·x = b by Cramer's rule, each determinant a scalar triple product
// summed by Dot2. It is the closed form for a tiny system and needs no allocation; MatSolve (NUM-LINALG-096)
// is the stable general choice when m is ill-conditioned.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrSingular when det m = 0; ErrOverflow when a component is
// outside the float64 range.
//
// ID:      NUM-LINALG-145
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (xⱼ = det Bⱼ / det A)
// Tags:    matrix, 3x3, cramer, solve, linear-system, linear-algebra
// Units:   m in unit A, b in unit B; x in B/A
// Example: x + y + z = 6, 2y + 5z = −4, 2x + 5y − z = 27 → x = (5, 3, −2) (hand-checked)
func Mat3SolveCramer(m Mat3, b Vec3) (Vec3, error) {
	if !reduce.AllFinite(m[:]...) || !finite3(b) {
		return Vec3{}, errNotFinite("m or b")
	}
	d := tripleDet(row3(&m, 0), row3(&m, 1), row3(&m, 2))
	if d == 0 {
		return Vec3{}, fmt.Errorf("%w: det m = 0", ErrSingular)
	}
	var x Vec3
	for j := range 3 {
		mj := m
		for i := range 3 {
			mj[3*i+j] = b[i]
		}
		x[j] = tripleDet(row3(&mj, 0), row3(&mj, 1), row3(&mj, 2)) / d
	}
	if !finite3(x) {
		return Vec3{}, errOverflow("x")
	}
	return x, nil
}

// Mat2CharPoly returns the characteristic polynomial det(λI − m) = λ² − tr(m)·λ + det(m) of a 2×2 matrix
// as ascending coefficients (D-006): [det m, −tr m, 1].
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a coefficient is outside the float64
// range.
//
// ID:      NUM-LINALG-146
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (det(A − λI) = λ² − (trace)λ + det for 2×2)
// Tags:    matrix, 2x2, characteristic-polynomial, eigenvalues, trace, determinant, linear-algebra
// Units:   m in unit A; coefficient k in A^(2−k)
// Example: [4 1; 2 3] → λ² − 7λ + 10, i.e. [10, −7, 1] (hand-checked)
func Mat2CharPoly(m Mat2) ([3]float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return [3]float64{}, errNotFinite("m")
	}
	p := [3]float64{det2(m[0], m[1], m[2], m[3]), -(m[0] + m[3]), 1}
	if !reduce.AllFinite(p[:]...) {
		return [3]float64{}, errOverflow("coefficient")
	}
	return p, nil
}

// Mat3CharPoly returns the characteristic polynomial det(λI − m) = λ³ − tr(m)·λ² + E₂·λ − det(m) of a 3×3
// matrix as ascending coefficients (D-006): [−det m, E₂, −tr m, 1], where E₂ is the sum of the three
// principal 2×2 minors.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when a coefficient is outside the float64
// range.
//
// ID:      NUM-LINALG-147
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §6.1 "Introduction to Eigenvalues" (det(A − λI) = 0; returned here as the monic det(λI − A) = −det(A − λI), same roots: the λ² coefficient is −trace, the λ coefficient the sum of the principal 2×2 minors, the constant −det)
// Tags:    matrix, 3x3, characteristic-polynomial, eigenvalues, principal-minors, linear-algebra
// Units:   m in unit A; coefficient k in A^(3−k)
// Example: diag(1, 2, 3) → (λ − 1)(λ − 2)(λ − 3) = λ³ − 6λ² + 11λ − 6, i.e. [−6, 11, −6, 1] (hand-checked)
func Mat3CharPoly(m Mat3) ([4]float64, error) {
	if !reduce.AllFinite(m[:]...) {
		return [4]float64{}, errNotFinite("m")
	}
	var e2 dot2
	e2.add(1, minor3(&m, 0, 0))
	e2.add(1, minor3(&m, 1, 1))
	e2.add(1, minor3(&m, 2, 2))
	det := tripleDet(row3(&m, 0), row3(&m, 1), row3(&m, 2))
	p := [4]float64{-det, e2.result(), -(m[0] + m[4] + m[8]), 1}
	if !reduce.AllFinite(p[:]...) {
		return [4]float64{}, errOverflow("coefficient")
	}
	return p, nil
}
