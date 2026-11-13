package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// QuatFromMat3 returns the unit quaternion (with W ≥ 0) of the rotation matrix m (active, column
// vectors). Of the four candidate formulas it uses the one whose square root has the largest argument
// (1 + trace, or 1 + 2mᵢᵢ − trace), so no division is by a small number (Shepperd's method). m is assumed
// to be a rotation (orthogonal, det +1); it is not checked, and the result is normalised, so a matrix
// slightly off orthogonal from rounding gives the nearest-looking unit quaternion.
// Errors: ErrNotFinite for a NaN or ±Inf element; ErrOverflow when sums of elements overflow (elements
// near the float64 maximum, which no rotation matrix has).
//
// ID:      NUM-LINALG-110
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shepperd, "Quaternion from Rotation Matrix", Journal of Guidance and Control, 1978, 1(3):223–224, doi:10.2514/3.55767b
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.5, eqs. (131)–(145) (for Diebel's frame matrix; the active matrix here is its transpose)
// Tags:    quaternion, rotation-matrix, conversion, shepperd, linear-algebra
// Units:   dimensionless
// Example: Mat3RotationZ(π/2) → (√2/2, 0, 0, √2/2) (hand-checked)
func QuatFromMat3(m Mat3) (Quat, error) {
	if !reduce.AllFinite(m[:]...) {
		return Quat{}, errNotFinite("m")
	}
	m00, m01, m02 := m[0], m[1], m[2]
	m10, m11, m12 := m[3], m[4], m[5]
	m20, m21, m22 := m[6], m[7], m[8]
	tr := m00 + m11 + m22
	var q Quat
	switch {
	case tr >= m00 && tr >= m11 && tr >= m22:
		w := math.Sqrt(1+tr) / 2
		f := 1 / (4 * w)
		q = Quat{w, (m21 - m12) * f, (m02 - m20) * f, (m10 - m01) * f}
	case m00 >= m11 && m00 >= m22:
		x := math.Sqrt(1+m00-m11-m22) / 2
		f := 1 / (4 * x)
		q = Quat{(m21 - m12) * f, x, (m01 + m10) * f, (m02 + m20) * f}
	case m11 >= m22:
		y := math.Sqrt(1-m00+m11-m22) / 2
		f := 1 / (4 * y)
		q = Quat{(m02 - m20) * f, (m01 + m10) * f, y, (m12 + m21) * f}
	default:
		z := math.Sqrt(1-m00-m11+m22) / 2
		f := 1 / (4 * z)
		q = Quat{(m10 - m01) * f, (m02 + m20) * f, (m12 + m21) * f, z}
	}
	if q.W < 0 {
		q = Quat{-q.W, -q.X, -q.Y, -q.Z}
	}
	if !finiteQ(q) {
		return Quat{}, errOverflow("quaternion of m")
	}
	// The largest square-root argument is at least 1 (the four arguments sum to 4), so q ≠ 0.
	u, _ := unitQ(q)
	return u, nil
}

// QuatToMat3 returns the active rotation matrix of q (normalised first, so any nonzero q works):
// with s = 2/|q|²,
// [1 − s(Y²+Z²)  s(XY − WZ)  s(XZ + WY); s(XY + WZ)  1 − s(X²+Z²)  s(YZ − WX); s(XZ − WY)  s(YZ + WX)  1 − s(X²+Y²)].
// It is the transpose of Diebel's frame-rotation matrix for the same q.
// Errors: ErrZeroVector when q = 0; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-111
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Diebel, "Representing Attitude: Euler Angles, Unit Quaternions, and Rotation Vectors", Stanford University, 20 October 2006, §6.4, eq. (125) (transposed: Diebel's matrix is the frame rotation)
// Source:  Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.2 "How quaternions rotate", p. 248 (the same matrix, for row vectors)
// Tags:    quaternion, rotation-matrix, conversion, linear-algebra
// Units:   dimensionless
// Example: (√2/2, 0, 0, √2/2) → [0 −1 0; 1 0 0; 0 0 1] (to rounding; hand-checked)
func QuatToMat3(q Quat) (Mat3, error) {
	if !finiteQ(q) {
		return Mat3{}, errNotFinite("q")
	}
	u, ok := unitQ(q)
	if !ok {
		return Mat3{}, fmt.Errorf("%w: q", ErrZeroVector)
	}
	return rotationMatrix(u), nil
}

// rotationMatrix returns the active rotation matrix of the unit quaternion u.
func rotationMatrix(u Quat) Mat3 {
	w, x, y, z := u.W, u.X, u.Y, u.Z
	return Mat3{
		1 - 2*(y*y+z*z), 2 * (x*y - w*z), 2 * (x*z + w*y),
		2 * (x*y + w*z), 1 - 2*(x*x+z*z), 2 * (y*z - w*x),
		2 * (x*z - w*y), 2 * (y*z + w*x), 1 - 2*(x*x+y*y),
	}
}
