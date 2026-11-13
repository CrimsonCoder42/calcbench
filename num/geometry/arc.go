package geometry

import (
	"fmt"
	"math"
)

// checkCentralAngle returns ErrNotFinite or ErrOutOfRange wrapped with name unless angleRad is finite
// and lies in (0, 2π]. Angles beyond a full turn are rejected, not wrapped.
func checkCentralAngle(angleRad float64, name string) error {
	if err := checkFinite(angleRad, name); err != nil {
		return err
	}
	if angleRad <= 0 || angleRad > 2*math.Pi {
		return fmt.Errorf("%w: %s", ErrOutOfRange, name)
	}
	return nil
}

// ArcLength returns the length of a circular arc of the given radius subtending the given central angle:
// r·θ.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when
// angleRad is outside (0, 2π]; ErrOverflow when the length is beyond float64.
//
// ID:      NUM-GEOM-037
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 ("The length of the arc of circle subtended by an angle θ … is rθ")
// Tags:    geometry, plane, circle, arc
// Units:   radius in a length unit; angleRad in radians; result in radius's unit
// Example: 2, π/3 → 2π/3 ≈ 2.094395102393 (hand-derived)
func ArcLength(radius, angleRad float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	if err := checkCentralAngle(angleRad, "angleRad"); err != nil {
		return 0, err
	}
	return checkResult(product(radius, angleRad, 1, 1), "length")
}

// ChordLength returns the length of the chord of a circle of the given radius subtending the given central
// angle: 2r·sin(θ/2). The half-angle sine form keeps full accuracy for small θ, where the law-of-cosines
// form √(2r²(1 − cos θ)) loses digits.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when
// centralAngleRad is outside (0, 2π]; ErrOverflow when the chord is beyond float64.
//
// ID:      NUM-GEOM-038
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 (c = 2R sin ½θ)
// Tags:    geometry, plane, circle, chord
// Units:   radius in a length unit; centralAngleRad in radians; result in radius's unit
// Example: 2, π/3 → 2 (hand-derived: an equilateral triangle)
func ChordLength(radius, centralAngleRad float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	if err := checkCentralAngle(centralAngleRad, "centralAngleRad"); err != nil {
		return 0, err
	}
	return checkResult(product(2, radius, math.Sin(centralAngleRad/2), 1), "chord")
}

// SectorArea returns the area of a circular sector of the given radius and central angle: ½·r²·θ.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when
// angleRad is outside (0, 2π]; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-039
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 (area of sector = ½Rs = ½R²θ)
// Tags:    geometry, plane, circle, sector
// Units:   radius in a length unit; angleRad in radians; result in that unit squared
// Example: 2, π/3 → 2π/3 ≈ 2.094395102393 (hand-derived)
func SectorArea(radius, angleRad float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	if err := checkCentralAngle(angleRad, "angleRad"); err != nil {
		return 0, err
	}
	return checkResult(product(0.5, radius, radius, angleRad), "area")
}
