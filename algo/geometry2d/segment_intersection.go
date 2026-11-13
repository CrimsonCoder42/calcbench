package geometry2d

// Segment is the closed line segment from A to B. A == B is a single point.
type Segment struct {
	A, B Point
}

// SegmentsIntersect reports whether the closed segments s and t share at least one point: a proper
// crossing, an endpoint touching the other segment, or a collinear overlap. A zero-length segment is a
// point and intersects when it lies on the other. Exact for every int64 coordinate.
//
// ID:      ALG-GEO2D-005
// Time:    O(1)
// Space:   O(1) extra, no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.1 (SEGMENTS-INTERSECT, ON-SEGMENT)
// Tags:    geometry, segment-intersection, predicate, exact
func SegmentsIntersect(s, t Segment) bool {
	d1 := Orientation(t.A, t.B, s.A)
	d2 := Orientation(t.A, t.B, s.B)
	d3 := Orientation(s.A, s.B, t.A)
	d4 := Orientation(s.A, s.B, t.B)
	if d1*d2 < 0 && d3*d4 < 0 {
		return true
	}
	return (d1 == 0 && inBox(t.A, t.B, s.A)) ||
		(d2 == 0 && inBox(t.A, t.B, s.B)) ||
		(d3 == 0 && inBox(s.A, s.B, t.A)) ||
		(d4 == 0 && inBox(s.A, s.B, t.B))
}
