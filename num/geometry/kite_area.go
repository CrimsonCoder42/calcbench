package geometry

// KiteArea returns the area of a quadrilateral whose diagonals are perpendicular (a kite, a rhombus, a
// square) from the diagonal lengths: ½·d₁·d₂, the general ½pq·sin θ with θ = π/2.
// Errors: ErrNotFinite when a diagonal is NaN or ±Inf; ErrNonPositive when a diagonal ≤ 0; ErrOverflow
// when the area is beyond float64.
//
// ID:      NUM-GEOM-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.2 "Quadrilaterals", eq. (4.5.9), p. 322 (area = ½pq sin θ) and p. 323 (rhombus: area = ½pq)
// Tags:    geometry, plane, quadrilateral, kite, rhombus, square
// Units:   diagonals in one length unit; result in that unit squared
// Example: 6, 4 → 12 (hand-derived)
func KiteArea(diagonal1, diagonal2 float64) (float64, error) {
	if err := checkPositive(diagonal1, "diagonal1"); err != nil {
		return 0, err
	}
	if err := checkPositive(diagonal2, "diagonal2"); err != nil {
		return 0, err
	}
	return checkResult(product(0.5, diagonal1, diagonal2, 1), "area")
}
