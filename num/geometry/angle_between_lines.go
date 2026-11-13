package geometry

import (
	"fmt"
	"math"
)

// AngleBetweenLinesRad returns the acute angle in radians, in [0, π/2], between the line through a and b
// and the line through c and d: atan2(|u×v|, |u·v|) with u = b−a and v = d−c. This is the source's
// tan θ = (a₀b₁ − a₁b₀)/(a₀a₁ + b₀b₁) in vector form; unlike the form in slopes, it is finite for vertical
// and perpendicular lines, and unlike acos(u·v/|u||v|) it keeps full accuracy near 0 and π/2. Each
// direction is scaled by a power of two (exact) so its larger component lies in [0.5, 1): no product
// overflows, and exactly parallel lines give exactly 0.
// Errors: ErrNotFinite when a coordinate is NaN or ±Inf; ErrDegenerate when a == b or c == d.
//
// ID:      NUM-GEOM-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.4.3 "Angles", eq. (4.4.10), p. 316 (angle between two lines; the form in slopes, tan⁻¹((m₁−m₀)/(1+m₀m₁)), on the same page)
// Tags:    geometry, plane, coordinates, interface-angle, line
// Units:   coordinates in one length unit; result in radians
// Example: the line y = 0 and the line through (0, 0), (√3, 1) → π/6 (hand-derived)
func AngleBetweenLinesRad(a, b, c, d Point) (float64, error) {
	for _, pt := range [...]struct {
		p    Point
		name string
	}{{a, "a"}, {b, "b"}, {c, "c"}, {d, "d"}} {
		if err := checkPoint(pt.p, pt.name); err != nil {
			return 0, err
		}
	}
	if a == b {
		return 0, fmt.Errorf("%w: a == b", ErrDegenerate)
	}
	if c == d {
		return 0, fmt.Errorf("%w: c == d", ErrDegenerate)
	}
	ux, uy := direction(a, b)
	vx, vy := direction(c, d)
	cross := ux*vy - uy*vx
	dot := ux*vx + uy*vy
	return math.Atan2(math.Abs(cross), math.Abs(dot)), nil
}

// direction returns the vector from p to q (p ≠ q, both finite) scaled by a power of two so that its
// larger component has magnitude in [0.5, 1). The scaling is exact, so parallel inputs stay exactly
// parallel. When a difference overflows, the halved coordinates give the same direction.
func direction(p, q Point) (x, y float64) {
	x, y = q.X-p.X, q.Y-p.Y
	if math.IsInf(x, 0) || math.IsInf(y, 0) {
		x, y = q.X/2-p.X/2, q.Y/2-p.Y/2
	}
	_, e := math.Frexp(math.Max(math.Abs(x), math.Abs(y)))
	return math.Ldexp(x, -e), math.Ldexp(y, -e)
}
