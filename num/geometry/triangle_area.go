package geometry

import (
	"fmt"
	"math"
)

// product returns x·y·z·w for finite non-negative factors without a spurious intermediate overflow or
// underflow: the mantissas (each in [½, 1)) and the exponents are multiplied separately, so the result is
// +Inf only when the true product is beyond float64 and 0 only when it is below the smallest subnormal.
// The three mantissa products round as the plain product would. Pass 1 for unused factors.
func product(x, y, z, w float64) float64 {
	mx, ex := math.Frexp(x)
	my, ey := math.Frexp(y)
	mz, ez := math.Frexp(z)
	mw, ew := math.Frexp(w)
	return math.Ldexp(mx*my*mz*mw, ex+ey+ez+ew)
}

// checkOpenAngle returns ErrNotFinite or ErrOutOfRange wrapped with name unless angleRad is finite and
// lies in the open interval (0, upperRad).
func checkOpenAngle(angleRad, upperRad float64, name string) error {
	if err := checkFinite(angleRad, name); err != nil {
		return err
	}
	if angleRad <= 0 || angleRad >= upperRad {
		return fmt.Errorf("%w: %s", ErrOutOfRange, name)
	}
	return nil
}

// TriangleAreaBaseHeight returns the area of a triangle from a base and the height perpendicular to it:
// ½·base·height.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the area is beyond float64.
//
// ID:      NUM-GEOM-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.1 "Triangles", p. 319 (area = ½h_c·c)
// Tags:    geometry, plane, triangle
// Units:   base and height in one length unit; result in that unit squared
// Example: base 6, height 4 → 12 (hand-derived)
func TriangleAreaBaseHeight(base, height float64) (float64, error) {
	if err := checkPositive(base, "base"); err != nil {
		return 0, err
	}
	if err := checkPositive(height, "height"); err != nil {
		return 0, err
	}
	return checkResult(product(0.5, base, height, 1), "area")
}

// TriangleAreaHeron returns the area of the triangle with the given three sides by Kahan's stable form of
// Heron's formula: with a ≥ b ≥ c, ∆ = ¼√((a+(b+c))(c−(a−b))(c+(a−b))(a+(b−c))). The parentheses are
// essential: unlike √(s(s−a)(s−b)(s−c)), which loses every digit for needle-shaped triangles, this is
// accurate to a few ulps for every triangle. The sides are scaled by a power of two (exact) so the factors
// cannot overflow, and the square root is taken of (first·fourth factor) and of the two middle factors
// separately so a needle triangle's tiny factors cannot underflow. A degenerate triangle returns 0.
// Errors: ErrNotFinite when a side is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrNotTriangle when
// c − (a − b) < 0; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §2 "How to compute ∆", p. 3, and §5 Table 1 "Area ∆ and Angle C", p. 4, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.1 "Triangles", p. 319 (area = √(s(s−a)(s−b)(s−c)), Heron)
// Tags:    geometry, plane, triangle, heron, needle-triangle, education
// Units:   sides in one length unit; result in that unit squared
// Example: (10, 10, 10) → 43.30127020; (100000, 99999.99979, 0.00029) → 9.999999990; (10000, 5000.000001, 15000) → 612.3724358 (Kahan 2014, Table 1, "Accurate ∆")
func TriangleAreaHeron(sideA, sideB, sideC float64) (float64, error) {
	a, b, c, err := triangleSides(sideA, sideB, sideC)
	if err != nil {
		return 0, err
	}
	_, e := math.Frexp(a)
	a, b, c = math.Ldexp(a, -e), math.Ldexp(b, -e), math.Ldexp(c, -e)
	outer := (a + (b + c)) * (a + (b - c))
	area := 0.25 * math.Sqrt(outer) * math.Sqrt(c-(a-b)) * math.Sqrt(c+(a-b))
	return checkResult(math.Ldexp(area, 2*e), "area")
}

// TriangleAreaSAS returns the area of a triangle from two sides and the angle between them:
// ½·a·b·sin C.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrOutOfRange when
// includedAngleRad is outside (0, π); ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.1 "Triangles", p. 319 (area = ½ab sin C)
// Tags:    geometry, plane, triangle
// Units:   sides in one length unit; includedAngleRad in radians; result in that unit squared
// Example: 3, 4, π/6 → 3 (hand-derived)
func TriangleAreaSAS(sideA, sideB, includedAngleRad float64) (float64, error) {
	if err := checkPositive(sideA, "sideA"); err != nil {
		return 0, err
	}
	if err := checkPositive(sideB, "sideB"); err != nil {
		return 0, err
	}
	if err := checkOpenAngle(includedAngleRad, math.Pi, "includedAngleRad"); err != nil {
		return 0, err
	}
	return checkResult(product(0.5, sideA, sideB, math.Sin(includedAngleRad)), "area")
}
