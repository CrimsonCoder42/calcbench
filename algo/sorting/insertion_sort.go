package sorting

import "cmp"

// InsertionSort returns values sorted ascending in cmp.Compare order, by straight insertion. It is the
// sort for small inputs (about 15 elements or fewer) and for nearly sorted ones, where it runs in time
// proportional to n plus the number of out-of-order pairs. For large unsorted input use slices.Sort or
// MergeSortStable. Stable. The result is written into dst (reused when cap(dst) >= len(values), else one
// new slice); dst may be values itself to sort in place. values is never modified unless it is dst.
//
// ID:      ALG-SORT-002
// Time:    O(n + I)   n = len(values), I = number of inversions (pairs i < j with values[i] > values[j]); O(n²) worst
// Space:   O(1) extra; O(n) for the result (none when dst has the capacity)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.1, Algorithm S (straight insertion)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §2.1 (INSERTION-SORT)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.3 (insertion sort for tiny arrays, cutoff 5 to 15)
// Tags:    sorting, insertion-sort, stable, small-n, adaptive, nearly-sorted
func InsertionSort[T cmp.Ordered](dst, values []T) []T {
	out := resize(dst, len(values))
	copy(out, values)
	insertionSortInPlace(out)
	return out
}

// insertionSortInPlace sorts s ascending in cmp.Compare order by straight insertion (shifting, not swapping).
func insertionSortInPlace[T cmp.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		v := s[i]
		j := i
		for j > 0 && cmp.Less(v, s[j-1]) {
			s[j] = s[j-1]
			j--
		}
		s[j] = v
	}
}
