package interval

import (
	"cmp"
	"fmt"
	"slices"
)

// MinMeetingRooms returns the fewest rooms that hold every interval with no two in one room at the same
// time: the depth of the set, the largest number of intervals that contain a common point. Intervals are
// closed, so one that ends at 3 and one that starts at 3 need two rooms; a caller with half-open
// bookings passes End minus one unit (doc/ASSUMPTIONS.md A-005). No intervals need 0 rooms.
// intervals is read, never modified.
// Errors: ErrInverted, with the index, when an interval has End < Start.
//
// ID:      ALG-INTVL-004
// Time:    O(n log n)   n = len(intervals)
// Space:   O(n) extra (sorted starts and ends, one allocation); O(1) for the result
// Source:  Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §4.1, "A Related Problem: Scheduling All Intervals" (interval partitioning; rooms needed = depth)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §16.1, Exercise 16.1-4 (lecture halls; interval-graph coloring)
// Tags:    interval, sweep-line, partitioning, depth, max-overlap, scheduling, rooms, resources
func MinMeetingRooms[T cmp.Ordered](intervals []Interval[T]) (int, error) {
	for i, iv := range intervals {
		if cmp.Less(iv.End, iv.Start) {
			return 0, fmt.Errorf("%w: intervals[%d] = [%v, %v]", ErrInverted, i, iv.Start, iv.End)
		}
	}
	n := len(intervals)
	if n == 0 {
		return 0, nil
	}
	buf := make([]T, 2*n)
	starts, ends := buf[:n], buf[n:]
	for i, iv := range intervals {
		starts[i], ends[i] = iv.Start, iv.End
	}
	slices.SortFunc(starts, cmp.Compare[T])
	slices.SortFunc(ends, cmp.Compare[T])
	// Sweep the starts in order. A room is free for a start s when some earlier interval ended strictly
	// before s (closed intervals); otherwise a new room opens. Rooms opened = the maximum depth.
	rooms, ended := 0, 0
	for _, s := range starts {
		if cmp.Less(ends[ended], s) {
			ended++
		} else {
			rooms++
		}
	}
	return rooms, nil
}
