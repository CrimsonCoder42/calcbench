package trig

import (
	"fmt"
	"math"
)

// SolveRightTriangleHypotenuseLeg returns the right triangle with hypotenuse c and leg a: C = π/2, the
// other leg b = √((c−a)(c+a)) (not √(c² − a²), which cancels when a ≈ c; sides beyond 2^±500 are scaled
// by an exact power of two first), A = atan2(a, b) and B = atan2(b, a).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a length ≤ 0; ErrNotTriangle
// when legA ≥ hypotenuse (legA = hypotenuse is the degenerate, flat case).
//
// ID:      NUM-TRIG-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.1 (Fig. 4.42.1, C the right angle) and 4.42.5 with C = π/2, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147, p. 78
// Tags:    trig, triangle, right-triangle, solve
// Units:   hypotenuse, legA in one length unit; result sides in that unit, angles in radians
// Example: (c 5, a 3) → b = 4, A = 0.6435011087932844 (hand-derived from 4.42.1)
func SolveRightTriangleHypotenuseLeg(hypotenuse, legA float64) (Triangle, error) {
	if err := checkLength(hypotenuse, "hypotenuse"); err != nil {
		return Triangle{}, err
	}
	if err := checkLength(legA, "legA"); err != nil {
		return Triangle{}, err
	}
	if legA >= hypotenuse {
		return Triangle{}, fmt.Errorf("%w: legA ≥ hypotenuse", ErrNotTriangle)
	}
	k := sideShift(hypotenuse)
	c, a := math.Ldexp(hypotenuse, k), math.Ldexp(legA, k)
	legB := math.Ldexp(math.Sqrt((c-a)*(c+a)), -k)
	return Triangle{legA, legB, hypotenuse, math.Atan2(legA, legB), math.Atan2(legB, legA), math.Pi / 2}, nil
}
