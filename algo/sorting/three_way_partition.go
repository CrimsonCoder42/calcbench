package sorting

import "cmp"

// ThreeWayPartition returns values rearranged into three blocks around pivot, in cmp.Compare order:
// partitioned[:less] < pivot, partitioned[less:notGreater] == pivot, partitioned[notGreater:] > pivot
// (Dijkstra's Dutch national flag). The order inside each block is unspecified (not stable).
// The result is written into dst (reused when cap(dst) >= len(values), else one new slice); dst may be
// values itself to partition in place. values is never modified unless it is dst.
//
// ID:      ALG-SORT-010
// Time:    O(n)   n = len(values); one pass, exactly n comparisons and at most n swaps
// Space:   O(1) extra; O(n) for the result (none when dst has the capacity)
// Source:  Dijkstra, A Discipline of Programming, Prentice-Hall, 1976, ch. 14 (the problem of the Dutch national flag)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.3 (3-way partitioning: a[lo..lt-1] < v, a[lt..i-1] = v, a[gt+1..hi] > v)
// Tags:    sorting, partition, three-way, dutch-national-flag, quicksort
func ThreeWayPartition[T cmp.Ordered](dst, values []T, pivot T) (partitioned []T, less, notGreater int) {
	out := resize(dst, len(values))
	copy(out, values)
	// Invariant: out[:lt] < pivot, out[lt:i] == pivot, out[i:gt] unread, out[gt:] > pivot.
	lt, i, gt := 0, 0, len(out)
	for i < gt {
		switch c := cmp.Compare(out[i], pivot); {
		case c < 0:
			out[lt], out[i] = out[i], out[lt]
			lt++
			i++
		case c > 0:
			gt--
			out[i], out[gt] = out[gt], out[i]
		default:
			i++
		}
	}
	return out, lt, gt
}
