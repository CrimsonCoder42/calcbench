package geometry

import (
	"errors"
	"fmt"
	"math"
)

// ErrOutOfRange: an angle, ratio or length is outside the range the formula is defined on.
var ErrOutOfRange = errors.New("geometry: value outside the formula's range")

// ComplementaryAngleRad returns the complement of an angle in radians, π/2 − θ: the angle that adds to θ
// to make a right angle. θ must lie in [0, π/2]; negative and obtuse angles are rejected, not wrapped.
// Errors: ErrNotFinite when angleRad is NaN or ±Inf; ErrOutOfRange when angleRad is outside [0, π/2].
//
// ID:      NUM-GEOM-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Marecek, Anthony-Smith, Honeycutt Mathis, Prealgebra, 2nd ed., OpenStax (Rice University), 2020, §9.3 "Use Properties of Angles, Triangles, and the Pythagorean Theorem" (complementary angles: measures that add to 90°), https://openstax.org/books/prealgebra-2e/pages/9-3-use-properties-of-angles-triangles-and-the-pythagorean-theorem (read 2026-11-12)
// Tags:    geometry, plane, angle, complementary
// Units:   angleRad and the result in radians
// Example: π/6 → π/3 (hand-derived)
func ComplementaryAngleRad(angleRad float64) (float64, error) {
	if err := checkFinite(angleRad, "angleRad"); err != nil {
		return 0, err
	}
	if angleRad < 0 || angleRad > math.Pi/2 {
		return 0, fmt.Errorf("%w: angleRad", ErrOutOfRange)
	}
	return math.Pi/2 - angleRad, nil
}
