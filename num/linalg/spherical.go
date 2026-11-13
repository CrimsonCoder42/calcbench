package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Spherical holds spherical coordinates in the physics convention of NIST DLMF eq. 1.5.16: R is the
// distance from the origin (R ≥ 0), ThetaRad the polar angle measured from the +z axis, in [0, π], and
// PhiRad the azimuth from the +x axis toward +y, in [0, 2π).
// Mathematicians often swap the names θ and φ (DLMF §1.5(ii) notes this); geographers use latitude
// π/2 − θ instead of the polar angle.
type Spherical struct {
	R, ThetaRad, PhiRad float64
}

// CartesianToSpherical returns the spherical coordinates (R, θ, φ) of v = (x, y, z): R = |v| (without
// overflow in the squares), θ = atan2(√(x² + y²), z) in [0, π] (accurate at the poles, where
// arccos(z/R) loses half its digits), φ = atan2(y, x) in [0, 2π). At the origin both angles are
// undefined and 0 is returned; on the z axis φ is undefined and 0 is returned.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when R exceeds the float64 range.
//
// ID:      NUM-LINALG-054
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.16 (inverted), https://dlmf.nist.gov/1.5 (read 2026-11-12)
// Tags:    coordinates, spherical, polar-angle, azimuth, cartesian, conversion, 3d, linear-algebra
// Units:   x, y, z in one length unit; R in that unit, ThetaRad and PhiRad in radians
// Example: (0, 1, 1) → (√2, π/4, π/2) (hand-checked)
func CartesianToSpherical(v Vec3) (Spherical, error) {
	if !finite3(v) {
		return Spherical{}, errNotFinite("v")
	}
	r := norm3(v)
	if !reduce.AllFinite(r) {
		return Spherical{}, errOverflow("R")
	}
	if r == 0 {
		return Spherical{}, nil // the origin: both angles undefined, 0 (also for a −0 z, which atan2 would make π)
	}
	theta := math.Atan2(math.Hypot(v[0], v[1]), v[2])
	return Spherical{R: r, ThetaRad: theta, PhiRad: azimuth(v[1], v[0])}, nil
}

// SphericalToCartesian returns the point (R sin θ cos φ, R sin θ sin φ, R cos θ). Any finite angles are
// accepted.
// Errors: ErrOutOfRange when R < 0; ErrNotFinite for a NaN or ±Inf coordinate.
//
// ID:      NUM-LINALG-055
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.16, https://dlmf.nist.gov/1.5 (read 2026-11-12)
// Tags:    coordinates, spherical, polar-angle, azimuth, cartesian, conversion, 3d, linear-algebra
// Units:   R in a length unit, ThetaRad and PhiRad in radians; the result in that length unit
// Example: (2, π/2, 0) → (2, 0, 0) (hand-checked; the computed z is 2·cos(π/2) ≈ 1.2e−16)
func SphericalToCartesian(s Spherical) (Vec3, error) {
	if !reduce.AllFinite(s.R) {
		return Vec3{}, errNotFinite("R")
	}
	if !reduce.AllFinite(s.ThetaRad) {
		return Vec3{}, errNotFinite("ThetaRad")
	}
	if !reduce.AllFinite(s.PhiRad) {
		return Vec3{}, errNotFinite("PhiRad")
	}
	if s.R < 0 {
		return Vec3{}, fmt.Errorf("%w: R = %v (need R ≥ 0)", ErrOutOfRange, s.R)
	}
	st, ct := math.Sincos(s.ThetaRad)
	sp, cp := math.Sincos(s.PhiRad)
	rs := s.R * st
	return Vec3{rs * cp, rs * sp, s.R * ct}, nil
}
