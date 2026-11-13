package interval

import (
	"cmp"
	"fmt"
	"slices"
)

// IntervalSchedulingMaxNonOverlap returns the indices of a largest set of pairwise disjoint intervals,
// in order of End (ties by index): the greedy choice of the interval that finishes first. Intervals are
// closed, so two that share an endpoint conflict; a caller with half-open bookings passes End minus one
// unit (doc/ASSUMPTIONS.md A-005). No intervals give an empty, non-nil slice. intervals is read, never
// modified.
// Errors: ErrInverted, with the index, when an interval has End < Start.
//
// ID:      ALG-INTVL-008
// Time:    O(n log n)   n = len(intervals)
// Space:   O(log n) extra (the sort's stack); O(n) for the result (one allocation of n indices, trimmed)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §16.1 (activity selection; Theorem 16.1, GREEDY-ACTIVITY-SELECTOR)
// Source:  Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §4.1 (interval scheduling: earliest finish time, "the greedy algorithm stays ahead")
// Tags:    interval, scheduling, greedy, activity-selection, non-overlapping, maximum-set
func IntervalSchedulingMaxNonOverlap[T cmp.Ordered](intervals []Interval[T]) ([]int, error) {
	for i, iv := range intervals {
		if cmp.Less(iv.End, iv.Start) {
			return nil, fmt.Errorf("%w: intervals[%d] = [%v, %v]", ErrInverted, i, iv.Start, iv.End)
		}
	}
	order := make([]int, len(intervals))
	for i := range order {
		order[i] = i
	}
	slices.SortFunc(order, func(p, q int) int {
		return cmp.Or(cmp.Compare(intervals[p].End, intervals[q].End), cmp.Compare(p, q))
	})
	// The kept indices overwrite the front of order: the write index never passes the read index.
	kept := 0
	for _, i := range order {
		if kept == 0 || cmp.Less(intervals[order[kept-1]].End, intervals[i].Start) {
			order[kept] = i
			kept++
		}
	}
	return order[:kept:kept], nil
}
