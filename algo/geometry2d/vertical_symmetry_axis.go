package geometry2d

// HasVerticalSymmetryAxis reports whether some vertical line x = c reflects the set of points onto
// itself. Points are a set: duplicates are ignored. Such a line, when it exists, is x = (minX + maxX)/2
// (the reflection must swap the least and the greatest X), so the set is symmetric exactly when every
// point's mirror (minX + maxX − X, Y) is also in it. The empty set is symmetric. points is read, never
// modified. The mirror is computed without overflow.
//
// ID:      ALG-GEO2D-008
// Time:    O(n) expected   n = len(points); Ω(n log n) for a deterministic comparison-based method
// Space:   O(n) extra (one hash set)
// Source:  Atallah, "On symmetry detection", IEEE Transactions on Computers C-34(7):663–666, 1985, doi:10.1109/TC.1985.1676605 (axes of symmetry of a planar point set)
// Source:  the vertical-axis specialisation is proved in doc/algorithms_core_2026-11-12/catalog/geometry2d.md (ALG-GEO2D-008)
// Tags:    geometry, symmetry, reflection, mirror, line-reflection, hashing
func HasVerticalSymmetryAxis(points []Point) bool {
	if len(points) == 0 {
		return true
	}
	minX, maxX := points[0].X, points[0].X
	set := make(map[Point]struct{}, len(points))
	for _, p := range points {
		minX, maxX = min(minX, p.X), max(maxX, p.X)
		set[p] = struct{}{}
	}
	for p := range set {
		// minX + (maxX − X) lies in [minX, maxX]; the wrapped uint64 arithmetic gives it exactly.
		mirror := Point{int64(uint64(minX) + (uint64(maxX) - uint64(p.X))), p.Y}
		if _, ok := set[mirror]; !ok {
			return false
		}
	}
	return true
}
