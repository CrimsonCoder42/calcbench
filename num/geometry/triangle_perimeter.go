package geometry

import (
	"errors"
	"fmt"
)

var (
	// ErrNonPositive: a length, area or count that must be > 0 is ≤ 0.
	ErrNonPositive = errors.New("geometry: value must be positive")
	// ErrNotTriangle: the lengths violate the triangle inequality (or a leg is not shorter than the
	// hypotenuse).
	ErrNotTriangle = errors.New("geometry: lengths do not form a triangle")
)

// checkPositive returns ErrNotFinite or ErrNonPositive wrapped with name unless x is finite and > 0.
func checkPositive(x float64, name string) error {
	if err := checkFinite(x, name); err != nil {
		return err
	}
	if x <= 0 {
		return fmt.Errorf("%w: %s", ErrNonPositive, name)
	}
	return nil
}

// checkNonNegative returns ErrNotFinite or ErrNegative wrapped with name unless x is finite and ≥ 0.
func checkNonNegative(x float64, name string) error {
	if err := checkFinite(x, name); err != nil {
		return err
	}
	if x < 0 {
		return fmt.Errorf("%w: %s", ErrNegative, name)
	}
	return nil
}

// triangleSides validates three side lengths and returns them sorted so that a ≥ b ≥ c. It applies
// Kahan's test: the sides form a (possibly degenerate) triangle exactly when c − (a − b) ≥ 0, a test that
// roundoff cannot get wrong because a − b is exact whenever it matters (Kahan 2014 §2, §4).
func triangleSides(sideA, sideB, sideC float64) (a, b, c float64, err error) {
	if err := checkPositive(sideA, "sideA"); err != nil {
		return 0, 0, 0, err
	}
	if err := checkPositive(sideB, "sideB"); err != nil {
		return 0, 0, 0, err
	}
	if err := checkPositive(sideC, "sideC"); err != nil {
		return 0, 0, 0, err
	}
	a, b, c = sideA, sideB, sideC
	if a < b {
		a, b = b, a
	}
	if b < c {
		b, c = c, b
	}
	if a < b {
		a, b = b, a
	}
	if c-(a-b) < 0 {
		return 0, 0, 0, fmt.Errorf("%w: sideA, sideB, sideC", ErrNotTriangle)
	}
	return a, b, c, nil
}

// TrianglePerimeter returns the perimeter a + b + c of the triangle with the given sides, after checking
// that they form a triangle with the same test as TriangleAreaHeron, so the two agree on what a triangle
// is. A degenerate triangle (one side equal to the sum of the others) is accepted.
// Errors: ErrNotFinite when a side is NaN or ±Inf; ErrNonPositive when a side ≤ 0; ErrNotTriangle when
// the triangle inequality fails; ErrOverflow when the perimeter is beyond float64.
//
// ID:      NUM-GEOM-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.1 "Triangles", p. 319 (semi-perimeter s = ½(a + b + c)) and p. 321 (additional fact 4: each side is less than the sum of the other two)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §2 "How to compute ∆" (sort a ≥ b ≥ c; the data are not a triangle if c − (a − b) < 0), p. 3, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
// Tags:    geometry, plane, triangle
// Units:   sides in one length unit; result in that unit
// Example: 3, 4, 5 → 12 (hand-derived)
func TrianglePerimeter(sideA, sideB, sideC float64) (float64, error) {
	a, b, c, err := triangleSides(sideA, sideB, sideC)
	if err != nil {
		return 0, err
	}
	return checkResult(a+(b+c), "perimeter")
}
