package geometry2d

import "math/bits"

// MaxPointsOnLine returns the largest number of points of points that lie on one straight line. Points
// are a multiset: every copy of a point counts, and lies on every line through it. Empty input gives 0;
// one point gives 1. Exact for every int64 coordinate. points is read, never modified.
//
// ID:      ALG-GEO2D-006
// Time:    O(n²) expected   n = len(points); deciding three collinear points is 3SUM-hard, so o(n²) is not expected
// Space:   O(n) extra (one hash map, reused for every anchor)
// Source:  Gajentaan, Overmars, "On a class of O(n²) problems in computational geometry", Computational Geometry 5(3):165–185, 1995, doi:10.1016/0925-7721(95)00022-2 (the problem; its O(n²) class)
// Source:  per-anchor direction counting, proved in doc/algorithms_core_2026-11-12/catalog/geometry2d.md (ALG-GEO2D-006); directions reduced by Stein's binary gcd: Stein, "Computational problems associated with Racah algebra", Journal of Computational Physics 1(3):397–405, 1967, doi:10.1016/0021-9991(67)90047-2
// Tags:    geometry, collinear, lines, gcd, hashing
func MaxPointsOnLine(points []Point) int {
	n := len(points)
	if n <= 2 {
		return n
	}
	counts := make(map[direction]int, n)
	best := 1
	// Every line through two or more points has a first point (least index); anchor there and count
	// the later points by their direction from it.
	for i := 0; i < n && best < n-i; i++ {
		clear(counts)
		copies, largest := 0, 0
		for j := i + 1; j < n; j++ {
			d, same := directionBetween(points[i], points[j])
			if same {
				copies++
				continue
			}
			counts[d]++
			largest = max(largest, counts[d])
		}
		best = max(best, 1+copies+largest)
	}
	return best
}

// direction is a line's direction up to sign: |dx| and |dy| divided by their gcd, and whether dx and dy
// have opposite signs (false when either is 0).
type direction struct {
	dx, dy   uint64
	opposite bool
}

// directionBetween returns the reduced direction of the line through p and q, or same = true when the
// points are equal.
func directionBetween(p, q Point) (d direction, same bool) {
	dx, dxNeg := coordDiff(q.X, p.X)
	dy, dyNeg := coordDiff(q.Y, p.Y)
	if dx == 0 && dy == 0 {
		return direction{}, true
	}
	g := binaryGCD(dx, dy)
	return direction{dx / g, dy / g, dx != 0 && dy != 0 && dxNeg != dyNeg}, false
}

// binaryGCD returns gcd(u, v) by Stein's binary algorithm (Stein 1967); gcd(0, v) = v. Kept private:
// the exported gcd belongs to num/integer (lead decision L4).
func binaryGCD(u, v uint64) uint64 {
	if u == 0 {
		return v
	}
	if v == 0 {
		return u
	}
	shift := bits.TrailingZeros64(u | v)
	u >>= bits.TrailingZeros64(u)
	for {
		v >>= bits.TrailingZeros64(v)
		if u > v {
			u, v = v, u
		}
		v -= u
		if v == 0 {
			return u << shift
		}
	}
}
