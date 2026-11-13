package trig

import (
	"fmt"
	"math"
)

// LawOfCosinesSide returns side c of a triangle from sides a and b and the included angle C (opposite c):
// c² = a² + b² − 2ab cos C, evaluated in Kahan's cancellation-free form c = √((a−b)² + 4ab sin²(C/2)),
// as hypot(|a−b|, 2 sin(C/2)·√a·√b) so that ab is never formed and cannot overflow. The textbook form
// loses digits when c is much smaller than a and b.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrDomain when
// angleCRad is outside (0, π); ErrOverflow when c is beyond float64.
//
// ID:      NUM-TRIG-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §7 ("c := √((a−b)² + 4 a b sin²(C/2))"), p. 6, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Tags:    trig, triangle, law-of-cosines
// Units:   sideA, sideB in one length unit; angleCRad in radians; result in that length unit
// Example: a = b = 10, C = 60° → 10; a = 3, b = 4, C = 90° → 5 (hand-derived from 4.42.5)
func LawOfCosinesSide(sideA, sideB, angleCRad float64) (float64, error) {
	if err := checkLength(sideA, "sideA"); err != nil {
		return 0, err
	}
	if err := checkLength(sideB, "sideB"); err != nil {
		return 0, err
	}
	if err := checkAngle(angleCRad, math.Pi, "angleCRad"); err != nil {
		return 0, err
	}
	c := cosinesSide(sideA, sideB, angleCRad)
	if math.IsInf(c, 0) {
		return 0, fmt.Errorf("%w: sideC", ErrOverflow)
	}
	return c, nil
}

// cosinesSide is Kahan 2014 §7 on validated inputs. It returns +Inf when c is beyond float64.
func cosinesSide(a, b, angleC float64) float64 {
	return math.Hypot(math.Abs(a-b), 2*math.Sin(angleC/2)*math.Sqrt(a)*math.Sqrt(b))
}
