package geometry

import (
	"fmt"
	"math"
)

// SlopeAngleRad returns the angle in radians, in (−π, π], from the positive x-axis to the direction
// (run, rise): atan2(rise, run), whose tangent is the slope rise/run. It is exact in the quadrant and
// finite for a vertical line (run = 0). For a right triangle, rise = opposite and run = adjacent give
// the acute angle. −π (from rise = −0) is returned as π.
// Errors: ErrNotFinite when rise or run is NaN or ±Inf; ErrDegenerate when rise == run == 0.
//
// ID:      NUM-GEOM-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.4.3 "Angles", p. 316 (the angle of a line of slope m from the x-axis is tan⁻¹ m)
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.1.4.1 "Relations between Cartesian and polar coordinates", p. 303 (θ from x and y, sin θ = y/√(x²+y²), cos θ = x/√(x²+y²))
// Tags:    geometry, plane, coordinates, slope, rise-over-run, interface-angle
// Units:   rise and run in one length unit (or a ratio with run = 1); result in radians
// Example: rise 1, run 1 → π/4 (hand-derived)
func SlopeAngleRad(rise, run float64) (float64, error) {
	if err := checkFinite(rise, "rise"); err != nil {
		return 0, err
	}
	if err := checkFinite(run, "run"); err != nil {
		return 0, err
	}
	if rise == 0 && run == 0 {
		return 0, fmt.Errorf("%w: rise == run == 0", ErrDegenerate)
	}
	return halfOpenAngle(math.Atan2(rise, run)), nil
}

// halfOpenAngle maps the −π that math.Atan2 returns for a −0 ordinate to π, so angles lie in (−π, π].
func halfOpenAngle(theta float64) float64 {
	if theta == -math.Pi {
		return math.Pi
	}
	return theta
}
