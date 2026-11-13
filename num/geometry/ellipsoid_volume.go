package geometry

import "math"

// checkSemiAxes3 validates the three semi-axes of an ellipsoid: each finite and > 0.
func checkSemiAxes3(semiAxisA, semiAxisB, semiAxisC float64) error {
	if err := checkPositive(semiAxisA, "semiAxisA"); err != nil {
		return err
	}
	if err := checkPositive(semiAxisB, "semiAxisB"); err != nil {
		return err
	}
	return checkPositive(semiAxisC, "semiAxisC")
}

// EllipsoidVolume returns the volume of an ellipsoid from its three semi-axes (any order): (4/3)πabc.
// Errors: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow
// when the volume is beyond float64.
//
// ID:      NUM-GEOM-069
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.18 "Quadrics", p. 364 (the volume of the ellipsoid with semi-axes a, b, c is (4/3)πabc)
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), §19.33 (V = 4πabc/3, beside eq. 19.33.1), https://dlmf.nist.gov/19.33 (read 2026-11-12)
// Tags:    geometry, solid, quadric, ellipsoid
// Units:   semi-axes in one length unit; result in that unit cubed
// Example: 3, 2, 1 → 8π ≈ 25.132741228718 (hand-derived)
func EllipsoidVolume(semiAxisA, semiAxisB, semiAxisC float64) (float64, error) {
	if err := checkSemiAxes3(semiAxisA, semiAxisB, semiAxisC); err != nil {
		return 0, err
	}
	return checkResult(product(4*math.Pi/3, semiAxisA, semiAxisB, semiAxisC), "volume")
}
