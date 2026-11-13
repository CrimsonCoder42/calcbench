package linalg

import (
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// dot3 returns a·b with fused multiply-adds (two roundings fewer than the plain sum of products).
func dot3(a, b Vec3) float64 {
	return math.FMA(a[0], b[0], math.FMA(a[1], b[1], a[2]*b[2]))
}

// cross3 returns a × b; each component is a 2×2 determinant by Kahan's algorithm.
func cross3(a, b Vec3) Vec3 {
	return Vec3{
		det2(a[1], a[2], b[1], b[2]),
		det2(a[2], a[0], b[2], b[0]),
		det2(a[0], a[1], b[0], b[1]),
	}
}

// Vec3Dot returns the dot product a·b = a₀b₀ + a₁b₁ + a₂b₂, accumulated with fused multiply-adds.
// For an ill-conditioned dot product (heavy cancellation) use VecDotCompensated.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the result is outside the float64
// range.
//
// ID:      NUM-LINALG-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Tags:    vector, 3d, dot-product, inner-product, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: (1, 2, 3)·(4, −5, 6) → 12 (hand-checked)
func Vec3Dot(a, b Vec3) (float64, error) {
	if !finite3(a) {
		return 0, errNotFinite("a")
	}
	if !finite3(b) {
		return 0, errNotFinite("b")
	}
	d := dot3(a, b)
	if !reduce.AllFinite(d) {
		return 0, errOverflow("a·b")
	}
	return d, nil
}

// Vec3Cross returns the cross product a × b, perpendicular to both by the right-hand rule, with length
// |a||b| sin θ. Each component is a 2×2 determinant computed by Kahan's fused multiply-add algorithm,
// within 2 ulps.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a component is outside the float64
// range.
//
// ID:      NUM-LINALG-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (cross product)
// Source:  Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8
// Tags:    vector, 3d, cross-product, right-hand-rule, normal, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: x̂ × ŷ → ẑ; (1, 2, 3) × (4, 5, 6) → (−3, 6, −3) (hand-checked)
func Vec3Cross(a, b Vec3) (Vec3, error) {
	if !finite3(a) {
		return Vec3{}, errNotFinite("a")
	}
	if !finite3(b) {
		return Vec3{}, errNotFinite("b")
	}
	c := cross3(a, b)
	if !finite3(c) {
		return Vec3{}, errOverflow("a × b")
	}
	return c, nil
}

// Vec3ScalarTripleProduct returns a·(b × c), the determinant of the matrix with rows a, b, c: the
// signed volume of the parallelepiped they span (positive when a, b, c are right-handed). It is the
// "mixed product".
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the result is outside the float64
// range.
//
// ID:      NUM-LINALG-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (triple product = determinant = volume)
// Tags:    vector, 3d, triple-product, mixed-product, determinant, volume, linear-algebra
// Units:   a, b, c in units A, B, C; the result in A·B·C
// Example: x̂, ŷ, ẑ → 1; (1, 2, 3), (4, 5, 6), (7, 8, 10) → −3 (hand-checked)
func Vec3ScalarTripleProduct(a, b, c Vec3) (float64, error) {
	if !finite3(a) {
		return 0, errNotFinite("a")
	}
	if !finite3(b) {
		return 0, errNotFinite("b")
	}
	if !finite3(c) {
		return 0, errNotFinite("c")
	}
	bc := cross3(b, c)
	if !finite3(bc) {
		return 0, errOverflow("b × c")
	}
	d := dot3(a, bc)
	if !reduce.AllFinite(d) {
		return 0, errOverflow("a·(b × c)")
	}
	return d, nil
}

// Vec3VectorTripleProduct returns a × (b × c). It lies in the plane of b and c and equals
// b(a·c) − c(a·b); it is computed here as two cross products, each component within 2 ulps of the
// rounded inner product.
// Errors: as Vec3ScalarTripleProduct.
//
// ID:      NUM-LINALG-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (cross product, applied twice)
// Tags:    vector, 3d, triple-product, vector-triple-product, cross-product, linear-algebra
// Units:   a, b, c in units A, B, C; the result in A·B·C
// Example: a = x̂, b = x̂, c = ŷ → x̂ × ẑ = (0, −1, 0) (hand-checked)
func Vec3VectorTripleProduct(a, b, c Vec3) (Vec3, error) {
	if !finite3(a) {
		return Vec3{}, errNotFinite("a")
	}
	if !finite3(b) {
		return Vec3{}, errNotFinite("b")
	}
	if !finite3(c) {
		return Vec3{}, errNotFinite("c")
	}
	bc := cross3(b, c)
	if !finite3(bc) {
		return Vec3{}, errOverflow("b × c")
	}
	r := cross3(a, bc)
	if !finite3(r) {
		return Vec3{}, errOverflow("a × (b × c)")
	}
	return r, nil
}
