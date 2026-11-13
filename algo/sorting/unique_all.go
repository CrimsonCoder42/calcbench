package sorting

import "cmp"

// UniqueAll returns the distinct values of values in ascending order, and for each: firstIndex, the index
// of its first occurrence; counts, how many times it occurs; and inverse, one entry per input, the
// position of that input's value in uniques, so values[i] == uniques[inverse[i]] (up to -0/+0 and NaN).
// Values compare with cmp.Compare: all NaNs count as one value (sorted first), and -0 and +0 are one
// value, represented by whichever occurs first. Empty input gives empty, non-nil slices. values is read,
// never modified.
//
// ID:      ALG-SORT-012
// Time:    O(n log n)   n = len(values)
// Space:   O(n) extra (the sorting permutation and merge buffer); O(n) for the four results
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2 (sorting brings equal keys together, which finds the distinct keys and their multiplicities) and §5.2.4 (stable merge sort)
// Tags:    sorting, unique, distinct, counts, inverse, first-occurrence, label-encoding
func UniqueAll[T cmp.Ordered](values []T) (uniques []T, firstIndex, inverse, counts []int) {
	return uniqueSorted(values, true, true, true)
}

// uniqueSorted does the work of UniqueAll, UniqueCounts and UniqueInverse: a stable argsort, then one walk
// over equal runs. Stability puts the smallest index of each run first. Results not wanted are nil.
func uniqueSorted[T cmp.Ordered](values []T, wantFirst, wantInverse, wantCounts bool) (uniques []T, firstIndex, inverse, counts []int) {
	n := len(values)
	idx, _ := SortIndices(nil, values, cmp.Compare[T]) // compare is non-nil: no error.
	uniques = []T{}
	if wantFirst {
		firstIndex = []int{}
	}
	if wantInverse {
		inverse = make([]int, n)
	}
	if wantCounts {
		counts = []int{}
	}
	for r := 0; r < n; {
		start := r
		v := values[idx[start]]
		for r < n && cmp.Compare(values[idx[r]], v) == 0 {
			if wantInverse {
				inverse[idx[r]] = len(uniques)
			}
			r++
		}
		uniques = append(uniques, v)
		if wantFirst {
			firstIndex = append(firstIndex, idx[start])
		}
		if wantCounts {
			counts = append(counts, r-start)
		}
	}
	return uniques, firstIndex, inverse, counts
}
