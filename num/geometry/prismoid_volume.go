package geometry

// PrismoidVolume returns the volume of a solid between two parallel planes by the prismoidal formula:
// V = (h/6)(A₁ + 4·Aₘ + A₂), where A₁ and A₂ are the end areas and Aₘ the area of the section midway
// between them. It is exact when the area of every parallel cross-section is a quadratic function of its
// distance from one end (prismoids, pyramids and their frustums, cones, spheres, ellipsoids, wedges); for
// other solids it is Simpson's rule, an approximation. It is computed as (A₁/6 + (2/3)Aₘ + A₂/6)·h so the
// sum cannot overflow. An end area may be 0 (a pyramid, a cone, a wedge).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when height ≤ 0; ErrNegative when an
// area < 0; ErrOverflow when the volume is beyond float64.
//
// ID:      NUM-GEOM-057
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Encyclopædia Britannica, 11th ed., 1911, vol. 18, "Mensuration", §58 "The Prismoidal Formula", p. 160 (volume ⅙H(S₀ + 4S₁ + S₂) when the cross-section area is a quadratic function of the distance), https://en.wikisource.org/wiki/Page:EB1911_-_Volume_18.djvu/160 (read 2026-11-12)
// Tags:    geometry, solid, prismoid, prismatoid, trapezoidal-volume, earthwork
// Units:   areas in a length unit squared; height in that unit; result in that unit cubed
// Example: square frustum with end sides 4 and 2, height 3: 16, 9, 4, 3 → 28 (hand-derived; equals PyramidFrustumVolume)
func PrismoidVolume(area1, areaMid, area2, height float64) (float64, error) {
	if err := checkNonNegative(area1, "area1"); err != nil {
		return 0, err
	}
	if err := checkNonNegative(areaMid, "areaMid"); err != nil {
		return 0, err
	}
	if err := checkNonNegative(area2, "area2"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	sum := area1/6 + areaMid*(2.0/3) + area2/6
	return checkResult(product(sum, height, 1, 1), "volume")
}
