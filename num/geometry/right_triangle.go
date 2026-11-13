package geometry

import (
	"fmt"
	"math"
)

// legFromHypot returns √(c² − a²) for finite 0 ≤ a ≤ c as √(c−a)·√(c+a): no cancellation when a ≈ c, and
// no c² to overflow. When c + a overflows, √(c+a) is taken as √2·√(c/2 + a/2).
func legFromHypot(c, a float64) float64 {
	s := c + a
	if math.IsInf(s, 1) {
		return math.Sqrt(c-a) * (math.Sqrt2 * math.Sqrt(c/2+a/2))
	}
	return math.Sqrt(c-a) * math.Sqrt(s)
}

// RightTriangleHypotenuse returns the hypotenuse of a right triangle from its two legs, √(a² + b²),
// computed with math.Hypot so the squares never overflow or underflow.
// Errors: ErrNotFinite when a leg is NaN or ±Inf; ErrNonPositive when a leg ≤ 0; ErrOverflow when the
// hypotenuse is beyond float64.
//
// ID:      NUM-GEOM-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.1 "Triangles", p. 320 (right triangle: c² = a² + b², Pythagoras)
// Tags:    geometry, plane, triangle, pythagorean, right-triangle
// Units:   legs in one length unit; result in that unit
// Example: 3, 4 → 5 (hand-derived)
func RightTriangleHypotenuse(legA, legB float64) (float64, error) {
	if err := checkPositive(legA, "legA"); err != nil {
		return 0, err
	}
	if err := checkPositive(legB, "legB"); err != nil {
		return 0, err
	}
	return checkResult(math.Hypot(legA, legB), "hypotenuse")
}

// RightTriangleLeg returns the other leg of a right triangle from its hypotenuse and one leg,
// √(c² − a²), computed as √(c−a)·√(c+a): exact cancellation when the leg is close to the hypotenuse, and
// no overflow of c².
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrNotTriangle
// when leg ≥ hypotenuse.
//
// ID:      NUM-GEOM-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.1 "Triangles", p. 320 (right triangle: c² = a² + b², Pythagoras)
// Tags:    geometry, plane, triangle, pythagorean, right-triangle
// Units:   hypotenuse and leg in one length unit; result in that unit
// Example: 5, 3 → 4 (hand-derived)
func RightTriangleLeg(hypotenuse, leg float64) (float64, error) {
	if err := checkPositive(hypotenuse, "hypotenuse"); err != nil {
		return 0, err
	}
	if err := checkPositive(leg, "leg"); err != nil {
		return 0, err
	}
	if leg >= hypotenuse {
		return 0, fmt.Errorf("%w: leg >= hypotenuse", ErrNotTriangle)
	}
	return legFromHypot(hypotenuse, leg), nil
}
