package geometry

import "math"

// coneSlantHalf returns half the slant height ℓ = √(r² + h²) of a right circular cone, computed from the
// halved legs so that it is always finite; callers multiply the 2 back in a product.
func coneSlantHalf(radius, height float64) float64 {
	return math.Hypot(radius/2, height/2)
}

// ConeVolume returns the volume of a circular cone from its base radius and perpendicular height:
// ⅓πr²h. The same formula holds for an oblique cone.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the volume is beyond float64.
//
// ID:      NUM-GEOM-058
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.16 "Cones", p. 362 (right circular cone: volume ⅓πr²h; eq. (4.16.1) for any base)
// Tags:    geometry, solid, cone
// Units:   radius and height in one length unit; result in that unit cubed
// Example: 3, 4 → 12π ≈ 37.699111843078 (hand-derived)
func ConeVolume(radius, height float64) (float64, error) {
	if err := checkRadiusHeight(radius, height); err != nil {
		return 0, err
	}
	return checkResult(product(math.Pi/3, radius, radius, height), "volume")
}

// ConeLateralArea returns the lateral (curved) area of a right circular cone: πr·ℓ, ℓ = √(r² + h²) the
// slant height (math.Hypot, so no square overflows).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-059
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.16 "Cones", p. 362 (ℓ = √(r² + h²), lateral area = πrℓ)
// Tags:    geometry, solid, cone
// Units:   radius and height in one length unit; result in that unit squared
// Example: 3, 4 → 15π ≈ 47.123889803847 (hand-derived: ℓ = 5)
func ConeLateralArea(radius, height float64) (float64, error) {
	if err := checkRadiusHeight(radius, height); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, radius, coneSlantHalf(radius, height), 1), "area")
}

// ConeSurfaceArea returns the total surface area of a right circular cone (the lateral area and the base):
// πr(r + ℓ), ℓ = √(r² + h²).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-060
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.16 "Cones", p. 362 (total area = πr(ℓ + r))
// Tags:    geometry, solid, cone
// Units:   radius and height in one length unit; result in that unit squared
// Example: 3, 4 → 24π ≈ 75.398223686155 (hand-derived)
func ConeSurfaceArea(radius, height float64) (float64, error) {
	if err := checkRadiusHeight(radius, height); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, radius, radius/2+coneSlantHalf(radius, height), 1), "area")
}
