package geometry

import (
	"fmt"
	"math"
)

// checkCap validates a sphere radius r and a cap height h: both finite and > 0, h ≤ 2r.
func checkCap(sphereRadius, capHeight float64) error {
	if err := checkPositive(sphereRadius, "sphereRadius"); err != nil {
		return err
	}
	if err := checkPositive(capHeight, "capHeight"); err != nil {
		return err
	}
	if capHeight/2 > sphereRadius {
		return fmt.Errorf("%w: capHeight > 2·sphereRadius", ErrOutOfRange)
	}
	return nil
}

// SphericalCapVolume returns the volume of a spherical cap of height h cut from a sphere of radius r:
// (πh²/3)(3r − h), computed as πh²(r − h/3) so 3r cannot overflow. h = 2r is the whole sphere.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// capHeight > 2·sphereRadius; ErrOverflow when the volume is beyond float64.
//
// ID:      NUM-GEOM-065
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.18.1.1 "Spherical cap", p. 367 (the volume of the cap is ⅓πh²(3r − h))
// Tags:    geometry, solid, quadric, sphere, cap
// Units:   sphereRadius and capHeight in one length unit; result in that unit cubed
// Example: 5, 2 → 52π/3 ≈ 54.454272662223 (hand-derived)
func SphericalCapVolume(sphereRadius, capHeight float64) (float64, error) {
	if err := checkCap(sphereRadius, capHeight); err != nil {
		return 0, err
	}
	return checkResult(product(math.Pi, capHeight, capHeight, sphereRadius-capHeight/3), "volume")
}

// SphericalCapArea returns the area of the curved surface of a spherical cap of height h cut from a
// sphere of radius r: 2πrh. The base disc is not included.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// capHeight > 2·sphereRadius; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-066
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.18.1.1 "Spherical cap", p. 367 (the area of the curved region is 2πrh)
// Tags:    geometry, solid, quadric, sphere, cap
// Units:   sphereRadius and capHeight in one length unit; result in that unit squared
// Example: 5, 2 → 20π ≈ 62.831853071796 (hand-derived)
func SphericalCapArea(sphereRadius, capHeight float64) (float64, error) {
	if err := checkCap(sphereRadius, capHeight); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, sphereRadius, capHeight, 1), "area")
}
