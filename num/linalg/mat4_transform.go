package linalg

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Homogeneous transforms here act on COLUMN vectors (x, y, z, 1): the translation sits in the last
// column. Strang §10.6 multiplies row vectors, so his matrices are the transposes of these.

// Mat4Translation returns the homogeneous matrix that moves every point by offset:
// [I t; 0 1] with t = offset in the last column.
// Errors: ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-084
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 495 (translation matrix T in homogeneous coordinates; transposed here for column vectors)
// Tags:    matrix, 4x4, translation, homogeneous, transform, computer-graphics, linear-algebra
// Units:   offset in a length unit; the matrix carries that unit in its last column
// Example: offset (1, 2, 3) applied to the point (0, 0, 0) → (1, 2, 3)
func Mat4Translation(offset Vec3) (Mat4, error) {
	if !finite3(offset) {
		return Mat4{}, errNotFinite("offset")
	}
	return Mat4{
		1, 0, 0, offset[0],
		0, 1, 0, offset[1],
		0, 0, 1, offset[2],
		0, 0, 0, 1,
	}, nil
}

// Mat4Scale returns the homogeneous matrix diag(sx, sy, sz, 1) that rescales the x, y and z directions.
// Errors: ErrNotFinite for a NaN or ±Inf factor.
//
// ID:      NUM-LINALG-085
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", pp. 495–496 (scaling matrix S in homogeneous coordinates)
// Tags:    matrix, 4x4, scale, scaling, homogeneous, transform, computer-graphics, linear-algebra
// Units:   factors dimensionless
// Example: factors (2, 3, 4) applied to (1, 1, 1) → (2, 3, 4)
func Mat4Scale(factors Vec3) (Mat4, error) {
	if !finite3(factors) {
		return Mat4{}, errNotFinite("factors")
	}
	return Mat4{
		factors[0], 0, 0, 0,
		0, factors[1], 0, 0,
		0, 0, factors[2], 0,
		0, 0, 0, 1,
	}, nil
}

// Mat4FromMat3 returns the homogeneous matrix [m 0; 0 1]: the linear map m (a rotation, say) in the
// upper-left block, no translation.
// Errors: ErrNotFinite for a NaN or ±Inf element.
//
// ID:      NUM-LINALG-086
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497 (the rotation Q in the upper-left block of the 4 by 4 matrix R)
// Tags:    matrix, 4x4, homogeneous, rotation, transform, computer-graphics, linear-algebra
// Units:   the unit of m
// Example: Mat3RotationZ(π/2) embedded, applied to (1, 0, 0, 1) → (0, 1, 0, 1)
func Mat4FromMat3(m Mat3) (Mat4, error) {
	if !reduce.AllFinite(m[:]...) {
		return Mat4{}, errNotFinite("m")
	}
	return Mat4{
		m[0], m[1], m[2], 0,
		m[3], m[4], m[5], 0,
		m[6], m[7], m[8], 0,
		0, 0, 0, 1,
	}, nil
}

// Mat4TransformPoint returns the point p transformed by m: (x′, y′, z′, w′) = m·(p, 1), then divided by
// w′ (exactly 1 for an affine m; the perspective divide otherwise).
// Errors: ErrOutOfRange when w′ = 0 (the point maps to infinity); ErrNotFinite for a NaN or ±Inf input;
// ErrOverflow when a component of the result is outside the float64 range.
//
// ID:      NUM-LINALG-087
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 495 (points as (x, y, z, 1) in homogeneous coordinates)
// Tags:    matrix, 4x4, homogeneous, transform, point, computer-graphics, linear-algebra
// Units:   p in the length unit of m's translation; the result in that unit
// Example: translation by (1, 2, 3) applied to (1, 1, 1) → (2, 3, 4)
func Mat4TransformPoint(m Mat4, p Vec3) (Vec3, error) {
	if !reduce.AllFinite(m[:]...) {
		return Vec3{}, errNotFinite("m")
	}
	if !finite3(p) {
		return Vec3{}, errNotFinite("p")
	}
	var r Vec4
	for i := range 4 {
		r[i] = dot4(m[4*i], m[4*i+1], m[4*i+2], m[4*i+3], p[0], p[1], p[2], 1)
	}
	if r[3] == 0 {
		return Vec3{}, fmt.Errorf("%w: w = 0", ErrOutOfRange)
	}
	out := Vec3{r[0], r[1], r[2]}
	if r[3] != 1 {
		out = Vec3{r[0] / r[3], r[1] / r[3], r[2] / r[3]}
	}
	if !finite3(out) || !reduce.AllFinite(r[3]) {
		return Vec3{}, errOverflow("transformed point")
	}
	return out, nil
}

// Mat4TransformDirection returns the direction d transformed by m: the upper-left 3×3 block of m applied
// to d, ignoring translation (a direction is (x, y, z, 0) in homogeneous coordinates). For normals under
// a non-uniform scale, transform by the inverse transpose instead.
// Errors: ErrNotFinite for a NaN or ±Inf input; ErrOverflow when a component of the result is outside the
// float64 range.
//
// ID:      NUM-LINALG-088
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §10.6 "Computer Graphics", p. 497 (the linear part in the upper-left block; translation does not act on vectors)
// Tags:    matrix, 4x4, homogeneous, transform, direction, computer-graphics, linear-algebra
// Units:   d in any unit; the result in that unit
// Example: translation by (1, 2, 3) applied to the direction (1, 1, 1) → (1, 1, 1)
func Mat4TransformDirection(m Mat4, d Vec3) (Vec3, error) {
	if !reduce.AllFinite(m[:]...) {
		return Vec3{}, errNotFinite("m")
	}
	if !finite3(d) {
		return Vec3{}, errNotFinite("d")
	}
	r := Vec3{
		dot3(Vec3{m[0], m[1], m[2]}, d),
		dot3(Vec3{m[4], m[5], m[6]}, d),
		dot3(Vec3{m[8], m[9], m[10]}, d),
	}
	if !finite3(r) {
		return Vec3{}, errOverflow("transformed direction")
	}
	return r, nil
}
