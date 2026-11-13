package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// The rotation matrices here are ACTIVE: applied to a column vector they turn the vector
// counterclockwise about the axis as seen from the axis' positive end (right-hand rule). Diebel 2006
// writes frame (passive) rotations, which are the transposes of these.

// Mat3RotationX returns the matrix that turns a vector by angleRad about the +x axis:
// [1 0 0; 0 cos −sin; 0 sin cos].
// Errors: ErrNotFinite when angleRad is NaN or ±Inf.
//
// ID:      NUM-LINALG-072
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (14) (R₁, transposed: Diebel's matrices are frame rotations)
// Tags:    matrix, 3x3, rotation, rotation-x, roll, linear-algebra
// Units:   angleRad in radians; the result is dimensionless
// Example: π/2 turns ŷ into ẑ
func Mat3RotationX(angleRad float64) (Mat3, error) {
	if !reduce.AllFinite(angleRad) {
		return Mat3{}, errNotFinite("angleRad")
	}
	s, c := math.Sincos(angleRad)
	return Mat3{1, 0, 0, 0, c, -s, 0, s, c}, nil
}

// Mat3RotationY returns the matrix that turns a vector by angleRad about the +y axis:
// [cos 0 sin; 0 1 0; −sin 0 cos].
// Errors: ErrNotFinite when angleRad is NaN or ±Inf.
//
// ID:      NUM-LINALG-073
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (15) (R₂, transposed: Diebel's matrices are frame rotations)
// Tags:    matrix, 3x3, rotation, rotation-y, pitch, linear-algebra
// Units:   angleRad in radians; the result is dimensionless
// Example: π/2 turns ẑ into x̂
func Mat3RotationY(angleRad float64) (Mat3, error) {
	if !reduce.AllFinite(angleRad) {
		return Mat3{}, errNotFinite("angleRad")
	}
	s, c := math.Sincos(angleRad)
	return Mat3{c, 0, s, 0, 1, 0, -s, 0, c}, nil
}

// Mat3RotationZ returns the matrix that turns a vector by angleRad about the +z axis:
// [cos −sin 0; sin cos 0; 0 0 1].
// Errors: ErrNotFinite when angleRad is NaN or ±Inf.
//
// ID:      NUM-LINALG-074
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §2.4, eq. (16) (R₃, transposed: Diebel's matrices are frame rotations)
// Tags:    matrix, 3x3, rotation, rotation-z, yaw, linear-algebra
// Units:   angleRad in radians; the result is dimensionless
// Example: π/2 turns x̂ into ŷ
func Mat3RotationZ(angleRad float64) (Mat3, error) {
	if !reduce.AllFinite(angleRad) {
		return Mat3{}, errNotFinite("angleRad")
	}
	s, c := math.Sincos(angleRad)
	return Mat3{c, -s, 0, s, c, 0, 0, 0, 1}, nil
}

// Mat3RotationAxisAngle returns the matrix that turns a vector by angleRad about axis (right-hand rule):
// with k = axis/|axis|, R = cos θ·I + (1 − cos θ)·kkᵀ + sin θ·[k]×, where [k]× is the cross-product
// matrix (Rodrigues' formula in matrix form). 1 − cos θ is computed as 2 sin²(θ/2) for accuracy at small
// angles.
// Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf input.
//
// ID:      NUM-LINALG-075
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497, eq. (1) (rotation about the unit vector a; Strang multiplies row vectors, so the matrix here is its transpose)
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §7.7, eqs. (243)–(244) (transposed: Diebel's matrix is the frame rotation)
// Tags:    matrix, 3x3, rotation, axis-angle, rodrigues, linear-algebra
// Units:   axis in any unit, angleRad in radians; the result is dimensionless
// Example: axis ẑ, π/2 → [0 −1 0; 1 0 0; 0 0 1] (to within cos(π/2) ≈ 6.1e−17)
func Mat3RotationAxisAngle(axis Vec3, angleRad float64) (Mat3, error) {
	if !finite3(axis) {
		return Mat3{}, errNotFinite("axis")
	}
	if !reduce.AllFinite(angleRad) {
		return Mat3{}, errNotFinite("angleRad")
	}
	k, ok := unit3(axis)
	if !ok {
		return Mat3{}, fmt.Errorf("%w: axis", ErrZeroVector)
	}
	s, c := math.Sincos(angleRad)
	h := math.Sin(angleRad / 2)
	t := 2 * h * h
	x, y, z := k[0], k[1], k[2]
	return Mat3{
		math.FMA(t*x, x, c), math.FMA(t*x, y, -s*z), math.FMA(t*x, z, s*y),
		math.FMA(t*y, x, s*z), math.FMA(t*y, y, c), math.FMA(t*y, z, -s*x),
		math.FMA(t*z, x, -s*y), math.FMA(t*z, y, s*x), math.FMA(t*z, z, c),
	}, nil
}
