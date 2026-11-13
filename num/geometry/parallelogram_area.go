package geometry

import "math"

// ParallelogramAreaSAS returns the area of a parallelogram from two adjacent sides and the angle between
// them: a·b·sin θ (h = a sin θ is the height on side b). With equal sides it is a rhombus.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrOutOfRange when
// angleRad is outside (0, π); ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.2 "Quadrilaterals", p. 322 (parallelogram: h = a sin A, area = bh)
// Tags:    geometry, plane, quadrilateral, parallelogram, rhombus
// Units:   sides in one length unit; angleRad in radians; result in that unit squared
// Example: 3, 4, π/6 → 6 (hand-derived)
func ParallelogramAreaSAS(sideA, sideB, angleRad float64) (float64, error) {
	if err := checkPositive(sideA, "sideA"); err != nil {
		return 0, err
	}
	if err := checkPositive(sideB, "sideB"); err != nil {
		return 0, err
	}
	if err := checkOpenAngle(angleRad, math.Pi, "angleRad"); err != nil {
		return 0, err
	}
	return checkResult(product(sideA, sideB, math.Sin(angleRad), 1), "area")
}
