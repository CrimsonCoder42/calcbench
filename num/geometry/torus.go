package geometry

import (
	"fmt"
	"math"
)

// checkTorus validates a ring torus: both radii finite and > 0, minorRadius ≤ majorRadius (a horn torus,
// r = R, is allowed; a spindle torus, r > R, intersects itself and the formulas no longer hold).
func checkTorus(majorRadius, minorRadius float64) error {
	if err := checkPositive(majorRadius, "majorRadius"); err != nil {
		return err
	}
	if err := checkPositive(minorRadius, "minorRadius"); err != nil {
		return err
	}
	if minorRadius > majorRadius {
		return fmt.Errorf("%w: minorRadius > majorRadius", ErrOutOfRange)
	}
	return nil
}

// TorusVolume returns the volume of a ring torus from R, the distance from the axis to the centre of the
// tube, and r, the tube's radius: 2π²Rr² (Pappus: the disc πr² swept along a circle of length 2πR).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// minorRadius > majorRadius; ErrOverflow when the volume is beyond float64.
//
// ID:      NUM-GEOM-067
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.17 "Surfaces of revolution: the torus", p. 363 (its volume is 2π²Rr²; the Pappus–Guldinus theorem)
// Tags:    geometry, solid, quadric, torus
// Units:   majorRadius and minorRadius in one length unit; result in that unit cubed
// Example: 3, 1 → 6π² ≈ 59.217626406536 (hand-derived)
func TorusVolume(majorRadius, minorRadius float64) (float64, error) {
	if err := checkTorus(majorRadius, minorRadius); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi*math.Pi, majorRadius, minorRadius, minorRadius), "volume")
}

// TorusSurfaceArea returns the surface area of a ring torus from R, the distance from the axis to the
// centre of the tube, and r, the tube's radius: 4π²Rr.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// minorRadius > majorRadius; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-068
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.17 "Surfaces of revolution: the torus", p. 363 (the area of the torus is 4π²Rr)
// Tags:    geometry, solid, quadric, torus
// Units:   majorRadius and minorRadius in one length unit; result in that unit squared
// Example: 3, 1 → 12π² ≈ 118.435252813072 (hand-derived)
func TorusSurfaceArea(majorRadius, minorRadius float64) (float64, error) {
	if err := checkTorus(majorRadius, minorRadius); err != nil {
		return 0, err
	}
	return checkResult(product(4*math.Pi*math.Pi, majorRadius, minorRadius, 1), "area")
}
