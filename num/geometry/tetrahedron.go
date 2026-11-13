package geometry

import "math"

// sqrt3 is √3 to 50 significant digits; the compiler rounds every constant expression that uses it once.
const sqrt3 = 1.7320508075688772935274463415058723669428052538104

// TetrahedronVolume returns the volume of a regular tetrahedron (4 equilateral triangles) from its edge
// length a: (√2/12)·a³.
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume
// is beyond float64.
//
// ID:      NUM-GEOM-071
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, tetrahedron
// Units:   edge in a length unit; result in that unit cubed
// Example: 1 → 0.117851 (CRC §4.14.1 table, V/a³)
func TetrahedronVolume(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(math.Sqrt2/12, edge, edge, edge), "volume")
}

// TetrahedronSurfaceArea returns the total surface area of a regular tetrahedron (4 equilateral triangles)
// from its edge length a: √3·a².
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area
// is beyond float64.
//
// ID:      NUM-GEOM-072
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, tetrahedron
// Units:   edge in a length unit; result in that unit squared
// Example: 1 → 1.73205 (CRC §4.14.1 table, S/a²)
func TetrahedronSurfaceArea(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(sqrt3, edge, edge, 1), "area")
}
