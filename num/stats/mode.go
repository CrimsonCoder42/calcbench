package stats

import (
	"cmp"
	"fmt"
	"slices"
)

// Mode returns every value of values that occurs with the greatest frequency, ascending, and that
// frequency. The mode is not necessarily unique: all tied values are returned, and when every value is
// distinct every value is a mode with count 1. Values are compared with cmp.Compare, a total order in
// which NaN equals NaN and sorts first, and −0 equals +0 (the first of them in sorted order is returned).
// values is read, never modified: a copy is sorted.
// Errors: ErrEmpty when values is empty.
//
// ID:      NUM-STATS-024
// Time:    O(n log n) worst case   n = len(values); one sort of the copy, then two linear scans
// Space:   O(n) extra (the copy) plus O(k) for the result, k = number of modes; two allocations
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12)
// Tags:    mode, location, central-tendency, frequency, descriptive
// Units:   values of any ordered type; modes in that type; count dimensionless
// Example: 2, 3, 3, 5, 5, 5 → [5], 3; 1, 1, 2, 2, 3 → [1 2], 2 (hand-checked)
func Mode[T cmp.Ordered](values []T) (modes []T, count int, err error) {
	if len(values) == 0 {
		return nil, 0, fmt.Errorf("%w: values", ErrEmpty)
	}
	work := slices.Clone(values)
	slices.Sort(work)
	// Pass 1: the greatest run length and how many runs reach it.
	runs := 0
	for i := 0; i < len(work); {
		j := runEnd(work, i)
		switch l := j - i; {
		case l > count:
			count, runs = l, 1
		case l == count:
			runs++
		}
		i = j
	}
	// Pass 2: collect the runs of that length, already ascending.
	modes = make([]T, 0, runs)
	for i := 0; i < len(work); {
		j := runEnd(work, i)
		if j-i == count {
			modes = append(modes, work[i])
		}
		i = j
	}
	return modes, count, nil
}

// runEnd returns the index just past the run of elements equal (by cmp.Compare) to sorted[i].
func runEnd[T cmp.Ordered](sorted []T, i int) int {
	j := i + 1
	for j < len(sorted) && cmp.Compare(sorted[j], sorted[i]) == 0 {
		j++
	}
	return j
}
