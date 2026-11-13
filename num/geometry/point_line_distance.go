package geometry

import (
	"fmt"
	"math"
)

// PointLineDistance returns the perpendicular distance from p to the infinite line through a and b:
// |(b−a)×(p−a)| / |b−a|. The direction b−a is first scaled to a unit vector, so no product overflows,
// and the cross-product form keeps its accuracy when the line is far from the origin (the a·x+b·y+c form
// does not). When a coordinate difference overflows, the halved coordinates are used and the result
// doubled. The distance to the segment ab (clamped) is a different quantity.
// Errors: ErrNotFinite when a coordinate is NaN or ±Inf; ErrDegenerate when a == b; ErrOverflow when the
// distance is beyond float64.
//
// ID:      NUM-GEOM-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.4.2 "Distances", eq. (4.4.9), p. 316 (distance from a point to a line), with the line through two points of §4.4.1, eq. (4.4.3), p. 315
// Tags:    geometry, plane, coordinates, distance, line
// Units:   coordinates of p, a, b in one length unit; result in that unit
// Example: p (0, 0), a (0, 1), b (1, 0) → √2/2 (hand-derived)
func PointLineDistance(p, a, b Point) (float64, error) {
	if err := checkPoint(p, "p"); err != nil {
		return 0, err
	}
	if err := checkPoint(a, "a"); err != nil {
		return 0, err
	}
	if err := checkPoint(b, "b"); err != nil {
		return 0, err
	}
	if a == b {
		return 0, fmt.Errorf("%w: a == b", ErrDegenerate)
	}
	ux, uy := b.X-a.X, b.Y-a.Y
	wx, wy := p.X-a.X, p.Y-a.Y
	scale := 1.0
	if math.IsInf(ux, 0) || math.IsInf(uy, 0) || math.IsInf(wx, 0) || math.IsInf(wy, 0) {
		ux, uy = b.X/2-a.X/2, b.Y/2-a.Y/2
		wx, wy = p.X/2-a.X/2, p.Y/2-a.Y/2
		scale = 2
	}
	ex, ey := unitVector(ux, uy)
	return checkResult(math.Abs(ex*wy-ey*wx)*scale, "distance")
}

// unitVector returns (x, y) scaled to length 1. (x, y) must be finite and not (0, 0). Dividing by the
// larger component first keeps math.Hypot's argument in [1, √2], so nothing overflows or underflows.
func unitVector(x, y float64) (ux, uy float64) {
	m := math.Max(math.Abs(x), math.Abs(y))
	x, y = x/m, y/m
	n := math.Hypot(x, y)
	return x / n, y / n
}
