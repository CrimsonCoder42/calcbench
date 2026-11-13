package geometry2d

import "fmt"

// Location is where a point lies relative to a polygon.
type Location int

const (
	// Outside: the point is not in the closed polygon.
	Outside Location = iota
	// Inside: the point is in the polygon's interior.
	Inside
	// OnBoundary: the point lies on an edge or is a vertex.
	OnBoundary
)

// String returns "Outside", "Inside" or "OnBoundary".
func (l Location) String() string {
	switch l {
	case Outside:
		return "Outside"
	case Inside:
		return "Inside"
	case OnBoundary:
		return "OnBoundary"
	}
	return fmt.Sprintf("Location(%d)", int(l))
}

// PointInPolygon returns whether q is Inside, Outside or OnBoundary of the polygon whose vertices are
// given in order (either orientation; the last vertex joins the first). The boundary is decided exactly
// and first; otherwise the crossing number of a ray from q towards +X decides, with the half-open rule at
// vertices (an edge counts when one endpoint is strictly above q and the other is not). A
// self-intersecting polygon is read by the even–odd rule. polygon is read, never modified.
// Errors: ErrTooFewPoints when len(polygon) < 3.
//
// ID:      ALG-GEO2D-004
// Time:    O(n)   n = len(polygon); Θ(n) for one query without preprocessing
// Space:   O(1) extra, no allocation
// Source:  Haines, "Point in Polygon Strategies", in Heckbert (ed.), Graphics Gems IV, Academic Press, 1994, pp. 24–46 (crossings test)
// Source:  O'Rourke, Computational Geometry in C, 2nd ed., Cambridge University Press, 1998, §7.4.2 (ray crossings)
// Tags:    geometry, point-in-polygon, crossing-number, ray-casting, even-odd, point-location
func PointInPolygon(polygon []Point, q Point) (Location, error) {
	n := len(polygon)
	if n < 3 {
		return Outside, fmt.Errorf("%w: polygon has %d vertices, need 3", ErrTooFewPoints, n)
	}
	inside := false
	a := polygon[n-1]
	for _, b := range polygon {
		turn := Orientation(a, b, q)
		if turn == 0 && inBox(a, b, q) {
			return OnBoundary, nil
		}
		if (a.Y > q.Y) != (b.Y > q.Y) {
			// The edge meets the line y = q.Y at one point; it is on the ray when q lies on the
			// edge's −X side: left of an upward edge, right of a downward one. turn != 0 here,
			// since a collinear q inside the edge's Y range would be on the edge.
			if (b.Y > a.Y) == (turn > 0) {
				inside = !inside
			}
		}
		a = b
	}
	if inside {
		return Inside, nil
	}
	return Outside, nil
}

// inBox reports whether q lies in the closed bounding box of a and b; with Orientation(a, b, q) == 0
// it means q is on the closed segment ab (CLRS §33.1, ON-SEGMENT).
func inBox(a, b, q Point) bool {
	return min(a.X, b.X) <= q.X && q.X <= max(a.X, b.X) && min(a.Y, b.Y) <= q.Y && q.Y <= max(a.Y, b.Y)
}
