package trig

import "math"

// SolveTriangleAAS returns all six parts of the triangle with angles A and B and side a (opposite A, so
// not between them): C = π − A − B, then b = a·sin B / sin A and c = a·sin C / sin A by the law of sines.
// C is formed with A + B exact (TwoSum) and π carried to about 107 bits, so it is accurate to a few ulps
// even when A or A + B is near π and C is tiny; sin C is taken as sin(A + B) when A + B < π/2, which
// stays accurate when C is near π.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when sideA ≤ 0; ErrDomain when an
// angle is outside (0, π) or angleARad + angleBRad ≥ π; ErrOverflow when a side is beyond float64.
//
// ID:      NUM-TRIG-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.148, p. 79
// Tags:    trig, triangle, solve, aas
// Units:   angleARad, angleBRad in radians; sideA in any length unit; result sides in that unit
// Example: A = atan2(3, 4), B = π/2, a = 3 → b = 5, c = 4 (hand-derived from 4.42.4)
func SolveTriangleAAS(angleARad, angleBRad, sideA float64) (Triangle, error) {
	if err := checkTwoAngles(angleARad, angleBRad); err != nil {
		return Triangle{}, err
	}
	if err := checkLength(sideA, "sideA"); err != nil {
		return Triangle{}, err
	}
	angleC, sinC := thirdAngle(angleARad, angleBRad)
	sinA := math.Sin(angleARad)
	t := Triangle{sideA, sideA * (math.Sin(angleBRad) / sinA), sideA * (sinC / sinA), angleARad, angleBRad, angleC}
	if err := checkOverflow(t); err != nil {
		return Triangle{}, err
	}
	return t, nil
}
