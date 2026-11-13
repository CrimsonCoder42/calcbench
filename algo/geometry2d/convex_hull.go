package geometry2d

import (
	"cmp"
	"slices"
)

// ConvexHull appends to dst the vertices of the convex hull of points, counterclockwise, starting at the
// point with the least X (the least Y among those), and returns the extended slice. Only strict convex
// vertices are returned: duplicates and points inside a hull edge are dropped. When every point is
// collinear the two extreme points are returned; one distinct point gives that point; no points give dst
// unchanged. points is read, never modified; dst == nil allocates once.
//
// ID:      ALG-GEO2D-002
// Time:    O(n log n)   n = len(points); Ω(n log n) in the algebraic decision-tree model
// Space:   O(n) extra (a sorted copy); O(h) appended to dst, h = number of hull vertices
// Source:  Andrew, "Another efficient algorithm for convex hulls in two dimensions", Information Processing Letters 9(5):216–219, 1979, doi:10.1016/0020-0190(79)90072-3 (monotone chain)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.3 (convex hull; lower bound)
// Tags:    geometry, convex-hull, monotone-chain, sort
func ConvexHull(dst, points []Point) []Point {
	sorted := slices.Clone(points)
	slices.SortFunc(sorted, comparePoints)
	sorted = slices.Compact(sorted)
	n := len(sorted)
	if n <= 2 {
		return append(dst, sorted...)
	}
	dst = slices.Grow(dst, n+1)
	base := len(dst)
	// Lower chain, left to right: pop while the last two and p do not make a strict left turn.
	for _, p := range sorted {
		for len(dst)-base >= 2 && Orientation(dst[len(dst)-2], dst[len(dst)-1], p) <= 0 {
			dst = dst[:len(dst)-1]
		}
		dst = append(dst, p)
	}
	// Upper chain, right to left, never popping into the lower chain.
	floor := len(dst) + 1
	for i := n - 2; i >= 0; i-- {
		p := sorted[i]
		for len(dst) >= floor && Orientation(dst[len(dst)-2], dst[len(dst)-1], p) <= 0 {
			dst = dst[:len(dst)-1]
		}
		dst = append(dst, p)
	}
	// The last point appended is the first point again.
	return dst[:len(dst)-1]
}

// comparePoints orders points by X, then by Y.
func comparePoints(p, q Point) int {
	if c := cmp.Compare(p.X, q.X); c != 0 {
		return c
	}
	return cmp.Compare(p.Y, q.Y)
}
