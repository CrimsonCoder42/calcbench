package trig

import "math"

// SolveRightTriangleLegs returns the right triangle with legs a and b: C = π/2, the hypotenuse
// c = √(a² + b²) by math.Hypot (no overflow in the squares), A = atan2(a, b) and B = atan2(b, a), each
// accurate on its own (B is not formed as π/2 − A).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a leg ≤ 0; ErrOverflow when the
// hypotenuse is beyond float64.
//
// ID:      NUM-TRIG-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1–4.42.3 (Fig. 4.42.1, C the right angle) and 4.42.5 with C = π/2, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78
// Tags:    trig, triangle, right-triangle, solve
// Units:   legA, legB in one length unit; result sides in that unit, angles in radians
// Example: (3, 4) → c = 5, A = 0.6435011087932844 (36.8698976458°) (hand-derived from 4.42.3)
func SolveRightTriangleLegs(legA, legB float64) (Triangle, error) {
	if err := checkLength(legA, "legA"); err != nil {
		return Triangle{}, err
	}
	if err := checkLength(legB, "legB"); err != nil {
		return Triangle{}, err
	}
	t := Triangle{legA, legB, math.Hypot(legA, legB), math.Atan2(legA, legB), math.Atan2(legB, legA), math.Pi / 2}
	if err := checkOverflow(t); err != nil {
		return Triangle{}, err
	}
	return t, nil
}
