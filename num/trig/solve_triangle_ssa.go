package trig

import "math"

// SolveTriangleSSA returns the triangles with sides a and b and the non-included angle A (opposite a),
// the ambiguous case: count is 0, 1 or 2 and solutions[:count] holds them, solutions[0] being the one
// with the acute (smaller) B; unused slots are the zero Triangle. count 0 (no triangle) is not an error.
// B is found as by LawOfSinesAngle (Kahan's B(a, A, b)), then C = π − A − B and c = a·sin C / sin A. C is
// formed with A + B exact (TwoSum) and π carried to about 107 bits, so it is accurate when A is near π and
// C is tiny; sin C is taken as sin(A + B) when A + B < π/2, which stays accurate when C is near π. For
// the second triangle (B2 = π − B1), C is taken as B1 − A, which is exact in real arithmetic.
// The result is a fixed array, not a slice, so nothing is allocated.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrDomain when
// angleARad is outside (0, π); ErrOverflow when c is beyond float64.
//
// ID:      NUM-TRIG-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.4, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §8 (B(a, A, b)), p. 6, and §9 Table 2, p. 8, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Tags:    trig, triangle, solve, ssa, ambiguous-case
// Units:   sideA, sideB in one length unit; angleARad in radians; result sides in that unit, angles in radians
// Example: a = 3, b = 17, A = 10° → two triangles, B = 79.73894317° or 100.26105683° (Kahan 2014 Table 2, row 1)
func SolveTriangleSSA(sideA, sideB, angleARad float64) (solutions [2]Triangle, count int, err error) {
	b1, b2, n, err := LawOfSinesAngle(sideA, sideB, angleARad)
	if err != nil {
		return solutions, 0, err
	}
	sinA := math.Sin(angleARad)
	for i := range n {
		var angleB, angleC, sinC float64
		if i == 0 {
			angleB = b1
			angleC, sinC = thirdAngle(angleARad, b1)
		} else {
			angleB = b2
			angleC = b1 - angleARad // = π − A − (π − B1), exact in real arithmetic
			sinC = math.Sin(angleC)
		}
		t := Triangle{sideA, sideB, sideA * (sinC / sinA), angleARad, angleB, angleC}
		if err := checkOverflow(t); err != nil {
			return [2]Triangle{}, 0, err
		}
		solutions[i] = t
	}
	return solutions, n, nil
}
