package sorting

import (
	"errors"
	"fmt"
)

// ErrNilCompare: the compare func is nil.
var ErrNilCompare = errors.New("sorting: compare func is nil")

// smallRun is the length at or below which a range is sorted by insertion sort. Sedgewick, Wayne,
// Algorithms, 4th ed., §2.3: "any value between 5 and 15 is likely to work well"; not measured here.
const smallRun = 12

// MergeSortStable returns values sorted ascending by compare, keeping the input order of elements that
// compare equal (stable). The result is written into dst (reused when cap(dst) >= len(values), else one
// new slice); dst may be values itself to sort in place. values is never modified unless it is dst.
// compare follows the slices.SortFunc contract: negative when a < b, zero when equal, positive when a > b.
// Errors: ErrNilCompare when compare is nil (dst is returned unwritten).
//
// ID:      ALG-SORT-001
// Time:    O(n log n)   n = len(values); worst case, compares and moves
// Space:   O(n) extra (one allocation, the merge buffer; none when n <= 12); O(n) for the result
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.4 (sorting by merging)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.2 (bottom-up mergesort; between 1/2 N lg N and N lg N compares)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §2.3.1 (MERGE)
// Tags:    sorting, merge-sort, stable, comparison-sort, bottom-up
func MergeSortStable[T any](dst, values []T, compare func(a, b T) int) ([]T, error) {
	if compare == nil {
		return dst, fmt.Errorf("%w: compare", ErrNilCompare)
	}
	out := resize(dst, len(values))
	copy(out, values) // a no-op when out and values share their first element.
	var buf []T
	if len(out) > smallRun {
		buf = make([]T, len(out))
	}
	mergeSortFunc(out, buf, compare)
	return out, nil
}

// resize returns dst[:n] when dst has the capacity, otherwise a new slice of length n.
func resize[T any](dst []T, n int) []T {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]T, n)
}

// mergeSortFunc sorts a stably by compare: insertion sort on runs of smallRun, then bottom-up merging
// of run pairs, alternating between a and buf. buf must have len(buf) >= len(a) when len(a) > smallRun.
// Iterative: no recursion.
func mergeSortFunc[E any](a, buf []E, compare func(x, y E) int) {
	n := len(a)
	for lo := 0; lo < n; lo += smallRun {
		insertionSortFunc(a[lo:min(lo+smallRun, n)], compare)
	}
	if n <= smallRun {
		return
	}
	src, to := a, buf[:n]
	for width := smallRun; width < n; width *= 2 {
		for lo := 0; lo < n; lo += 2 * width {
			mid := min(lo+width, n)
			hi := min(mid+width, n)
			mergeRuns(to[lo:hi], src[lo:mid], src[mid:hi], compare)
		}
		src, to = to, src
	}
	if &src[0] != &a[0] {
		copy(a, src)
	}
}

// mergeRuns merges the ascending runs left and right into out (len(out) = len(left)+len(right)).
// On a tie the left element goes first, which is what makes the sort stable.
func mergeRuns[E any](out, left, right []E, compare func(x, y E) int) {
	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		if compare(right[j], left[i]) < 0 {
			out[k] = right[j]
			j++
		} else {
			out[k] = left[i]
			i++
		}
		k++
	}
	k += copy(out[k:], left[i:])
	copy(out[k:], right[j:])
}

// insertionSortFunc sorts s stably by compare (straight insertion, Knuth §5.2.1 Algorithm S).
func insertionSortFunc[E any](s []E, compare func(x, y E) int) {
	for i := 1; i < len(s); i++ {
		v := s[i]
		j := i
		for j > 0 && compare(v, s[j-1]) < 0 {
			s[j] = s[j-1]
			j--
		}
		s[j] = v
	}
}
