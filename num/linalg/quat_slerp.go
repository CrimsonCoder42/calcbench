package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// sinc returns sin(x)/x, with sinc(0) = 1.
func sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Sin(x) / x
}

// QuatSlerp returns the spherical linear interpolation between the rotations q0 and q1 at t ∈ [0, 1]:
// Slerp = sin((1−t)θ)/sin θ · q0 + sin(tθ)/sin θ · q1, with cos θ = q0·q1, which turns at a constant
// angular rate along the great arc. Both inputs are normalised first, and q1 is negated when q0·q1 < 0 so
// the interpolation takes the shorter of the two arcs (q and −q are the same rotation); at t = 1 the
// result may therefore be −q1. θ is computed as 2·atan2(|q0 − q1|, |q0 + q1|) (accurate for nearly equal
// rotations), and the weights as s·sinc(sθ)/sinc(θ), which equal (1−t, t) exactly when θ = 0, so no
// threshold is needed. At t = 0 the result is exactly the normalised q0.
// Errors: ErrOutOfRange when t is outside [0, 1]; ErrZeroVector when q0 or q1 is zero; ErrNotFinite for a
// NaN or ±Inf input.
//
// ID:      NUM-LINALG-113
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shoemake, "Animating Rotation with Quaternion Curves", Computer Graphics (SIGGRAPH '85), 1985, 19(3):245–254, doi:10.1145/325165.325242, §3.3 "Great arc in-betweening", p. 248 (Slerp)
// Tags:    quaternion, slerp, interpolation, rotation, animation, linear-algebra
// Units:   dimensionless quaternions, t dimensionless
// Example: identity to π/2 about ẑ at t = 0.5 → π/4 about ẑ (hand-checked)
func QuatSlerp(q0, q1 Quat, t float64) (Quat, error) {
	if !finiteQ(q0) {
		return Quat{}, errNotFinite("q0")
	}
	if !finiteQ(q1) {
		return Quat{}, errNotFinite("q1")
	}
	if !reduce.AllFinite(t) {
		return Quat{}, errNotFinite("t")
	}
	if t < 0 || t > 1 {
		return Quat{}, fmt.Errorf("%w: t = %v (need 0 ≤ t ≤ 1)", ErrOutOfRange, t)
	}
	a, ok := unitQ(q0)
	if !ok {
		return Quat{}, fmt.Errorf("%w: q0", ErrZeroVector)
	}
	b, ok := unitQ(q1)
	if !ok {
		return Quat{}, fmt.Errorf("%w: q1", ErrZeroVector)
	}
	if a.W*b.W+a.X*b.X+a.Y*b.Y+a.Z*b.Z < 0 {
		b = Quat{-b.W, -b.X, -b.Y, -b.Z}
	}
	d := norm4(Quat{a.W - b.W, a.X - b.X, a.Y - b.Y, a.Z - b.Z})
	s := norm4(Quat{a.W + b.W, a.X + b.X, a.Y + b.Y, a.Z + b.Z})
	theta := 2 * math.Atan2(d, s)
	st := sinc(theta)
	w0 := (1 - t) * sinc((1-t)*theta) / st
	w1 := t * sinc(t*theta) / st
	return Quat{
		math.FMA(w0, a.W, w1*b.W),
		math.FMA(w0, a.X, w1*b.X),
		math.FMA(w0, a.Y, w1*b.Y),
		math.FMA(w0, a.Z, w1*b.Z),
	}, nil
}
