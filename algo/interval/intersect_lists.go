package interval

import (
	"cmp"
)

// IntersectLists returns the intersection of two sorted, disjoint lists of closed intervals, as a
// sorted, disjoint list: every non-empty [max(starts), min(ends)] of a pair, one from each list. Closed
// intervals that share only an endpoint meet in a single point: [1,3] ∩ [3,5] = [3,3]. No overlap gives
// an empty, non-nil slice. a and b are read, never modified.
// Errors: ErrInverted or ErrNotSortedDisjoint (with the list's name and index) when a or b is not a
// sorted, disjoint list of valid intervals.
//
// ID:      ALG-INTVL-003
// Time:    O(n + m)   n = len(a), m = len(b)
// Space:   O(1) extra; O(n + m) for the result (one allocation, sized n + m)
// Source:  Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4, Algorithm M (two-way merge: advance the list whose current item ends first)
// Tags:    interval, intersection, two-pointer, merge, sorted, ranges
func IntersectLists[T cmp.Ordered](a, b []Interval[T]) ([]Interval[T], error) {
	if err := checkSortedDisjoint("a", a); err != nil {
		return nil, err
	}
	if err := checkSortedDisjoint("b", b); err != nil {
		return nil, err
	}
	// Each step advances one list, and each pair contributes at most one piece, so n + m bounds the result.
	out := make([]Interval[T], 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		lo := max(a[i].Start, b[j].Start)
		hi := min(a[i].End, b[j].End)
		if cmp.Compare(lo, hi) <= 0 {
			out = append(out, Interval[T]{lo, hi})
		}
		if cmp.Less(a[i].End, b[j].End) {
			i++
		} else {
			j++
		}
	}
	return out, nil
}
