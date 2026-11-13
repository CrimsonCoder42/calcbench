package geometry

// CubeVolume returns the volume of a cube (6 squares) from its edge length a: a³.
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume
// is beyond float64.
//
// ID:      NUM-GEOM-073
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, cube
// Units:   edge in a length unit; result in that unit cubed
// Example: 1 → 1 (CRC §4.14.1 table, V/a³)
func CubeVolume(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(edge, edge, edge, 1), "volume")
}

// CubeSurfaceArea returns the total surface area of a cube (6 squares) from its edge length a: 6a².
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area
// is beyond float64.
//
// ID:      NUM-GEOM-074
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, cube
// Units:   edge in a length unit; result in that unit squared
// Example: 1 → 6 (CRC §4.14.1 table, S/a²)
func CubeSurfaceArea(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(6, edge, edge, 1), "area")
}
