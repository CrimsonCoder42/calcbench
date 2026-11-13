// Package interval works on closed ranges of ordered values.
//
// ID short form: INTVL (IDs are ALG-INTVL-NNN).
//
// An Interval is closed, [Start, End], with Start <= End (doc/ASSUMPTIONS.md A-005). Ordering is
// cmp.Compare's total order, in which a floating-point NaN sorts before every other value and -0
// equals +0 (doc/DECISIONS.md D-002). No function modifies its input.
// Closed means two intervals that share an endpoint overlap: they merge, intersect in a point, need two
// rooms, and conflict in scheduling. A caller with half-open bookings passes End minus one unit. The
// exception is Skyline, whose buildings raise the outline on Left <= x < Right, the usual reading of a
// key point. "Sorted and disjoint" means each interval starts strictly after the previous one ends,
// which is what MergeOverlapping returns. Skyline uses algo/heap's Heap.
//
// Functions:
//
//	ALG-INTVL-001  MergeOverlapping  the union of closed intervals as disjoint sorted intervals
//	ALG-INTVL-002  InsertInterval  a sorted, disjoint list with one interval merged in
//	ALG-INTVL-003  IntersectLists  the intersection of two sorted, disjoint lists
//	ALG-INTVL-004  MinMeetingRooms  the fewest rooms that hold every interval (the depth)
//	ALG-INTVL-005  MissingRanges  the integers of [lower, upper] absent from an ascending slice, as ranges
//	ALG-INTVL-006  SummaryRanges  ascending integers as ranges of consecutive values
//	ALG-INTVL-007  Covers  whether a union of intervals contains a target interval
//	ALG-INTVL-008  IntervalSchedulingMaxNonOverlap  a largest set of pairwise disjoint intervals (greedy by End)
//	ALG-INTVL-009  Skyline  the outline of rectangular buildings as key points
package interval
