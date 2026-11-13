package geometry

// IcosahedronVolume returns the volume of a regular icosahedron (20 equilateral triangles) from its edge
// length a: (5(3 + √5)/12)·a³.
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume
// is beyond float64.
//
// ID:      NUM-GEOM-079
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, icosahedron
// Units:   edge in a length unit; result in that unit cubed
// Example: 1 → 2.181695 (CRC §4.14.1 table, V/a³)
func IcosahedronVolume(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(5*(3+sqrt5)/12, edge, edge, edge), "volume")
}

// IcosahedronSurfaceArea returns the total surface area of a regular icosahedron (20 equilateral triangles)
// from its edge length a: 5√3·a².
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area
// is beyond float64.
//
// ID:      NUM-GEOM-080
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, icosahedron
// Units:   edge in a length unit; result in that unit squared
// Example: 1 → 8.66025 (CRC §4.14.1 table, S/a²)
func IcosahedronSurfaceArea(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(5*sqrt3, edge, edge, 1), "area")
}
