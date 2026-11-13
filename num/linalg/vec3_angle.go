package linalg

import (
	"fmt"
	"math"
)

// Vec3AngleBetween returns the unoriented angle between a and b, in [0, π] radians, by Kahan's formula
// θ = 2·atan2(|â − b̂|, |â + b̂|) with â = a/|a|, b̂ = b/|b|. Unlike arccos of the normalised dot
// product it keeps full accuracy for nearly parallel or antiparallel vectors, and unlike arcsin of the
// cross product it keeps it near π/2.
// Errors: ErrZeroVector when a or b is zero; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12)
// Tags:    vector, 3d, angle, angle-between, linear-algebra
// Units:   a and b in any units; the result in radians
// Example: (1, 0, 0) and (0, 1, 1) → π/2; (1, 0, 0) and (1, 1, 0) → π/4 (hand-checked)
func Vec3AngleBetween(a, b Vec3) (float64, error) {
	if !finite3(a) {
		return 0, errNotFinite("a")
	}
	if !finite3(b) {
		return 0, errNotFinite("b")
	}
	ua, ok := unit3(a)
	if !ok {
		return 0, fmt.Errorf("%w: a", ErrZeroVector)
	}
	ub, ok := unit3(b)
	if !ok {
		return 0, fmt.Errorf("%w: b", ErrZeroVector)
	}
	return kahanAngle3(ua, ub), nil
}

// kahanAngle3 returns 2·atan2(|u − w|, |u + w|) for unit vectors u and w.
func kahanAngle3(u, w Vec3) float64 {
	d := Vec3{u[0] - w[0], u[1] - w[1], u[2] - w[2]}
	s := Vec3{u[0] + w[0], u[1] + w[1], u[2] + w[2]}
	return 2 * math.Atan2(math.Sqrt(dot3(d, d)), math.Sqrt(dot3(s, s)))
}

// Vec3DirectionAngles returns the angles (α, β, γ) between v and the +x, +y and +z axes, each in
// [0, π] radians; their cosines are the direction cosines v/|v|. Each angle is computed as
// atan2(|v × eᵢ|, v·eᵢ), e.g. α = atan2(√(y² + z²), x), which is accurate at every angle (the same
// value Kahan's formula gives for v and the axis).
// Errors: ErrZeroVector when v = 0; ErrNotFinite for a NaN or ±Inf component.
//
// ID:      NUM-LINALG-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products" (cosine formula with w a coordinate axis)
// Source:  Kahan, "How Futile are Mindless Assessments of Roundoff in Floating-Point Computation?", University of California, Berkeley, 11 January 2006, §12 "Mangled Angles", pp. 46–47, https://people.eecs.berkeley.edu/~wkahan/Mindless.pdf (read 2026-11-12)
// Tags:    vector, 3d, direction-angles, direction-cosines, angle, linear-algebra
// Units:   v in any unit; the result in radians
// Example: (1, 1, 0) → (π/4, π/4, π/2) (hand-checked)
func Vec3DirectionAngles(v Vec3) (Vec3, error) {
	if !finite3(v) {
		return Vec3{}, errNotFinite("v")
	}
	u, ok := unit3(v)
	if !ok {
		return Vec3{}, fmt.Errorf("%w: v", ErrZeroVector)
	}
	return Vec3{
		math.Atan2(math.Hypot(u[1], u[2]), u[0]),
		math.Atan2(math.Hypot(u[0], u[2]), u[1]),
		math.Atan2(math.Hypot(u[0], u[1]), u[2]),
	}, nil
}
