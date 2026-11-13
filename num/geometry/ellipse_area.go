package geometry

import "math"

// EllipseArea returns the area of an ellipse from its two semi-axes (either order): π·a·b.
// Errors: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow
// when the area is beyond float64.
//
// ID:      NUM-GEOM-043
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.3 "Additional properties of ellipses", eq. (4.6.11), p. 330 (area of C = πab)
// Tags:    geometry, plane, ellipse, conic
// Units:   semi-axes in one length unit; result in that unit squared
// Example: 5, 3 → 15π ≈ 47.123889803847 (hand-derived)
func EllipseArea(semiAxisA, semiAxisB float64) (float64, error) {
	if err := checkPositive(semiAxisA, "semiAxisA"); err != nil {
		return 0, err
	}
	if err := checkPositive(semiAxisB, "semiAxisB"); err != nil {
		return 0, err
	}
	return checkResult(product(math.Pi, semiAxisA, semiAxisB, 1), "area")
}
