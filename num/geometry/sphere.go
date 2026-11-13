package geometry

import "math"

// SphereVolume returns the volume of a sphere of the given radius: (4/3)πr³.
// Errors: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the
// volume is beyond float64.
//
// ID:      NUM-GEOM-063
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.18.1 "Spheres", item 3, p. 367 (the volume is (4/3)πr³)
// Tags:    geometry, solid, quadric, sphere
// Units:   radius in a length unit; result in that unit cubed
// Example: 3 → 36π ≈ 113.097335529233 (hand-derived)
func SphereVolume(radius float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	return checkResult(product(4*math.Pi/3, radius, radius, radius), "volume")
}

// SphereSurfaceArea returns the surface area of a sphere of the given radius: 4πr².
// Errors: ErrNotFinite when radius is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOverflow when the
// area is beyond float64.
//
// ID:      NUM-GEOM-064
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.18.1 "Spheres", item 3, p. 367 (the area of a sphere of radius r is 4πr²)
// Tags:    geometry, solid, quadric, sphere
// Units:   radius in a length unit; result in that unit squared
// Example: 3 → 36π ≈ 113.097335529233 (hand-derived)
func SphereSurfaceArea(radius float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	return checkResult(product(4*math.Pi, radius, radius, 1), "area")
}
