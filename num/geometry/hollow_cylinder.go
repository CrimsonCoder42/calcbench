package geometry

import "math"

// HollowCylinderVolume returns the volume of a hollow circular cylinder (a tube or pipe wall) from its
// outer and inner radii and its height: the prism of an annular base, π(R² − r²)h, computed as
// 2π(R − r)(R/2 + r/2)h with no cancellation for a thin wall and no R² to overflow.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// innerRadius ≥ outerRadius; ErrOverflow when the volume is beyond float64.
//
// ID:      NUM-GEOM-052
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.15 "Cylinders", eq. (4.15.1), p. 361 (V = Ah), with the area of the annular base from §4.6.6 "Circles", p. 334 (area πr²)
// Tags:    geometry, solid, cylinder, tube, pipe
// Units:   radii and height in one length unit; result in that unit cubed
// Example: 3, 2, 5 → 25π ≈ 78.539816339745 (hand-derived)
func HollowCylinderVolume(outerRadius, innerRadius, height float64) (float64, error) {
	if err := checkRadii(outerRadius, innerRadius); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, outerRadius-innerRadius, outerRadius/2+innerRadius/2, height), "volume")
}

// HollowCylinderSurfaceArea returns the total surface area of a hollow circular cylinder: the outer and
// inner walls and the two annular ends, 2π(R + r)h + 2π(R − r)(R + r), computed as the single product
// 4π(R/2 + r/2)(h + (R − r)).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// innerRadius ≥ outerRadius; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-053
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.15 "Cylinders", p. 361 (lateral area 2πrh, for each wall), with the annular ends from §4.6.6 "Circles", p. 334 (area πr²)
// Tags:    geometry, solid, cylinder, tube, pipe
// Units:   radii and height in one length unit; result in that unit squared
// Example: 3, 2, 5 → 60π ≈ 188.495559215388 (hand-derived)
func HollowCylinderSurfaceArea(outerRadius, innerRadius, height float64) (float64, error) {
	if err := checkRadii(outerRadius, innerRadius); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	return checkResult(product(4*math.Pi, outerRadius/2+innerRadius/2, height+(outerRadius-innerRadius), 1), "area")
}
