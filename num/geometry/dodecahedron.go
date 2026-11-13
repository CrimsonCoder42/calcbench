package geometry

const (
	// sqrt5 is √5 to 50 significant digits.
	sqrt5 = 2.2360679774997896964091736687312762354406183596115
	// sqrt25Plus10Sqrt5 is √(25 + 10√5) to 50 significant digits.
	sqrt25Plus10Sqrt5 = 6.8819096023558676910360479095544383976294966800408
)

// DodecahedronVolume returns the volume of a regular dodecahedron (12 regular pentagons) from its edge
// length a: ((15 + 7√5)/4)·a³.
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the volume
// is beyond float64.
//
// ID:      NUM-GEOM-077
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of V/a³ for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, dodecahedron
// Units:   edge in a length unit; result in that unit cubed
// Example: 1 → 7.663119 (CRC §4.14.1 table, V/a³)
func DodecahedronVolume(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product((15+7*sqrt5)/4, edge, edge, edge), "volume")
}

// DodecahedronSurfaceArea returns the total surface area of a regular dodecahedron (12 regular pentagons)
// from its edge length a: 3√(25 + 10√5)·a².
// Errors: ErrNotFinite when edge is NaN or ±Inf; ErrNonPositive when edge ≤ 0; ErrOverflow when the area
// is beyond float64.
//
// ID:      NUM-GEOM-078
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.14.1 "Convex regular polyhedra", p. 360 (table of S/a² for the five Platonic solids)
// Tags:    geometry, solid, polyhedron, platonic-solid, dodecahedron
// Units:   edge in a length unit; result in that unit squared
// Example: 1 → 20.64573 (CRC §4.14.1 table, S/a²)
func DodecahedronSurfaceArea(edge float64) (float64, error) {
	if err := checkPositive(edge, "edge"); err != nil {
		return 0, err
	}
	return checkResult(product(3*sqrt25Plus10Sqrt5, edge, edge, 1), "area")
}
