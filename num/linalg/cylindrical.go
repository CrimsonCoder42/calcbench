package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Cylindrical holds cylindrical coordinates in the convention of NIST DLMF eq. 1.5.14: R is the distance
// from the z axis (R ≥ 0), PhiRad the azimuth from the +x axis toward +y in radians, Z the height.
type Cylindrical struct {
	R, PhiRad, Z float64
}

// azimuth returns atan2(y, x) moved into [0, 2π), the range of DLMF §1.5(ii) with 2π folded onto 0.
func azimuth(y, x float64) float64 {
	phi := math.Atan2(y, x)
	if phi < 0 {
		phi += 2 * math.Pi
		if phi >= 2*math.Pi { // a tiny negative angle rounds up to 2π
			phi = 0
		}
	}
	return phi + 0 // turns −0 into +0
}

// CartesianToCylindrical returns the cylindrical coordinates (R, φ, z) of the point v = (x, y, z):
// R = √(x² + y²) (without overflow in the squares), φ = atan2(y, x) in [0, 2π), z unchanged. On the z
// axis (x = y = 0) the azimuth is undefined and 0 is returned.
// Errors: ErrNotFinite for a NaN or ±Inf component; ErrOverflow when R exceeds the float64 range.
//
// ID:      NUM-LINALG-052
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.14 (inverted), https://dlmf.nist.gov/1.5 (read 2026-11-12)
// Tags:    coordinates, cylindrical, cartesian, conversion, 3d, linear-algebra
// Units:   x, y, z in one length unit; R and Z in that unit, PhiRad in radians
// Example: (1, 1, 5) → (√2, π/4, 5) (hand-checked)
func CartesianToCylindrical(v Vec3) (Cylindrical, error) {
	if !finite3(v) {
		return Cylindrical{}, errNotFinite("v")
	}
	r := math.Hypot(v[0], v[1])
	if !reduce.AllFinite(r) {
		return Cylindrical{}, errOverflow("R")
	}
	return Cylindrical{R: r, PhiRad: azimuth(v[1], v[0]), Z: v[2]}, nil
}

// CylindricalToCartesian returns the point (R cos φ, R sin φ, z). Any finite φ is accepted (it need not
// be reduced to [0, 2π)).
// Errors: ErrOutOfRange when R < 0; ErrNotFinite for a NaN or ±Inf coordinate.
//
// ID:      NUM-LINALG-053
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, §1.5(ii) "Coordinate Systems", eq. 1.5.14, https://dlmf.nist.gov/1.5 (read 2026-11-12)
// Tags:    coordinates, cylindrical, cartesian, conversion, 3d, linear-algebra
// Units:   R and Z in one length unit, PhiRad in radians; the result in that length unit
// Example: (2, π/2, 3) → (0, 2, 3) (hand-checked; the computed x is 2·cos(π/2) ≈ 1.2e−16)
func CylindricalToCartesian(c Cylindrical) (Vec3, error) {
	if !reduce.AllFinite(c.R) {
		return Vec3{}, errNotFinite("R")
	}
	if !reduce.AllFinite(c.PhiRad) {
		return Vec3{}, errNotFinite("PhiRad")
	}
	if !reduce.AllFinite(c.Z) {
		return Vec3{}, errNotFinite("Z")
	}
	if c.R < 0 {
		return Vec3{}, fmt.Errorf("%w: R = %v (need R ≥ 0)", ErrOutOfRange, c.R)
	}
	s, co := math.Sincos(c.PhiRad)
	return Vec3{c.R * co, c.R * s, c.Z}, nil
}
