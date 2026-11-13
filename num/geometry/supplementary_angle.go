package geometry

import (
	"fmt"
	"math"
)

// SupplementaryAngleRad returns the supplement of an angle in radians, π − θ: the angle that adds to θ to
// make a straight angle. θ must lie in [0, π]; it is rejected outside, not wrapped.
// Errors: ErrNotFinite when angleRad is NaN or ±Inf; ErrOutOfRange when angleRad is outside [0, π].
//
// ID:      NUM-GEOM-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Marecek, Anthony-Smith, Honeycutt Mathis, Prealgebra, 2nd ed., OpenStax (Rice University), 2020, §9.3 "Use Properties of Angles, Triangles, and the Pythagorean Theorem" (supplementary angles: measures that add to 180°), https://openstax.org/books/prealgebra-2e/pages/9-3-use-properties-of-angles-triangles-and-the-pythagorean-theorem (read 2026-11-12)
// Tags:    geometry, plane, angle, supplementary
// Units:   angleRad and the result in radians
// Example: π/3 → 2π/3 (hand-derived)
func SupplementaryAngleRad(angleRad float64) (float64, error) {
	if err := checkFinite(angleRad, "angleRad"); err != nil {
		return 0, err
	}
	if angleRad < 0 || angleRad > math.Pi {
		return 0, fmt.Errorf("%w: angleRad", ErrOutOfRange)
	}
	return math.Pi - angleRad, nil
}
