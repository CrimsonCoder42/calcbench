package trig

import "math"

// SolveTriangleSAS returns all six parts of the triangle with sides a and b and the included angle C:
// c by Kahan's form of the law of cosines (as LawOfCosinesSide). The angle opposite the shorter of a and
// b (say b ≤ a) is B = atan2(b·sin C, (a − b) + 2b·sin²(C/2)), from c·sin B = b·sin C (law of sines) and
// c·cos B = a − b·cos C (projection formula) with 1 − cos C = 2 sin²(C/2); every term is non-negative, so
// nothing cancels and B is accurate even when C is near π or near 0, and it does not use the rounded c.
// The other angle is A = π − C − B, formed with C + B exact and π carried to about 107 bits, so it too is
// accurate when C is near π. The law of sines alone is not used for an angle: it cannot tell B from π − B.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrDomain when
// angleCRad is outside (0, π); ErrOverflow when c is beyond float64.
//
// ID:      NUM-TRIG-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eqs. 4.42.4 (law of sines) and 4.42.6 (projection formula), which give tan B = b sin C / (a − b cos C), https://dlmf.nist.gov/4.42 (read 2026-11-12); eq. 4.21.28 (cos 2z = 1 − 2 sin² z), https://dlmf.nist.gov/4.21 (read 2026-11-12)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §7 (side c), p. 6, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Tags:    trig, triangle, solve, sas
// Units:   sideA, sideB in one length unit; angleCRad in radians; result sides in that unit, angles in radians
// Example: (3, 4, C = 90°) → c = 5, A = 0.6435011087932844, B = π/2 − A (hand-derived)
func SolveTriangleSAS(sideA, sideB, angleCRad float64) (Triangle, error) {
	c, err := LawOfCosinesSide(sideA, sideB, angleCRad)
	if err != nil {
		return Triangle{}, err
	}
	big, small := sideA, sideB
	if small > big {
		big, small = small, big
	}
	if k := sideShift(big); k != 0 {
		big, small = math.Ldexp(big, k), math.Ldexp(small, k)
	}
	h := math.Sin(angleCRad / 2)
	// The angle opposite the shorter side: atan2(c·sin, c·cos) with both legs free of cancellation.
	minor := math.Atan2(small*math.Sin(angleCRad), (big-small)+2*small*h*h)
	major, _ := thirdAngle(angleCRad, minor)
	angleA, angleB := major, minor
	switch {
	case sideA == sideB:
		angleA = minor
	case sideA < sideB:
		angleA, angleB = minor, major
	}
	return Triangle{sideA, sideB, c, angleA, angleB, angleCRad}, nil
}
