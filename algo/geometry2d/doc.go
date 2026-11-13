// Package geometry2d holds computational geometry in the plane: algorithms over point sets, polygons,
// segments and axis-aligned rectangles. Measurement formulas (polygon area, distance, the area of one
// rectangle) live in num/geometry; this package does not import num.
//
// ID short form: GEO2D (IDs are ALG-GEO2D-NNN).
//
// Coordinates: Point has exact int64 coordinates in any one length unit the caller chooses; areas are in
// its square. Every predicate is exact on the whole int64 range: coordinate differences are carried in 65
// bits, cross products in 128 and sums of squares in 192 (math/bits), so no answer is ever wrong near a
// degeneracy and no input is rejected for size. Callers with real-valued data scale to a fixed grid
// (for example micrometres). Floating-point points with Shewchuk's adaptive predicates were the
// alternative (Shewchuk, Discrete & Computational Geometry 18:305–363, 1997, doi:10.1007/PL00009321);
// see doc/algorithms_core_2026-11-12/catalog/geometry2d.md for why they lost.
//
// Orientation is +1 counterclockwise, −1 clockwise, 0 collinear. Segments and rectangles are closed.
// ConvexHull appends to dst (doc/DECISIONS.md D-003). No function modifies its input.
//
// Functions:
//
//	ALG-GEO2D-001  Orientation              the turn of three points: counterclockwise, clockwise or collinear
//	ALG-GEO2D-002  ConvexHull               the strict convex hull, counterclockwise (Andrew's monotone chain)
//	ALG-GEO2D-003  ClosestPair              the indices of two points at the least distance, O(n log n)
//	ALG-GEO2D-004  PointInPolygon           Inside, Outside or OnBoundary of a polygon (crossing number)
//	ALG-GEO2D-005  SegmentsIntersect        whether two closed segments share a point
//	ALG-GEO2D-006  MaxPointsOnLine          the most points on one straight line
//	ALG-GEO2D-007  RectanglesUnionArea      the area covered by two axis-aligned rectangles
//	ALG-GEO2D-008  HasVerticalSymmetryAxis  whether a vertical line reflects the point set onto itself
package geometry2d
