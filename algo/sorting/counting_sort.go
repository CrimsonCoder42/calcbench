package sorting

import (
	"errors"
	"fmt"
)

// ErrKeyRange: maxKey or a value is outside the range counting sort accepts.
var ErrKeyRange = errors.New("sorting: key out of range")

// maxCountingKey is the largest maxKey CountingSort accepts: a count array of 1<<24 + 1 ints (128 MiB on
// a 64-bit platform). It stops one wrong argument from requesting gigabytes; wider keys belong to
// RadixSortUint64. A policy of this package (doc/algorithms_core_2026-11-12/catalog/sorting.md).
const maxCountingKey = 1 << 24

// CountingSort returns values sorted ascending, for values that are all in [0, maxKey]. It counts how
// often each key occurs and writes each key that many times, so it never compares two values.
// The result is written into dst (reused when cap(dst) >= len(values), else one new slice); dst may be
// values itself. values is never modified unless it is dst.
// Errors: ErrKeyRange when maxKey < 0, maxKey > 1<<24, or a value is outside [0, maxKey]; every value is
// checked before dst is written, and on error dst is returned unwritten.
//
// ID:      ALG-SORT-003
// Time:    O(n + k)   n = len(values), k = maxKey + 1
// Space:   O(k) extra (one allocation, the counts); O(n) for the result (none when dst has the capacity)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.2 (COUNTING-SORT)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2, Algorithm D (distribution counting)
// Tags:    sorting, counting-sort, integer-sort, linear-time, non-comparison
func CountingSort(dst, values []int, maxKey int) ([]int, error) {
	if maxKey < 0 || maxKey > maxCountingKey {
		return dst, fmt.Errorf("%w: maxKey = %d, want 0..%d", ErrKeyRange, maxKey, maxCountingKey)
	}
	counts := make([]int, maxKey+1)
	for i, v := range values {
		if v < 0 || v > maxKey {
			return dst, fmt.Errorf("%w: values[%d] = %d, want 0..%d", ErrKeyRange, i, v, maxKey)
		}
		counts[v]++
	}
	out := resize(dst, len(values))
	i := 0
	for key, c := range counts {
		for range c {
			out[i] = key
			i++
		}
	}
	return out, nil
}
