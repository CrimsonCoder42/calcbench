package geometry2d

import (
	"errors"
	"fmt"
	"slices"
)

// ErrTooFewPoints: the input has fewer points than the function needs.
var ErrTooFewPoints = errors.New("geometry2d: too few points")

// ClosestPair returns the indices i < j of two points of points at the least Euclidean distance. Equal
// points are a pair at distance 0. When several pairs tie, one of them is returned (deterministically for
// a given input). Distances are compared exactly. points is read, never modified.
// Errors: ErrTooFewPoints when len(points) < 2.
//
// ID:      ALG-GEO2D-003
// Time:    O(n log n)   n = len(points); Θ(n log n) in the algebraic decision-tree model
// Space:   O(n) extra (two index slices); recursion depth ⌈log2 n⌉
// Source:  Shamos, Hoey, "Closest-point problems", 16th Annual Symposium on Foundations of Computer Science (SFCS 1975), IEEE, pp. 151–162, doi:10.1109/SFCS.1975.8
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.4 (divide and conquer, presorting, the strip)
// Tags:    geometry, closest-pair, nearest, divide-and-conquer
func ClosestPair(points []Point) (i, j int, err error) {
	n := len(points)
	if n < 2 {
		return 0, 0, fmt.Errorf("%w: points has %d, need 2", ErrTooFewPoints, n)
	}
	byX := make([]int, n)
	for k := range byX {
		byX[k] = k
	}
	slices.SortFunc(byX, func(a, b int) int {
		if c := comparePoints(points[a], points[b]); c != 0 {
			return c
		}
		return a - b
	})
	c := closestSearch{points: points}
	c.solve(byX, make([]int, n))
	return min(c.i, c.j), max(c.i, c.j), nil
}

// closestSearch holds the best pair found so far and its squared distance.
type closestSearch struct {
	points []Point
	best   uint192
	i, j   int
	found  bool
}

func (c *closestSearch) consider(a, b int) {
	d := squaredDistance(c.points[a], c.points[b])
	if !c.found || d.less(c.best) {
		c.best, c.i, c.j, c.found = d, a, b, true
	}
}

// squaredGapLess reports whether (u − v)² < best.
func (c *closestSearch) squaredGapLess(u, v int64) bool {
	m, _ := coordDiff(u, v)
	return widen128(mul64(m, m)).less(c.best)
}

// solve finds the closest pair among the points whose indices ix holds, sorted by X, and leaves ix sorted
// by Y (merge sort on the way back, CLRS §33.4). aux has the length of ix. len(ix) >= 2.
func (c *closestSearch) solve(ix, aux []int) {
	n := len(ix)
	if n <= 3 {
		for a := 0; a < n; a++ {
			for b := a + 1; b < n; b++ {
				c.consider(ix[a], ix[b])
			}
		}
		// Insertion sort by Y.
		for a := 1; a < n; a++ {
			for b := a; b > 0 && c.points[ix[b]].Y < c.points[ix[b-1]].Y; b-- {
				ix[b], ix[b-1] = ix[b-1], ix[b]
			}
		}
		return
	}
	mid := n / 2
	midX := c.points[ix[mid]].X
	c.solve(ix[:mid], aux[:mid])
	c.solve(ix[mid:], aux[mid:])
	// Merge the two halves by Y into aux, then back.
	left, right, k := 0, mid, 0
	for left < mid && right < n {
		if c.points[ix[right]].Y < c.points[ix[left]].Y {
			aux[k] = ix[right]
			right++
		} else {
			aux[k] = ix[left]
			left++
		}
		k++
	}
	k += copy(aux[k:], ix[left:mid])
	copy(aux[k:], ix[right:n])
	copy(ix, aux)
	if c.best == (uint192{}) {
		return // distance 0 cannot be beaten
	}
	// The strip: points within the best distance of the dividing line, in Y order. Each is compared with
	// the earlier strip points less than the best distance below it (at most 7, CLRS §33.4).
	strip := aux[:0]
	for _, q := range ix {
		p := c.points[q]
		if !c.squaredGapLess(p.X, midX) {
			continue
		}
		for t := len(strip) - 1; t >= 0; t-- {
			if !c.squaredGapLess(p.Y, c.points[strip[t]].Y) {
				break
			}
			c.consider(strip[t], q)
		}
		strip = append(strip, q)
	}
}
