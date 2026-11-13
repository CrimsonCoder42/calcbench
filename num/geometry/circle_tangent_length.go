package geometry

import "fmt"

// CircleTangentLength returns the length of a tangent from an external point to a circle, from the
// circle's radius r and the point's distance d from the centre: t = √(d² − r²), computed as √(d−r)·√(d+r).
// That is Euclid's tangent² = (whole secant)·(its external part) with the secant through the centre, and
// the factored form keeps every digit for a point just outside the circle (d − r ≪ r), where d² − r²
// loses them all. A point on the circle gives 0. It knows no planet: a horizon distance passes its own
// cited radius.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// distanceFromCenter < radius (the point is inside the circle).
//
// ID:      NUM-GEOM-042
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Euclid, Elements, Book III, Proposition 36 (the square on the tangent equals the rectangle of the whole secant and its part outside the circle); D. E. Joyce (ed.), Euclid's Elements, Clark University, https://mathcs.clarku.edu/~djoyce/java/elements/bookIII/propIII36.html (read 2026-11-12)
// Tags:    geometry, plane, circle, tangent, horizon
// Units:   radius and distanceFromCenter in one length unit; result in that unit
// Example: 6371000, 6371002 → √(2·6371000·2 + 4) ≈ 5048.168380706 (hand-derived)
func CircleTangentLength(radius, distanceFromCenter float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	if err := checkPositive(distanceFromCenter, "distanceFromCenter"); err != nil {
		return 0, err
	}
	if distanceFromCenter < radius {
		return 0, fmt.Errorf("%w: distanceFromCenter < radius", ErrOutOfRange)
	}
	return legFromHypot(distanceFromCenter, radius), nil
}
