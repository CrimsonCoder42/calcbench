package geometry

import "math"

// CircleCircumference returns the circumference of a circle of the given radius, 2πr. From a diameter,
// halve it first.
// Errors: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the
// circumference is beyond float64.
//
// ID:      NUM-GEOM-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 ("The length or circumference of a circle of radius r is 2πr")
// Tags:    geometry, plane, circle
// Units:   radius in a length unit; result in that unit
// Example: 1 → 2π ≈ 6.283185307180 (hand-derived)
func CircleCircumference(radius float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, radius, 1, 1), "circumference")
}

// CircleArea returns the area of a circle of the given radius, πr².
// Errors: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the
// area is beyond float64.
//
// ID:      NUM-GEOM-035
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 ("the area is πr²")
// Tags:    geometry, plane, circle
// Units:   radius in a length unit; result in that unit squared
// Example: 2 → 4π ≈ 12.566370614359 (hand-derived)
func CircleArea(radius float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	return checkResult(product(math.Pi, radius, radius, 1), "area")
}
