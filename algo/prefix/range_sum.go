package prefix

import "fmt"

// RangeSum returns values[lo] + … + values[hi-1] from prefix, the n+1 prefix sums of values made by
// PrefixSums: prefix[hi] - prefix[lo]. The range is half-open; lo == hi gives 0.
// prefix is read, never modified.
// Errors: ErrRange unless 0 <= lo <= hi < len(prefix); ErrOverflow when the range's sum does not fit T
// (each prefix fits, but their difference can be up to twice as large).
//
// ID:      ALG-PREFIX-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.1 "Static array queries", "Sum queries", p. 84 (sum(a, b) = sum(0, b) − sum(0, a−1)), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8 "Algorithm Design Techniques", §8.2 (sum of x[i..j] = cumarr[j] − cumarr[i−1])
// Tags:    prefix-sum, range-sum, range-query, integer, overflow-checked
func RangeSum[T integer](prefix []T, lo, hi int) (T, error) {
	if lo < 0 || hi < lo || hi >= len(prefix) {
		return 0, fmt.Errorf("%w: [lo, hi) = [%d, %d) with len(prefix) = %d", ErrRange, lo, hi, len(prefix))
	}
	d, ok := checkedSub(prefix[hi], prefix[lo])
	if !ok {
		return 0, fmt.Errorf("%w: sum of [%d, %d)", ErrOverflow, lo, hi)
	}
	return d, nil
}
