package prefix

import "fmt"

// MaxSubarray returns the largest sum of a non-empty contiguous subarray of values, and that subarray
// as the half-open range [start, end). Among subarrays with the largest sum it returns the one that
// ends first, and among those the shortest. When every value is negative the result is the single
// largest value.
// values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrOverflow when the largest sum does not fit T (a running sum
// overflows only when it is positive and still growing, so the true maximum is then out of range too).
//
// ID:      ALG-PREFIX-009
// Time:    O(n)   n = len(values); one pass, at most n−1 checked additions
// Space:   O(1), no allocation
// Source:  Bentley, "Programming Pearls: Algorithm Design Techniques", Communications of the ACM, 1984, 27(9):865–873, doi:10.1145/358234.381162 (the scanning algorithm, attributed to Kadane)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8 "Algorithm Design Techniques", §8.4 "A Scanning Algorithm" (Column 7 in the 1st ed., 1986)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §2.4 "Maximum subarray sum", Algorithm 3, pp. 21–23, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    kadane, maximum-subarray, scanning, dynamic-programming, integer, overflow-checked
func MaxSubarray[T integer](values []T) (sum T, start, end int, err error) {
	if len(values) == 0 {
		return 0, 0, 0, fmt.Errorf("%w: values", ErrEmpty)
	}
	best, bestStart, bestEnd := values[0], 0, 1
	cur, curStart := values[0], 0 // cur = largest sum of a subarray ending at i, starting at curStart
	for i := 1; i < len(values); i++ {
		v := values[i]
		if cur > 0 {
			s, ok := checkedAdd(cur, v)
			if !ok {
				return 0, 0, 0, fmt.Errorf("%w: sum of values[%d:%d]", ErrOverflow, curStart, i+1)
			}
			cur = s
		} else {
			// A running sum <= 0 never helps: restart here. Restarting on 0 keeps the subarray shortest.
			cur, curStart = v, i
		}
		if cur > best {
			best, bestStart, bestEnd = cur, curStart, i+1
		}
	}
	return best, bestStart, bestEnd, nil
}
