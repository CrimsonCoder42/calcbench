package geometry

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrDegenerate: the figure collapses (coincident points, a zero-length line, fewer vertices than the
	// figure needs, a polygon of zero area).
	ErrDegenerate = errors.New("geometry: degenerate figure")
	// ErrVertical: the slope of a vertical line is undefined.
	ErrVertical = errors.New("geometry: slope of a vertical line is undefined")
)

// Slope returns the slope (rise over run) of the line through p and q, (y₂−y₁)/(x₂−x₁). A vertical line
// has no slope: it is an error, never ±Inf (SlopeAngleRad handles vertical lines). When a coordinate
// difference overflows, both differences are taken of the halved coordinates; the ratio is unchanged.
// Errors: ErrNotFinite when a coordinate is NaN or ±Inf; ErrDegenerate when p == q; ErrVertical when
// x₁ == x₂ and y₁ ≠ y₂; ErrOverflow when the slope is beyond float64.
//
// ID:      NUM-GEOM-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.4.1 "Lines with prescribed properties", item 6 (slope of the line through two points), p. 315
// Tags:    geometry, plane, coordinates, slope, rise-over-run
// Units:   coordinates of p and q in one length unit; result dimensionless
// Example: (1, 1), (3, 5) → 2 (hand-derived)
func Slope(p, q Point) (float64, error) {
	if err := checkPoint(p, "p"); err != nil {
		return 0, err
	}
	if err := checkPoint(q, "q"); err != nil {
		return 0, err
	}
	if p == q {
		return 0, fmt.Errorf("%w: p == q", ErrDegenerate)
	}
	if p.X == q.X {
		return 0, fmt.Errorf("%w: q", ErrVertical)
	}
	dx, dy := q.X-p.X, q.Y-p.Y
	if math.IsInf(dx, 0) || math.IsInf(dy, 0) {
		dx, dy = q.X/2-p.X/2, q.Y/2-p.Y/2
	}
	return checkResult(dy/dx, "slope")
}
