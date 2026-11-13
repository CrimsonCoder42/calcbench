package sorting

import "cmp"

// DedupeSorted returns the distinct values of sorted, ascending: the first element of each run of
// elements equal under cmp.Compare (so every NaN is one value, and -0 and +0 are one value, the first
// met). sorted must be ascending; on unsorted input only adjacent repeats are removed, and nothing
// panics. The result is written into dst: when cap(dst) >= len(sorted) it is a prefix of dst; otherwise
// one new slice of exactly the distinct count is allocated. dst may be sorted itself (dedupe in place).
// sorted is never modified unless it is dst. Empty input returns dst[:0].
//
// ID:      ALG-SORT-007
// Time:    O(n)   n = len(sorted); two passes when dst must be allocated (count, then copy)
// Space:   O(1) extra; O(d) for the result, d = number of distinct values (none when dst has the capacity)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.5 "Duplicates" (sort, then one pass noting equal values that appear consecutively)
// Tags:    sorting, dedupe, unique, distinct, sorted
func DedupeSorted[T cmp.Ordered](dst, sorted []T) []T {
	n := len(sorted)
	if n == 0 {
		return dst[:0]
	}
	var out []T
	if cap(dst) >= n {
		out = dst[:n]
	} else {
		distinct := 1
		for i := 1; i < n; i++ {
			if cmp.Compare(sorted[i], sorted[i-1]) != 0 {
				distinct++
			}
		}
		out = make([]T, distinct)
	}
	// Compare with the last value kept, out[w-1], not sorted[i-1]: when out is sorted, the write index
	// w never passes the read index i, but sorted[i-1] may already have been overwritten.
	out[0] = sorted[0]
	w := 1
	for i := 1; i < n; i++ {
		if cmp.Compare(sorted[i], out[w-1]) != 0 {
			out[w] = sorted[i]
			w++
		}
	}
	return out[:w]
}
