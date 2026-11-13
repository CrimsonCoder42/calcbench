package interval

import (
	"cmp"
	"fmt"
	"slices"
)

// Covers reports whether the union of intervals contains every point of the closed interval target.
// Intervals are closed and their values continuous: [1,2] and [3,4] leave (2,3) open, so they do not
// cover [1,4]; integer callers who want neighbours to join pass End+1 (doc/ASSUMPTIONS.md A-005).
// No intervals cover nothing. intervals is read, never modified.
// Errors: ErrInverted when target or an interval (with its index) has End < Start.
//
// ID:      ALG-INTVL-007
// Time:    O(n log n)   n = len(intervals)
// Space:   O(n) extra (a sorted copy, one allocation)
// Source:  Klee, "Can the measure of ∪[aᵢ,bᵢ] be computed in less than O(n log n) steps?", American Mathematical Monthly, 1977, 84(4):284–285, doi:10.1080/00029890.1977.11994336 (sort by left end, sweep the union)
// Source:  Fredman, Weide, "On the complexity of computing the measure of ∪[aᵢ,bᵢ]", Communications of the ACM, 1978, 21(7):540–544, doi:10.1145/359545.359553 (Ω(n log n) for the union's measure)
// Tags:    interval, cover, coverage, union, sweep, gaps
func Covers[T cmp.Ordered](intervals []Interval[T], target Interval[T]) (bool, error) {
	if cmp.Less(target.End, target.Start) {
		return false, fmt.Errorf("%w: target = [%v, %v]", ErrInverted, target.Start, target.End)
	}
	for i, iv := range intervals {
		if cmp.Less(iv.End, iv.Start) {
			return false, fmt.Errorf("%w: intervals[%d] = [%v, %v]", ErrInverted, i, iv.Start, iv.End)
		}
	}
	sorted := slices.Clone(intervals)
	slices.SortFunc(sorted, func(x, y Interval[T]) int { return cmp.Compare(x.Start, y.Start) })
	// reach: every point of [target.Start, reach] is covered, once started.
	started := false
	var reach T
	for _, iv := range sorted {
		if !started {
			if cmp.Less(iv.End, target.Start) {
				continue // wholly before the target
			}
			if cmp.Less(target.Start, iv.Start) {
				return false, nil // target.Start itself is uncovered
			}
			started, reach = true, iv.End
		} else {
			if cmp.Less(reach, iv.Start) {
				return false, nil // the points between reach and iv.Start are uncovered
			}
			reach = max(reach, iv.End)
		}
		if cmp.Compare(reach, target.End) >= 0 {
			return true, nil
		}
	}
	return false, nil
}
