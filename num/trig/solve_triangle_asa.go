package trig

import "math"

// SolveTriangleASA returns all six parts of the triangle with angles A and B and the included side c:
// C = π − A − B, then a = c·sin A / sin C and b = c·sin B / sin C by the law of sines. C is formed with
// A + B exact (TwoSum) and π carried to about 107 bits, so it is accurate to a few ulps even when A or
// A + B is near π and C is tiny; sin C is taken as sin(A + B) when A + B < π/2, which stays accurate
// when C is near π.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when sideC ≤ 0; ErrDomain when an
// angle is outside (0, π) or angleARad + angleBRad ≥ π; ErrOverflow when a side is beyond float64.
//
// ID:      NUM-TRIG-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.148, p. 79
// Tags:    trig, triangle, solve, asa
// Units:   angleARad, angleBRad in radians; sideC in any length unit; result sides in that unit
// Example: A = B = 60°, c = 10 → the equilateral triangle, sides 10 (hand-derived from 4.42.4)
func SolveTriangleASA(angleARad, angleBRad, sideC float64) (Triangle, error) {
	if err := checkTwoAngles(angleARad, angleBRad); err != nil {
		return Triangle{}, err
	}
	if err := checkLength(sideC, "sideC"); err != nil {
		return Triangle{}, err
	}
	angleC, sinC := thirdAngle(angleARad, angleBRad)
	t := Triangle{sideC * (math.Sin(angleARad) / sinC), sideC * (math.Sin(angleBRad) / sinC), sideC, angleARad, angleBRad, angleC}
	if err := checkOverflow(t); err != nil {
		return Triangle{}, err
	}
	return t, nil
}
