package geometry

import "math"

// OctahedronVolume returns the volume of a regular octahedron (8 equilateral triangles) from its edge
// length a: (√2/3)·a³.
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume
// is beyond float64.
//
// ID:      NUM-GEOM-075
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, octahedron
// Units:   edge in a length unit; result in that unit cubed
// Example: 1 → 0.471405 (CRC §4.14.1 table, V/a³)
func OctahedronVolume(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(math.Sqrt2/3, edge, edge, edge), "volume")
}

// OctahedronSurfaceArea returns the total surface area of a regular octahedron (8 equilateral triangles)
// from its edge length a: 2√3·a².
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area
// is beyond float64.
//
// ID:      NUM-GEOM-076
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, octahedron
// Units:   edge in a length unit; result in that unit squared
// Example: 1 → 3.46410 (CRC §4.14.1 table, S/a²)
func OctahedronSurfaceArea(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(2*sqrt3, edge, edge, 1), "area")
}
