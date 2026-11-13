package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// QuatFromAxisAngle returns the unit quaternion [cos(θ/2), sin(θ/2)·k] (k = axis/|axis|) of the rotation
// by angleRad about axis, right-hand rule. This is the "quaternion of rotation".
// Errors: ErrZeroVector when axis = 0; ErrNotFinite for a NaN or ±Inf input.
//
// ID:      NUM-LINALG-108
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.12, eq. (175)
// Tags:    quaternion, rotation, axis-angle, versor, linear-algebra
// Units:   axis in any unit, angleRad in radians; the result is dimensionless
// Example: axis ẑ, π/2 → (√2/2, 0, 0, √2/2) (hand-checked)
func QuatFromAxisAngle(axis Vec3, angleRad float64) (Quat, error) {
	if !finite3(axis) {
		return Quat{}, errNotFinite("axis")
	}
	if !reduce.AllFinite(angleRad) {
		return Quat{}, errNotFinite("angleRad")
	}
	k, ok := unit3(axis)
	if !ok {
		return Quat{}, fmt.Errorf("%w: axis", ErrZeroVector)
	}
	s, c := math.Sincos(angleRad / 2)
	return Quat{c, s * k[0], s * k[1], s * k[2]}, nil
}

// QuatToAxisAngle returns the unit axis and the angle in [0, π] radians of the rotation represented by
// q (normalised first; q and −q give the same answer). The angle is 2·atan2(|q⃗|, |W|), accurate for
// small angles where Diebel's 2·acos(W) loses half its digits. For the identity rotation the axis is
// undefined and x̂ = (1, 0, 0) is returned with angle 0.
// Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-109
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.13, eqs. (198)–(199) (with acos replaced by the equivalent atan2 of eq. (175)'s sine and cosine)
// Tags:    quaternion, rotation, axis-angle, linear-algebra
// Units:   dimensionless quaternion; the axis is dimensionless, the angle in radians
// Example: (√2/2, 0, 0, √2/2) → axis ẑ, angle π/2 (hand-checked)
func QuatToAxisAngle(q Quat) (axis Vec3, angleRad float64, err error) {
	if !finiteQ(q) {
		return Vec3{}, 0, errNotFinite("q")
	}
	u, ok := unitQ(q)
	if !ok {
		return Vec3{}, 0, fmt.Errorf("%w: q", ErrZeroVector)
	}
	if u.W < 0 {
		u = Quat{-u.W, -u.X, -u.Y, -u.Z}
	}
	v := Vec3{u.X, u.Y, u.Z}
	s := norm3(v)
	if s == 0 {
		return Vec3{1, 0, 0}, 0, nil
	}
	return Vec3{v[0] / s, v[1] / s, v[2] / s}, 2 * math.Atan2(s, u.W), nil
}
