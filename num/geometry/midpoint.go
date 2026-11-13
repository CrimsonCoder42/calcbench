package geometry

// Midpoint returns the midpoint of the segment pq, ((x₁+x₂)/2, (y₁+y₂)/2), computed as x₁/2 + x₂/2 so
// that it never overflows; halving is exact except for subnormal coordinates, which lose at most one ulp.
// Errors: ErrNotFinite when a coordinate is NaN or ±Inf.
//
// ID:      NUM-GEOM-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.4.2 "Distances", the midpoint of P₀P₁ (after eq. (4.4.8)), p. 316
// Tags:    geometry, plane, coordinates
// Units:   coordinates of p and q in one length unit; result in that unit
// Example: (0, 0), (4, −2) → (2, −1) (hand-derived)
func Midpoint(p, q Point) (Point, error) {
	if err := checkPoint(p, "p"); err != nil {
		return Point{}, err
	}
	if err := checkPoint(q, "q"); err != nil {
		return Point{}, err
	}
	return Point{X: p.X/2 + q.X/2, Y: p.Y/2 + q.Y/2}, nil
}
