package geometry

// TrapezoidArea returns the area of a trapezoid from its two parallel sides and the distance between
// them: ½(a + b)·h, computed as (a/2 + b/2)·h so the sum cannot overflow. A base of 0 is a triangle
// (TriangleAreaBaseHeight).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.2 "Quadrilaterals", p. 323 (trapezoid: area = ½(AB + CD)h)
// Tags:    geometry, plane, quadrilateral, trapezoid
// Units:   baseA, baseB and height in one length unit; result in that unit squared
// Example: 3, 5, 2 → 8 (hand-derived)
func TrapezoidArea(baseA, baseB, height float64) (float64, error) {
	if err := checkPositive(baseA, "baseA"); err != nil {
		return 0, err
	}
	if err := checkPositive(baseB, "baseB"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	return checkResult(product(baseA/2+baseB/2, height, 1, 1), "area")
}
