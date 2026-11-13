package window

import (
	"errors"
	"fmt"
)

// ErrNegativeK: k is less than 0.
var ErrNegativeK = errors.New("window: k must not be negative")

// LongestWithAtMostKDistinct returns the longest run values[start : start+length] holding at most k
// distinct values, as (start, length). Among longest runs it returns the one that starts first. k = 0
// or empty values gives (0, 0). Equality is ==, so each floating-point NaN counts as a distinct value.
// Two pointers move right only, with a count per distinct value in the window.
// values is read, never modified.
// Errors: ErrNegativeK when k < 0.
//
// ID:      ALG-WINDOW-012
// Time:    O(n) expected   n = len(values); each pointer moves at most n times, one hash-map update per move
// Space:   O(min(n, k+1)) extra (counts of the distinct values in the window, at most k+1 at once)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", p. 77 (both pointers move one way only, so O(n) moves in total), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 "Hash Tables" (O(1) expected search, insert, delete)
// Tags:    window, sliding-window, two-pointers, variable-window, longest, distinct, hash-map
func LongestWithAtMostKDistinct[T comparable](values []T, k int) (start, length int, err error) {
	if k < 0 {
		return 0, 0, fmt.Errorf("%w: k = %d", ErrNegativeK, k)
	}
	if k == 0 || len(values) == 0 {
		return 0, 0, nil
	}
	counts := make(map[T]int, min(len(values), k)) // k+1 could overflow
	distinct, lo := 0, 0
	for i, v := range values {
		if v != v { // NaN: never equal to itself, so it cannot be counted in the map
			distinct++
		} else {
			if counts[v] == 0 {
				distinct++
			}
			counts[v]++
		}
		for distinct > k {
			u := values[lo]
			lo++
			if u != u {
				distinct--
				continue
			}
			if counts[u]--; counts[u] == 0 {
				delete(counts, u)
				distinct--
			}
		}
		if i+1-lo > length {
			start, length = lo, i+1-lo
		}
	}
	return start, length, nil
}
