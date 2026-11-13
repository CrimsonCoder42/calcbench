package trig

import (
	"fmt"
	"math"
)

// Triangle holds the six parts of a plane triangle. Angle A is opposite side a, B opposite b, C opposite
// c (DLMF Fig. 4.42.2); in a right triangle C is the right angle and c the hypotenuse (DLMF Fig. 4.42.1).
// Sides are in the caller's length unit; angles are in radians.
type Triangle struct {
	SideA, SideB, SideC             float64
	AngleARad, AngleBRad, AngleCRad float64
}

// SolveTriangleSSS returns all six parts of the triangle with sides a, b and c. Each angle is computed on
// its own by Kahan's formula (as LawOfCosinesAngle), so each is accurate even for needle-like triangles;
// they sum to π within a few ulps. A flat triangle is allowed (angles 0, 0 and π).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrNotTriangle when
// one side exceeds the sum of the other two.
//
// ID:      NUM-TRIG-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §3, p. 3, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Tags:    trig, triangle, solve, sss
// Units:   sideA, sideB, sideC in one length unit; angles in radians
// Example: (10, 10, 10) → all angles 60° (Kahan 2014 Table 1, row 1); (3, 4, 5) → C = π/2, A = 0.6435011087932844 (hand-derived)
func SolveTriangleSSS(sideA, sideB, sideC float64) (Triangle, error) {
	if err := checkSides(sideA, sideB, sideC); err != nil {
		return Triangle{}, err
	}
	angleC, ok := cosinesAngle(sideA, sideB, sideC, false)
	if !ok {
		return Triangle{}, fmt.Errorf("%w: sideA, sideB, sideC", ErrNotTriangle)
	}
	// The triangle inequality holds, so the other two cannot fail; clamp guards against rounding.
	angleA, _ := cosinesAngle(sideB, sideC, sideA, true)
	angleB, _ := cosinesAngle(sideA, sideC, sideB, true)
	return Triangle{sideA, sideB, sideC, angleA, angleB, angleC}, nil
}

// checkOverflow returns ErrOverflow wrapped with the first of a triangle's sides that is ±Inf.
func checkOverflow(t Triangle) error {
	switch {
	case math.IsInf(t.SideA, 0):
		return fmt.Errorf("%w: sideA", ErrOverflow)
	case math.IsInf(t.SideB, 0):
		return fmt.Errorf("%w: sideB", ErrOverflow)
	case math.IsInf(t.SideC, 0):
		return fmt.Errorf("%w: sideC", ErrOverflow)
	}
	return nil
}
