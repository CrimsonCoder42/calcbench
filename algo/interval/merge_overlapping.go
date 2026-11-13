package interval

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// ErrInverted: an interval's End is less than its Start.
var ErrInverted = errors.New("interval: End is less than Start")

// Interval is the closed range [Start, End] of an ordered type; valid when Start <= End.
// Its unit is the caller's (minutes, metres, index positions); both ends share it.
type Interval[T cmp.Ordered] struct {
	Start T
	End   T
}

// MergeOverlapping returns the union of intervals as disjoint intervals sorted by Start. Intervals that
// overlap or share an endpoint become one; [1,2] and [3,4] stay apart (they are closed, not adjacent
// integers). No intervals give an empty, non-nil slice. intervals is read, never modified.
// Errors: ErrInverted, with the index, when an interval has End < Start.
//
// ID:      ALG-INTVL-001
// Time:    O(n log n)   n = len(intervals)
// Space:   O(log n) extra (the sort's stack); O(n) for the result (one allocation)
// Source:  Klee, "Can the measure of ∪[aᵢ,bᵢ] be computed in less than O(n log n) steps?", American Mathematical Monthly, 1977, 84(4):284–285, doi:10.1080/00029890.1977.11994336
// Source:  Fredman, Weide, "On the complexity of computing the measure of ∪[aᵢ,bᵢ]", Communications of the ACM, 1978, 21(7):540–544, doi:10.1145/359545.359553 (Ω(n log n) lower bound)
// Tags:    interval, merge, union, overlap, sweep, sort, ranges
func MergeOverlapping[T cmp.Ordered](intervals []Interval[T]) ([]Interval[T], error) {
	for i, iv := range intervals {
		if cmp.Less(iv.End, iv.Start) {
			return nil, fmt.Errorf("%w: intervals[%d] = [%v, %v]", ErrInverted, i, iv.Start, iv.End)
		}
	}
	if len(intervals) == 0 {
		return []Interval[T]{}, nil
	}
	// The copy is the result buffer: sort it, then merge in place with a write index.
	merged := slices.Clone(intervals)
	slices.SortFunc(merged, func(x, y Interval[T]) int { return cmp.Compare(x.Start, y.Start) })
	last := 0
	for _, iv := range merged[1:] {
		if cmp.Compare(iv.Start, merged[last].End) <= 0 {
			if cmp.Less(merged[last].End, iv.End) {
				merged[last].End = iv.End
			}
		} else {
			last++
			merged[last] = iv
		}
	}
	return merged[:last+1], nil
}
