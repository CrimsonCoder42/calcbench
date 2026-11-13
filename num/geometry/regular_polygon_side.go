package geometry

import "math"

// RegularPolygonSideInscribed returns the side of a regular polygon with the given number of sides
// inscribed in a circle of the given radius (its circumradius): 2R·sin(π/n).
// Errors: ErrTooFewSides when sides < 3; ErrNotFinite when circumradius is NaN or ±Inf; ErrNonPositive
// when circumradius ≤ 0; ErrOverflow when the side is beyond float64.
//
// ID:      NUM-GEOM-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.3 "Regular polygons", p. 324 (a = 2R sin(180°/k))
// Tags:    geometry, plane, polygon, regular-polygon
// Units:   sides is a count; circumradius in a length unit; result in that unit
// Example: 4, 1 → √2; 6, 1 → 1 (hand-derived)
func RegularPolygonSideInscribed(sides int, circumradius float64) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	if err := checkPositive(circumradius, "circumradius"); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Sin(math.Pi/float64(sides)), circumradius, 1, 1), "side")
}

// RegularPolygonSideCircumscribed returns the side of a regular polygon with the given number of sides
// circumscribed about a circle of the given radius (its inradius, the apothem): 2r·tan(π/n).
// Errors: ErrTooFewSides when sides < 3; ErrNotFinite when inradius is NaN or ±Inf; ErrNonPositive when
// inradius ≤ 0; ErrOverflow when the side is beyond float64.
//
// ID:      NUM-GEOM-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.3 "Regular polygons", p. 324 (a = 2r tan(180°/k))
// Tags:    geometry, plane, polygon, regular-polygon
// Units:   sides is a count; inradius in a length unit; result in that unit
// Example: 4, 1 → 2 (hand-derived)
func RegularPolygonSideCircumscribed(sides int, inradius float64) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	if err := checkPositive(inradius, "inradius"); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Tan(math.Pi/float64(sides)), inradius, 1, 1), "side")
}
