package geometry

// RectangleArea returns the area of a rectangle, width·height. The same product is the area of a square
// (equal sides) and of a parallelogram or rhombus from a base and its perpendicular height.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.2 "Quadrilaterals", p. 322 (parallelogram: area = bh; rectangle: area = ab, eq. (4.5.10))
// Tags:    geometry, plane, quadrilateral, rectangle, square, parallelogram, rhombus, education
// Units:   width and height in one length unit; result in that unit squared
// Example: 3, 4 → 12 (hand-derived)
func RectangleArea(width, height float64) (float64, error) {
	if err := checkPositive(width, "width"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	return checkResult(width*height, "area")
}

// RectanglePerimeter returns 2(a + b), the perimeter of a quadrilateral with two pairs of equal sides a and
// b: a rectangle, square, parallelogram or rhombus (opposite sides equal) or a kite (adjacent sides equal).
// Errors: ErrNotFinite when a side is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrOverflow when the
// perimeter is beyond float64.
//
// ID:      NUM-GEOM-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.2 "Quadrilaterals", p. 322 (Figure 4.13: half-perimeter s = ½(a + b + c + d); in a parallelogram opposite sides have the same length)
// Tags:    geometry, plane, quadrilateral, rectangle, square, parallelogram, rhombus, kite
// Units:   sides in one length unit; result in that unit
// Example: 3, 4 → 14 (hand-derived)
func RectanglePerimeter(sideA, sideB float64) (float64, error) {
	if err := checkPositive(sideA, "sideA"); err != nil {
		return 0, err
	}
	if err := checkPositive(sideB, "sideB"); err != nil {
		return 0, err
	}
	return checkResult(2*(sideA+sideB), "perimeter")
}
