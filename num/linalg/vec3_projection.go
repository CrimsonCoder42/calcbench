package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// scaleV3 returns v·2^(−e) with its largest |component| brought into [1, 2) (e = 0 for v = 0), so that
// products and differences with v cannot overflow; results are scaled back by 2^e.
func scaleV3(v Vec3) (Vec3, int) {
	m := max(math.Abs(v[0]), math.Abs(v[1]), math.Abs(v[2]))
	if m == 0 {
		return v, 0
	}
	_, e := pow2Scale(m)
	return Vec3{math.Ldexp(v[0], -e), math.Ldexp(v[1], -e), math.Ldexp(v[2], -e)}, e
}

func unscaleV3(v Vec3, e int) Vec3 {
	return Vec3{math.Ldexp(v[0], e), math.Ldexp(v[1], e), math.Ldexp(v[2], e)}
}

// Vec3ComponentAlong returns the scalar component of v along axis, v·â with â = axis/|axis|: the
// signed length of v's projection onto the axis line.
// Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the
// result is outside the float64 range.
//
// ID:      NUM-LINALG-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (dot product with a unit vector)
// Tags:    vector, 3d, component, scalar-projection, linear-algebra
// Units:   v in any unit, axis in any unit; the result in v's unit
// Example: v = (3, 4, 5), axis = (0, 2, 0) → 4 (hand-checked)
func Vec3ComponentAlong(v, axis Vec3) (float64, error) {
	if !finite3(v) {
		return 0, errNotFinite("v")
	}
	if !finite3(axis) {
		return 0, errNotFinite("axis")
	}
	u, ok := unit3(axis)
	if !ok {
		return 0, fmt.Errorf("%w: axis", ErrZeroVector)
	}
	c := dot3(v, u)
	if !reduce.AllFinite(c) {
		return 0, errOverflow("v·â")
	}
	return c, nil
}

// Vec3Project returns the projection of v onto the line of onto: (v·n)n with n = onto/|onto|, the
// matrix nnᵀ applied to v.
// Errors: ErrZeroVector when onto = 0; ErrNotFinite for a NaN or ±Inf component; ErrOverflow when a
// component is outside the float64 range.
//
// ID:      NUM-LINALG-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §4.2 "Projections" (projection onto a line: p = a(aᵀb)/(aᵀa))
// Tags:    vector, 3d, projection, vector-projection, linear-algebra
// Units:   v in any unit, onto in any unit; the result in v's unit
// Example: v = (3, 4, 5), onto = (0, 2, 0) → (0, 4, 0) (hand-checked)
func Vec3Project(v, onto Vec3) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	if !finite3(onto) {
		return Vec3{}, errNotFinite("onto")
	}
	n, ok := unit3(onto)
	if !ok {
		return Vec3{}, fmt.Errorf("%w: onto", ErrZeroVector)
	}
	w, e := scaleV3(v)
	c := dot3(w, n)
	r := unscaleV3(Vec3{c * n[0], c * n[1], c * n[2]}, e)
	if !finite3(r) {
		return Vec3{}, errOverflow("projection")
	}
	return r, nil
}

// Vec3Reject returns the rejection of v from the line of from: v − (v·n)n with n = from/|from|, the
// component of v perpendicular to from (the projection I − nnᵀ onto the plane with normal n).
// Errors: as Vec3Project, with ErrZeroVector when from = 0.
//
// ID:      NUM-LINALG-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497 (projection onto the plane nᵀv = 0 is I − nnᵀ)
// Tags:    vector, 3d, rejection, perpendicular-component, projection, plane, linear-algebra
// Units:   v in any unit, from in any unit; the result in v's unit
// Example: v = (3, 4, 5), from = (0, 2, 0) → (3, 0, 5) (hand-checked)
func Vec3Reject(v, from Vec3) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	if !finite3(from) {
		return Vec3{}, errNotFinite("from")
	}
	n, ok := unit3(from)
	if !ok {
		return Vec3{}, fmt.Errorf("%w: from", ErrZeroVector)
	}
	w, e := scaleV3(v)
	c := dot3(w, n)
	r := unscaleV3(Vec3{w[0] - c*n[0], w[1] - c*n[1], w[2] - c*n[2]}, e)
	if !finite3(r) {
		return Vec3{}, errOverflow("rejection")
	}
	return r, nil
}

// Vec3Reflect returns the mirror image of v in the plane through the origin with the given normal:
// v − 2(v·n)n with n = normal/|normal|, the matrix I − 2nnᵀ.
// Errors: as Vec3Project, with ErrZeroVector when normal = 0.
//
// ID:      NUM-LINALG-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 498 (mirror matrix I − 2nnᵀ)
// Tags:    vector, 3d, reflection, mirror, householder, linear-algebra
// Units:   v in any unit, normal in any unit; the result in v's unit
// Example: v = (1, −1, 2), normal = (0, 3, 0) → (1, 1, 2) (hand-checked)
func Vec3Reflect(v, normal Vec3) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	if !finite3(normal) {
		return Vec3{}, errNotFinite("normal")
	}
	n, ok := unit3(normal)
	if !ok {
		return Vec3{}, fmt.Errorf("%w: normal", ErrZeroVector)
	}
	w, e := scaleV3(v)
	c := 2 * dot3(w, n)
	r := unscaleV3(Vec3{w[0] - c*n[0], w[1] - c*n[1], w[2] - c*n[2]}, e)
	if !finite3(r) {
		return Vec3{}, errOverflow("reflection")
	}
	return r, nil
}
