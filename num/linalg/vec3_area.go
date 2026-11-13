package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Vec3ParallelogramArea returns |a × b|, the area of the parallelogram spanned by a and b. In the
// plane, the area is |Vec2Cross(a, b)|.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the area exceeds the float64 range.
//
// ID:      NUM-LINALG-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (|u × v| is the area of the parallelogram)
// Tags:    vector, 3d, area, parallelogram, cross-product, linear-algebra
// Units:   a and b in one length unit; the result in that unit squared
// Example: (3, 0, 0) and (1, 2, 0) → 6 (hand-checked)
func Vec3ParallelogramArea(a, b Vec3) (float64, error) {
	if !finite3(a) {
		return 0, errNotFinite("a")
	}
	if !finite3(b) {
		return 0, errNotFinite("b")
	}
	c := cross3(a, b)
	if !finite3(c) {
		return 0, errOverflow("a × b")
	}
	n := norm3(c)
	if !reduce.AllFinite(n) {
		return 0, errOverflow("|a × b|")
	}
	return n, nil
}

// Vec3TriangleArea returns ½|a × b|, the area of the triangle with two sides a and b from a common
// vertex.
// Errors: as Vec3ParallelogramArea.
//
// ID:      NUM-LINALG-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (a triangle is half of the parallelogram)
// Tags:    vector, 3d, area, triangle, cross-product, linear-algebra
// Units:   a and b in one length unit; the result in that unit squared
// Example: (3, 0, 0) and (1, 2, 0) → 3 (hand-checked)
func Vec3TriangleArea(a, b Vec3) (float64, error) {
	if !finite3(a) {
		return 0, errNotFinite("a")
	}
	if !finite3(b) {
		return 0, errNotFinite("b")
	}
	c := cross3(a, b)
	if !finite3(c) {
		// The cross product overflowed; halve first (exact) and try again.
		c = cross3(Vec3{a[0] / 2, a[1] / 2, a[2] / 2}, b)
		if !finite3(c) {
			return 0, errOverflow("a × b")
		}
		n := norm3(c)
		if !reduce.AllFinite(n) {
			return 0, errOverflow("½|a × b|")
		}
		return n, nil
	}
	if n := norm3(c); reduce.AllFinite(n) {
		return n / 2, nil
	}
	// |a × b| exceeds the float64 range but half of it may not.
	n := norm3(Vec3{c[0] / 2, c[1] / 2, c[2] / 2})
	if !reduce.AllFinite(n) {
		return 0, errOverflow("½|a × b|")
	}
	return n, nil
}

// Vec3PlaneNormal returns the unit normal of the plane through p0, p1 and p2:
// (p1 − p0) × (p2 − p0) normalised. Its direction follows the right-hand rule: it points toward an
// observer who sees p0 → p1 → p2 counterclockwise.
// Errors: ErrZeroVector when the three points are collinear or coincide (no plane); ErrNotFinite for a
// NaN or ±Inf component; ErrOverflow when an edge vector is outside the float64 range.
//
// ID:      NUM-LINALG-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §5.3 "Cramer's Rule, Inverses, and Volumes" (u × v is perpendicular to u and v)
// Tags:    vector, 3d, plane, normal, cross-product, linear-algebra
// Units:   points in one length unit; the result is dimensionless
// Example: (0, 0, 0), (1, 0, 0), (0, 1, 0) → (0, 0, 1) (hand-checked)
func Vec3PlaneNormal(p0, p1, p2 Vec3) (Vec3, error) {
	if !finite3(p0) {
		return Vec3{}, errNotFinite("p0")
	}
	if !finite3(p1) {
		return Vec3{}, errNotFinite("p1")
	}
	if !finite3(p2) {
		return Vec3{}, errNotFinite("p2")
	}
	e1 := Vec3{p1[0] - p0[0], p1[1] - p0[1], p1[2] - p0[2]}
	e2 := Vec3{p2[0] - p0[0], p2[1] - p0[1], p2[2] - p0[2]}
	if !finite3(e1) || !finite3(e2) {
		return Vec3{}, errOverflow("edge vector")
	}
	// Scale each edge to unit size first (exact powers of two do not change the direction), so the
	// cross product cannot overflow.
	u1, ok1 := unit3(e1)
	u2, ok2 := unit3(e2)
	if !ok1 || !ok2 {
		return Vec3{}, fmt.Errorf("%w: coincident points", ErrZeroVector)
	}
	n, ok := unit3(cross3(u1, u2))
	if !ok {
		return Vec3{}, fmt.Errorf("%w: collinear points", ErrZeroVector)
	}
	return n, nil
}
