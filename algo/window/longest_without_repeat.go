package window

// LongestWithoutRepeat returns the longest run values[start : start+length] in which no value occurs
// twice, as (start, length). Among longest runs it returns the one that starts first. Empty values
// gives (0, 0). Equality is ==, so a floating-point NaN never repeats.
// Two pointers move right only: the map holds each value's last index, and when the entering value was
// last seen inside the window the left pointer jumps past that index.
// values is read, never modified.
//
// ID:      ALG-WINDOW-013
// Time:    O(n) expected   n = len(values); one hash-map lookup and update per value
// Space:   O(d) extra   d = number of distinct values (the last-index map)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", p. 77 (both pointers move one way only), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 "Hash Tables" (O(1) expected search and insert)
// Tags:    window, sliding-window, two-pointers, variable-window, longest, distinct, unique, hash-map
func LongestWithoutRepeat[T comparable](values []T) (start, length int) {
	last := make(map[T]int)
	lo := 0
	for i, v := range values {
		if v == v { // a NaN is never a repeat and cannot be found in the map
			if j, seen := last[v]; seen && j >= lo {
				lo = j + 1
			}
			last[v] = i
		}
		if i+1-lo > length {
			start, length = lo, i+1-lo
		}
	}
	return start, length
}
