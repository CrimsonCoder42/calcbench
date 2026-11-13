package geometry

import "math"

// RegularPolygonAreaFromSide returns the area of a regular polygon with the given number of sides and side
// length: n·s²/(4·tan(π/n)) (the source's ¼ka²·cot(180°/k)).
// Errors: ErrTooFewSides when sides < 3; ErrNotFinite when side is NaN or ±Inf; ErrNonPositive when
// side ≤ 0; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.3 "Regular polygons", p. 324 (area = ¼ka² cot(180°/k), and the table of areas for k = 3…12)
// Tags:    geometry, plane, polygon, regular-polygon
// Units:   sides is a count; side in a length unit; result in that unit squared
// Example: 6, 1 → 2.59808 (CRC, §4.5.3 table, regular hexagon: 2.59808a²)
func RegularPolygonAreaFromSide(sides int, side float64) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	if err := checkPositive(side, "side"); err != nil {
		return 0, err
	}
	n := float64(sides)
	return checkResult(product(n/(4*math.Tan(math.Pi/n)), side, side, 1), "area")
}

// RegularPolygonAreaFromApothem returns the area of a regular polygon with the given number of sides and
// apothem (the inradius, the distance from the centre to a side): n·r²·tan(π/n).
// Errors: ErrTooFewSides when sides < 3; ErrNotFinite when apothem is NaN or ±Inf; ErrNonPositive when
// apothem ≤ 0; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.3 "Regular polygons", p. 324 (area = kr² tan(180°/k), r the apothem)
// Tags:    geometry, plane, polygon, regular-polygon
// Units:   sides is a count; apothem in a length unit; result in that unit squared
// Example: 6, √3/2 → 3√3/2 ≈ 2.598076211353 (hand-derived; equals the hexagon of side 1)
func RegularPolygonAreaFromApothem(sides int, apothem float64) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	if err := checkPositive(apothem, "apothem"); err != nil {
		return 0, err
	}
	n := float64(sides)
	return checkResult(product(n*math.Tan(math.Pi/n), apothem, apothem, 1), "area")
}

// RegularPolygonAreaFromCircumradius returns the area of a regular polygon with the given number of sides
// and circumradius (the distance from the centre to a vertex): ½·n·R²·sin(2π/n).
// Errors: ErrTooFewSides when sides < 3; ErrNotFinite when circumradius is NaN or ±Inf; ErrNonPositive
// when circumradius ≤ 0; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.3 "Regular polygons", p. 324 (area = ½kR² sin(360°/k))
// Tags:    geometry, plane, polygon, regular-polygon
// Units:   sides is a count; circumradius in a length unit; result in that unit squared
// Example: 6, 1 → 3√3/2 ≈ 2.598076211353 (hand-derived)
func RegularPolygonAreaFromCircumradius(sides int, circumradius float64) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	if err := checkPositive(circumradius, "circumradius"); err != nil {
		return 0, err
	}
	n := float64(sides)
	return checkResult(product(0.5*n*math.Sin(2*math.Pi/n), circumradius, circumradius, 1), "area")
}
