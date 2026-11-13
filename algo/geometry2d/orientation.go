package geometry2d

// Point is a point of the plane with exact integer coordinates, in any one length unit the caller
// chooses. Every function of this package is exact on the whole int64 range.
type Point struct {
	X, Y int64
}

// Orientation returns the turn made by a, b, c: +1 when they turn counterclockwise (c lies to the left
// of the directed line a→b), −1 when clockwise, 0 when they are collinear (including when two of them
// are equal). It is the sign of the cross product (b − a) × (c − a), computed exactly in 128 bits.
//
// ID:      ALG-GEO2D-001
// Time:    O(1)
// Space:   O(1) extra, no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.1 (cross products, DIRECTION)
// Source:  O'Rourke, Computational Geometry in C, 2nd ed., Cambridge University Press, 1998, §1.3 (Area2, Left)
// Tags:    geometry, orientation, ccw, cross-product, predicate, exact
func Orientation(a, b, c Point) int {
	bx, bxNeg := coordDiff(b.X, a.X)
	by, byNeg := coordDiff(b.Y, a.Y)
	cx, cxNeg := coordDiff(c.X, a.X)
	cy, cyNeg := coordDiff(c.Y, a.Y)
	// (b − a) × (c − a) = bx·cy − by·cx.
	return compareSigned(signedProduct(bx, bxNeg, cy, cyNeg), signedProduct(by, byNeg, cx, cxNeg))
}
