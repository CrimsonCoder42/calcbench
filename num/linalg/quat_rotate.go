package linalg

import "fmt"

// QuatRotateVec3 returns v rotated by q: v′ = q v q⁻¹ with v taken as the pure quaternion [0, v] (active
// rotation, right-hand rule; q is normalised first, so any nonzero q works). For a unit q this is
// computed as t = 2(q⃗ × v), v′ = v + W·t + q⃗ × t, which costs two cross products. Shoemake and Diebel
// write the frame rotation q⁻¹ v q; this is its inverse.
// Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component
// of the result is outside the float64 range.
//
// ID:      NUM-LINALG-112
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.2 "How quaternions rotate", p. 247 (rotation by conjugation, x′ = q⁻¹xq; the active form is its inverse)
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.4, eqs. (121)–(127)
// Tags:    quaternion, rotate, rotation, vector, 3d, linear-algebra
// Units:   q dimensionless, v in any unit; the result in v's unit
// Example: q of π/2 about ẑ rotates x̂ to ŷ (hand-checked)
func QuatRotateVec3(q Quat, v Vec3) (Vec3, error) {
	if !finiteQ(q) {
		return Vec3{}, errNotFinite("q")
	}
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	u, ok := unitQ(q)
	if !ok {
		return Vec3{}, fmt.Errorf("%w: q", ErrZeroVector)
	}
	r := rotateByUnitQ(u, v)
	if !finite3(r) {
		return Vec3{}, errOverflow("rotated v")
	}
	return r, nil
}

func rotateByUnitQ(u Quat, v Vec3) Vec3 {
	qv := Vec3{u.X, u.Y, u.Z}
	t := cross3(qv, v)
	t = Vec3{2 * t[0], 2 * t[1], 2 * t[2]}
	c := cross3(qv, t)
	return Vec3{v[0] + u.W*t[0] + c[0], v[1] + u.W*t[1] + c[1], v[2] + u.W*t[2] + c[2]}
}
