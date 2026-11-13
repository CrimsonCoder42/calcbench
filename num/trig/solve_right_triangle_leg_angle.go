package trig

import "math"

// SolveRightTriangleLegAngle returns the right triangle with leg a and its opposite acute angle A:
// C = π/2, b = a·cos A / sin A, c = a / sin A (one math.Sincos) and B = π/2 − A. For the angle adjacent
// to the leg, pass π/2 minus it.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when legA ≤ 0; ErrDomain when
// angleARad is outside (0, π/2); ErrOverflow when a side is beyond float64.
//
// ID:      NUM-TRIG-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1 and 4.42.3 (Fig. 4.42.1, C the right angle), https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78
// Tags:    trig, triangle, right-triangle, solve
// Units:   legA in any length unit; angleARad in radians; result sides in that unit, angles in radians
// Example: a = 3, A = atan2(3, 4) → b = 4, c = 5 (hand-derived from 4.42.1 and 4.42.3)
func SolveRightTriangleLegAngle(legA, angleARad float64) (Triangle, error) {
	if err := checkLength(legA, "legA"); err != nil {
		return Triangle{}, err
	}
	if err := checkAngle(angleARad, math.Pi/2, "angleARad"); err != nil {
		return Triangle{}, err
	}
	s, co := math.Sincos(angleARad)
	t := Triangle{legA, legA * (co / s), legA / s, angleARad, math.Pi/2 - angleARad, math.Pi / 2}
	if err := checkOverflow(t); err != nil {
		return Triangle{}, err
	}
	return t, nil
}
