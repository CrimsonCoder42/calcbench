package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Vec3Rotate returns v turned by angleRad about axis through the origin, counterclockwise when the axis
// points at the viewer (right-hand rule). With k = axis/|axis| this is Rodrigues' formula
// v cos θ + (k × v) sin θ + k (k·v)(1 − cos θ), the active form of the axis-angle rotation matrix.
// Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a
// component of the result is outside the float64 range.
//
// ID:      NUM-LINALG-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497, eq. (1) (rotation about the unit vector a; Strang multiplies row vectors, so the column-vector matrix used here is its transpose)
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §7.7, eqs. (243)–(244) (transposed: Diebel's matrix is the frame rotation)
// Tags:    vector, 3d, rotate, rotation, rodrigues, axis-angle, linear-algebra
// Units:   v in any unit, axis in any unit, angleRad in radians; the result in v's unit
// Example: v = (1, 0, 0), axis = (0, 0, 1), angle π/2 → (0, 1, 0) (hand-checked)
func Vec3Rotate(v, axis Vec3, angleRad float64) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	if !finite3(axis) {
		return Vec3{}, errNotFinite("axis")
	}
	if !reduce.AllFinite(angleRad) {
		return Vec3{}, errNotFinite("angleRad")
	}
	k, ok := unit3(axis)
	if !ok {
		return Vec3{}, fmt.Errorf("%w: axis", ErrZeroVector)
	}
	r := rotate3(v, k, angleRad)
	if !finite3(r) {
		return Vec3{}, errOverflow("rotated v")
	}
	return r, nil
}

// rotate3 applies Rodrigues' formula for a unit axis k. 1 − cos θ is computed as 2 sin²(θ/2), which
// keeps its relative accuracy for small angles.
func rotate3(v, k Vec3, angleRad float64) Vec3 {
	s, c := math.Sincos(angleRad)
	h := math.Sin(angleRad / 2)
	oneMinusCos := 2 * h * h
	kv := cross3(k, v)
	kd := dot3(k, v) * oneMinusCos
	return Vec3{
		math.FMA(v[0], c, math.FMA(kv[0], s, k[0]*kd)),
		math.FMA(v[1], c, math.FMA(kv[1], s, k[1]*kd)),
		math.FMA(v[2], c, math.FMA(kv[2], s, k[2]*kd)),
	}
}
