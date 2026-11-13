package geometry

import "math"

// checkFrustum validates the two radii (either order) and the height of a conical frustum.
func checkFrustum(radius1, radius2, height float64) error {
	if err := checkPositive(radius1, "radius1"); err != nil {
		return err
	}
	if err := checkPositive(radius2, "radius2"); err != nil {
		return err
	}
	return checkPositive(height, "height")
}

// ConeFrustumVolume returns the volume of a frustum of a right circular cone from its two radii (either
// order) and height: ⅓πh(r₁² + r₂² + r₁r₂). The bracket is formed relative to the larger radius m, as
// m²·(1 + q + q²) with q = smaller/m ∈ (0, 1], so it cannot overflow before the volume does. A radius of 0
// is a cone (ConeVolume).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the volume is beyond float64.
//
// ID:      NUM-GEOM-061
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.16 "Cones", p. 363 (frustum of a right circular cone: volume = ⅓πh(r₁² + r₂² + r₁r₂))
// Tags:    geometry, solid, cone, frustum
// Units:   radii and height in one length unit; result in that unit cubed
// Example: 4, 1, 4 → 28π ≈ 87.964594300514 (hand-derived)
func ConeFrustumVolume(radius1, radius2, height float64) (float64, error) {
	if err := checkFrustum(radius1, radius2, height); err != nil {
		return 0, err
	}
	m, s := math.Max(radius1, radius2), math.Min(radius1, radius2)
	q := s / m
	return checkResult(product(math.Pi/3*(1+q+q*q), m, m, height), "volume")
}

// ConeFrustumSurfaceArea returns the total surface area of a frustum of a right circular cone (the lateral
// area and both ends): π(r₁ + r₂)ℓ + πr₁² + πr₂², with slant height ℓ = √((r₁ − r₂)² + h²). Each of the
// three terms is a non-negative product formed without intermediate overflow, so the sum overflows only
// when the area does.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-062
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.16 "Cones", p. 363 (frustum: ℓ = √((r₁ − r₂)² + h²), total area = π(r₁² + r₂² + (r₁ + r₂)ℓ))
// Tags:    geometry, solid, cone, frustum
// Units:   radii and height in one length unit; result in that unit squared
// Example: 4, 1, 4 → 42π ≈ 131.946891450771 (hand-derived: ℓ = 5)
func ConeFrustumSurfaceArea(radius1, radius2, height float64) (float64, error) {
	if err := checkFrustum(radius1, radius2, height); err != nil {
		return 0, err
	}
	halfSlant := math.Hypot(radius1/2-radius2/2, height/2)
	lateral := product(4*math.Pi, radius1/2+radius2/2, halfSlant, 1)
	ends := product(math.Pi, radius1, radius1, 1) + product(math.Pi, radius2, radius2, 1)
	return checkResult(lateral+ends, "area")
}
