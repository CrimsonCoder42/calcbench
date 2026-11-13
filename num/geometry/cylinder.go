package geometry

import "math"

// checkRadiusHeight validates a radius and a height: both finite and > 0.
func checkRadiusHeight(radius, height float64) error {
	if err := checkPositive(radius, "radius"); err != nil {
		return err
	}
	return checkPositive(height, "height")
}

// CylinderVolume returns the volume of a circular cylinder from its radius and perpendicular height:
// πr²h. The same formula holds for an oblique cylinder.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the volume is beyond float64.
//
// ID:      NUM-GEOM-049
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.15 "Cylinders", p. 361 (right circular cylinder: volume πr²h; eq. (4.15.1) for the oblique case)
// Tags:    geometry, solid, cylinder
// Units:   radius and height in one length unit; result in that unit cubed
// Example: 2, 5 → 20π ≈ 62.831853071796 (hand-derived)
func CylinderVolume(radius, height float64) (float64, error) {
	if err := checkRadiusHeight(radius, height); err != nil {
		return 0, err
	}
	return checkResult(product(math.Pi, radius, radius, height), "volume")
}

// CylinderLateralArea returns the lateral (curved) area of a right circular cylinder: 2πrh.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-050
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.15 "Cylinders", p. 361 (right circular cylinder: lateral area 2πrh)
// Tags:    geometry, solid, cylinder
// Units:   radius and height in one length unit; result in that unit squared
// Example: 2, 5 → 20π ≈ 62.831853071796 (hand-derived)
func CylinderLateralArea(radius, height float64) (float64, error) {
	if err := checkRadiusHeight(radius, height); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, radius, height, 1), "area")
}

// CylinderSurfaceArea returns the total surface area of a right circular cylinder (the lateral area and
// both ends): 2πr(r + h).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-051
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.15 "Cylinders", p. 361 (right circular cylinder: total area 2πr(r + h))
// Tags:    geometry, solid, cylinder
// Units:   radius and height in one length unit; result in that unit squared
// Example: 2, 5 → 28π ≈ 87.964594300514 (hand-derived)
func CylinderSurfaceArea(radius, height float64) (float64, error) {
	if err := checkRadiusHeight(radius, height); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, radius, radius+height, 1), "area")
}
