package geometry

import "math"

// PyramidVolume returns the volume of a pyramid (a cone over a polygonal base) of any base shape from its
// base area and perpendicular height: ⅓·A·h. A circular base is ConeVolume.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the volume is beyond float64.
//
// ID:      NUM-GEOM-054
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.16 "Cones", eq. (4.16.1), p. 362 (V = ⅓Ah), and §4.14 "Polyhedra", p. 357 (a cone with a polygonal base is a pyramid)
// Tags:    geometry, solid, pyramid
// Units:   baseArea in a length unit squared; height in that unit; result in that unit cubed
// Example: base 36 (a square of side 6), height 4 → 48 (hand-derived)
func PyramidVolume(baseArea, height float64) (float64, error) {
	if err := checkPositive(baseArea, "baseArea"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	return checkResult(product(1.0/3, baseArea, height, 1), "volume")
}

// RegularPyramidSurfaceArea returns the total surface area of a right pyramid whose base is a regular
// polygon, from the number of sides, the side length s and the height h: the base n·s²/(4 tan(π/n)) plus
// n triangular faces of area ½·s·ℓ, where the slant height ℓ = √(h² + r²) and r = s/(2 tan(π/n)) is the
// base's apothem.
// Errors: ErrTooFewSides when sides < 3; ErrNotFinite when side or height is NaN or ±Inf; ErrNonPositive
// when side or height ≤ 0; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-055
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.3 "Regular polygons", p. 324 (base area ¼ka² cot(180°/k) and apothem r = ½a cot(180°/k)), §4.5.1 "Triangles", p. 319 (face area ½h_c c) and p. 320 (Pythagoras, for the slant height)
// Tags:    geometry, solid, pyramid
// Units:   sides is a count; side and height in one length unit; result in that unit squared
// Example: 4, 6, 4 → 36 + 4·(½·6·5) = 96 (hand-derived: slant height √(16 + 9) = 5)
func RegularPyramidSurfaceArea(sides int, side, height float64) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	if err := checkPositive(side, "side"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	n := float64(sides)
	tan := math.Tan(math.Pi / n)
	base := product(n/(4*tan), side, side, 1)
	apothem := product(side, 0.5/tan, 1, 1)
	// Halving both legs keeps math.Hypot finite; the factor 2 goes back in the product.
	halfSlant := math.Hypot(height/2, apothem/2)
	lateral := product(n, side, halfSlant, 1)
	return checkResult(base+lateral, "area")
}

// PyramidFrustumVolume returns the volume of a frustum of a pyramid (or of any cone) from the areas of its
// two parallel bases and the distance between them: (h/3)(A₁ + A₂ + √(A₁A₂)), with √(A₁A₂) taken as
// √A₁·√A₂ so the product cannot overflow. A top area of 0 is a pyramid (PyramidVolume).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the volume is beyond float64.
//
// ID:      NUM-GEOM-056
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.16 "Cones", eq. (4.16.2), p. 362 (V = ⅓h(A + A′ + √(AA′)))
// Tags:    geometry, solid, pyramid, frustum
// Units:   base areas in a length unit squared; height in that unit; result in that unit cubed
// Example: 16, 4, 3 → 28 (hand-derived; PrismoidVolume gives the same for this solid)
func PyramidFrustumVolume(baseArea1, baseArea2, height float64) (float64, error) {
	if err := checkPositive(baseArea1, "baseArea1"); err != nil {
		return 0, err
	}
	if err := checkPositive(baseArea2, "baseArea2"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	sum := baseArea1/3 + baseArea2/3 + math.Sqrt(baseArea1)*math.Sqrt(baseArea2)/3
	return checkResult(product(sum, height, 1, 1), "volume")
}
