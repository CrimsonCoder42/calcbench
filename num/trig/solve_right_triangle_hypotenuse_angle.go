package trig

import "math"

// SolveRightTriangleHypotenuseAngle returns the right triangle with hypotenuse c and acute angle A:
// C = π/2, a = c·sin A, b = c·cos A (one math.Sincos) and B = π/2 − A.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when hypotenuse ≤ 0; ErrDomain when
// angleARad is outside (0, π/2).
//
// ID:      NUM-TRIG-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1–4.42.2 (Fig. 4.42.1, C the right angle), https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78
// Tags:    trig, triangle, right-triangle, solve
// Units:   hypotenuse in any length unit; angleARad in radians; result sides in that unit, angles in radians
// Example: c = 2, A = π/6 → a = 1, b = √3 (to rounding) (hand-derived from 4.42.1–4.42.2)
func SolveRightTriangleHypotenuseAngle(hypotenuse, angleARad float64) (Triangle, error) {
	if err := checkLength(hypotenuse, "hypotenuse"); err != nil {
		return Triangle{}, err
	}
	if err := checkAngle(angleARad, math.Pi/2, "angleARad"); err != nil {
		return Triangle{}, err
	}
	s, co := math.Sincos(angleARad)
	return Triangle{hypotenuse * s, hypotenuse * co, hypotenuse, angleARad, math.Pi/2 - angleARad, math.Pi / 2}, nil
}
