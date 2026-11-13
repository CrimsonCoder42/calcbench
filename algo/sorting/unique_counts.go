package sorting

import "cmp"

// UniqueCounts returns the distinct values of values in ascending order and how many times each occurs
// (counts[j] for uniques[j]; the counts sum to len(values)). Values compare with cmp.Compare: all NaNs
// count as one value, sorted first, and -0 and +0 are one value, represented by whichever occurs first.
// Empty input gives empty, non-nil slices. values is read, never modified.
//
// ID:      ALG-SORT-013
// Time:    O(n log n)   n = len(values)
// Space:   O(n) extra (the sorting permutation and merge buffer); O(u) for the results, u = distinct values
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2 (sorting brings equal keys together, which finds the distinct keys and their multiplicities)
// Tags:    sorting, unique, distinct, counts, histogram, multiplicity
func UniqueCounts[T cmp.Ordered](values []T) (uniques []T, counts []int) {
	uniques, _, _, counts = uniqueSorted(values, false, false, true)
	return uniques, counts
}
