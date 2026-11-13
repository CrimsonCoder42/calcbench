package interval

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/heap"
)

// ErrNegativeHeight: a building's height is below the ground (the zero value of its type) or is NaN.
var ErrNegativeHeight = errors.New("interval: height is below zero or NaN")

// Building is a rectangle standing on the ground (height zero) from Left to Right. All three share the
// caller's unit for each axis (metres, pixels); Left <= Right and Height >= 0.
type Building[T cmp.Ordered] struct {
	Left, Right, Height T
}

// KeyPoint is a corner of a skyline: from X on, the outline's height is Height.
type KeyPoint[T cmp.Ordered] struct {
	X, Height T
}

// Skyline returns the outline of buildings as key points sorted by X: each says the outline's height
// from X until the next key point, and the last one has Height zero. A building raises the outline on
// Left <= x < Right, so buildings that touch at an edge join without a dip. Consecutive key points
// differ in height. A building of zero width or zero height adds nothing; none give an empty, non-nil
// slice. buildings is read, never modified.
// Errors: ErrInverted when a building has Right < Left; ErrNegativeHeight when its Height is below zero
// (NaN included, D-002); both with the index.
//
// ID:      ALG-INTVL-009
// Time:    O(n log n)   n = len(buildings)
// Space:   O(n) extra (the edge list, the order and the heap); O(n) for the result (≤ 2n key points)
// Source:  Manber, Introduction to Algorithms: A Creative Approach, Addison-Wesley, 1989, §5.6 "The Skyline Problem" (the problem, and O(n log n) by divide and conquer)
// Source:  de Berg, Cheong, van Kreveld, Overmars, Computational Geometry: Algorithms and Applications, 3rd ed., Springer, 2008, §2.1 (plane sweep: an event queue in x and a status structure; here a max-heap with lazy deletion)
// Tags:    interval, skyline, sweep-line, heap, outline, silhouette, geometry
func Skyline[T cmp.Ordered](buildings []Building[T]) ([]KeyPoint[T], error) {
	var ground T
	live := make([]int, 0, len(buildings)) // buildings with width and height, in Left order once sorted
	for i, bd := range buildings {
		if cmp.Less(bd.Right, bd.Left) {
			return nil, fmt.Errorf("%w: buildings[%d] has Left = %v > Right = %v", ErrInverted, i, bd.Left, bd.Right)
		}
		if cmp.Less(bd.Height, ground) {
			return nil, fmt.Errorf("%w: buildings[%d].Height = %v", ErrNegativeHeight, i, bd.Height)
		}
		if cmp.Less(bd.Left, bd.Right) && cmp.Less(ground, bd.Height) {
			live = append(live, i)
		}
	}
	out := []KeyPoint[T]{}
	if len(live) == 0 {
		return out, nil
	}
	slices.SortFunc(live, func(p, q int) int { return cmp.Compare(buildings[p].Left, buildings[q].Left) })
	xs := make([]T, 0, 2*len(live))
	for _, i := range live {
		xs = append(xs, buildings[i].Left, buildings[i].Right)
	}
	slices.SortFunc(xs, cmp.Compare[T])
	xs = slices.CompactFunc(xs, func(x, y T) bool { return cmp.Compare(x, y) == 0 })

	// The heap holds the buildings that have started, tallest first. One that has ended stays until it
	// reaches the top, then is dropped (lazy deletion): only the top's height is ever read.
	tallest, err := heap.NewHeap(func(p, q int) bool { return cmp.Less(buildings[q].Height, buildings[p].Height) }, len(live))
	if err != nil {
		return nil, err // unreachable: the less function is not nil
	}
	next, height := 0, ground
	for _, x := range xs {
		for next < len(live) && cmp.Compare(buildings[live[next]].Left, x) <= 0 {
			tallest.Push(live[next])
			next++
		}
		for {
			top, ok := tallest.Peek()
			if !ok || cmp.Less(x, buildings[top].Right) {
				break
			}
			tallest.Pop()
		}
		h := ground
		if top, ok := tallest.Peek(); ok {
			h = buildings[top].Height
		}
		if cmp.Compare(h, height) != 0 {
			out = append(out, KeyPoint[T]{x, h})
			height = h
		}
	}
	return out, nil
}
