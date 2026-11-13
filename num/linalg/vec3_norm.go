package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// norm3 returns |v| without overflow or underflow in the squares: when the largest component is
// outside [2⁻⁵⁰⁰, 2⁵⁰⁰] the vector is scaled by an exact power of two first (Blue 1978). The result is
// +Inf only when the true length exceeds the float64 range.
func norm3(v Vec3) float64 {
	x, y, z := math.Abs(v[0]), math.Abs(v[1]), math.Abs(v[2])
	m := max(x, y, z)
	if m == 0 {
		return 0
	}
	if m > 0x1p-500 && m < 0x1p500 {
		return math.Sqrt(math.FMA(x, x, math.FMA(y, y, z*z)))
	}
	s, e := pow2Scale(m)
	x, y, z = x*s, y*s, z*s
	return math.Ldexp(math.Sqrt(math.FMA(x, x, math.FMA(y, y, z*z))), e)
}

// unit3 returns v/|v| and false for the zero vector; it scales first so that |v| cannot overflow.
func unit3(v Vec3) (Vec3, bool) {
	m := max(math.Abs(v[0]), math.Abs(v[1]), math.Abs(v[2]))
	if m == 0 {
		return Vec3{}, false
	}
	s, _ := pow2Scale(m)
	w := Vec3{v[0] * s, v[1] * s, v[2] * s}
	n := math.Sqrt(dot3(w, w))
	return Vec3{w[0] / n, w[1] / n, w[2] / n}, true
}

// Vec3Norm returns the Euclidean length |v| = √(x² + y² + z²), without overflow or underflow in the
// squares.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the length exceeds the float64
// range.
//
// ID:      NUM-LINALG-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Source:  Blue, "A Portable Fortran Program to Find the Euclidean Norm of a Vector", ACM Transactions on Mathematical Software, 1978, 4(1):15–23, doi:10.1145/355769.355771
// Tags:    vector, 3d, norm, length, magnitude, linear-algebra
// Units:   v in any unit; the result in that unit
// Example: (2, 3, 6) → 7 (hand-checked)
func Vec3Norm(v Vec3) (float64, error) {
	if !finite3(v) {
		return 0, errNotFinite("v")
	}
	n := norm3(v)
	if !reduce.AllFinite(n) {
		return 0, errOverflow("|v|")
	}
	return n, nil
}

// Vec3Normalize returns the unit vector v/|v|.
// Errors: ErrZeroVector when v = 0; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Tags:    vector, 3d, normalize, unit-vector, direction, linear-algebra
// Units:   v in any unit; the result is dimensionless
// Example: (2, 3, 6) → (2/7, 3/7, 6/7) (hand-checked)
func Vec3Normalize(v Vec3) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	u, ok := unit3(v)
	if !ok {
		return Vec3{}, fmt.Errorf("%w: v", ErrZeroVector)
	}
	return u, nil
}

// Vec3Distance returns |a − b|, the distance between the endpoints of a and b (two points in space).
// The 2D distance between points is num/geometry's PointDistance.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when the distance exceeds the float64
// range.
//
// ID:      NUM-LINALG-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (length of v − w)
// Tags:    vector, 3d, distance, euclidean-distance, linear-algebra
// Units:   a and b in one unit; the result in that unit
// Example: (1, 2, 3) and (4, 6, 15) → 13 (hand-checked: 3² + 4² + 12² = 169)
func Vec3Distance(a, b Vec3) (float64, error) {
	if !finite3(a) {
		return 0, errNotFinite("a")
	}
	if !finite3(b) {
		return 0, errNotFinite("b")
	}
	d := Vec3{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
	if !finite3(d) {
		return 0, errOverflow("a − b")
	}
	n := norm3(d)
	if !reduce.AllFinite(n) {
		return 0, errOverflow("|a − b|")
	}
	return n, nil
}
