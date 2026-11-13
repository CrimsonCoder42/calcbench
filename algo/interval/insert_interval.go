package interval

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrNotSortedDisjoint: a list that must be sorted and disjoint has an interval that starts at or
// before the end of the one before it (closed intervals that touch are not disjoint).
var ErrNotSortedDisjoint = errors.New("interval: list is not sorted and disjoint")

// checkSortedDisjoint reports the first interval of list that is inverted or does not start strictly
// after the previous one ends. name is the parameter name used in the error.
func checkSortedDisjoint[T cmp.Ordered](name string, list []Interval[T]) error {
	for i, iv := range list {
		if cmp.Less(iv.End, iv.Start) {
			return fmt.Errorf("%w: %s[%d] = [%v, %v]", ErrInverted, name, i, iv.Start, iv.End)
		}
		if i > 0 && cmp.Compare(iv.Start, list[i-1].End) <= 0 {
			return fmt.Errorf("%w: %s[%d] = [%v, %v] starts at or before %s[%d].End = %v",
				ErrNotSortedDisjoint, name, i, iv.Start, iv.End, name, i-1, list[i-1].End)
		}
	}
	return nil
}

// InsertInterval returns sorted, the sorted and disjoint list of closed intervals that MergeOverlapping
// produces, with added put in place: every interval of sorted that overlaps added or shares an endpoint
// with it is merged into one. The result is again sorted and disjoint. sorted is read, never modified.
// Errors: ErrInverted when added or an interval of sorted has End < Start; ErrNotSortedDisjoint when
// sorted is not sorted and disjoint (checked in the same pass).
//
// ID:      ALG-INTVL-002
// Time:    O(n)   n = len(sorted)
// Space:   O(1) extra; O(n) for the result (one allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4 (one-pass merge of sorted sequences; here a list and one element)
// Source:  Klee, "Can the measure of ∪[aᵢ,bᵢ] be computed in less than O(n log n) steps?", American Mathematical Monthly, 1977, 84(4):284–285, doi:10.1080/00029890.1977.11994336 (union of intervals by a sweep in Start order)
// Tags:    interval, insert, merge, union, sorted, ranges
func InsertInterval[T cmp.Ordered](sorted []Interval[T], added Interval[T]) ([]Interval[T], error) {
	if cmp.Less(added.End, added.Start) {
		return nil, fmt.Errorf("%w: added = [%v, %v]", ErrInverted, added.Start, added.End)
	}
	if err := checkSortedDisjoint("sorted", sorted); err != nil {
		return nil, err
	}
	out := make([]Interval[T], 0, len(sorted)+1)
	i := 0
	for i < len(sorted) && cmp.Less(sorted[i].End, added.Start) { // wholly before added
		out = append(out, sorted[i])
		i++
	}
	merged := added
	for i < len(sorted) && cmp.Compare(sorted[i].Start, added.End) <= 0 { // overlaps or touches added
		merged.Start = min(merged.Start, sorted[i].Start)
		merged.End = max(merged.End, sorted[i].End)
		i++
	}
	out = append(out, merged)
	return append(out, sorted[i:]...), nil
}
