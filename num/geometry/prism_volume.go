package geometry

// PrismVolume returns the volume of a prism (a cylinder over a polygonal base, right or oblique) from its
// base area and its perpendicular height: A·h. By Cavalieri's principle the base may have any shape; pass
// an area from RectangleArea, RegularPolygonAreaFromSide, PolygonArea and the like.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the volume is beyond float64.
//
// ID:      NUM-GEOM-048
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.15 "Cylinders", eq. (4.15.1), p. 361 (V = Ah, h the perpendicular distance between the planes), and §4.14 "Polyhedra", p. 357 (a cylinder with a polygonal base is a prism)
// Tags:    geometry, solid, prism
// Units:   baseArea in a length unit squared; height in that unit; result in that unit cubed
// Example: base 12, height 5 → 60 (hand-derived)
func PrismVolume(baseArea, height float64) (float64, error) {
	if err := checkPositive(baseArea, "baseArea"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	return checkResult(baseArea*height, "volume")
}
