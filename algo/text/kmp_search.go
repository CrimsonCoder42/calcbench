package text

import (
	"errors"
	"fmt"
)

// ErrEmptyPattern: a search was given an empty pattern, which would match at every position.
var ErrEmptyPattern = errors.New("text: pattern is empty")

// KMPSearch appends to dst the start of every occurrence of pattern in text, ascending and overlapping
// occurrences included, and returns the extended slice. Bytes are compared; offsets are byte offsets.
// text and pattern are read, never modified.
// Errors: ErrEmptyPattern when pattern is empty; dst is then returned unchanged.
//
// ID:      ALG-TEXT-002
// Time:    O(n + m)   n = len(text), m = len(pattern); at most 2n + 2m byte comparisons
// Space:   O(m) extra (the prefix function of pattern); O(k) for the result, k = occurrences
// Source:  Knuth, Morris, Pratt, "Fast Pattern Matching in Strings", SIAM Journal on Computing, 1977, 6(2):323–350, doi:10.1137/0206024
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.4, KMP-MATCHER
// Tags:    string, search, pattern-matching, kmp, find-all, substring
func KMPSearch[S ~string | ~[]byte](dst []int, text, pattern S) ([]int, error) {
	m := len(pattern)
	if m == 0 {
		return dst, fmt.Errorf("pattern: %w", ErrEmptyPattern)
	}
	if m > len(text) {
		return dst, nil
	}
	pi := PrefixFunction(make([]int, 0, m), pattern)
	q := 0 // bytes of pattern matched so far
	for i := 0; i < len(text); i++ {
		for q > 0 && pattern[q] != text[i] {
			q = pi[q-1]
		}
		if pattern[q] == text[i] {
			q++
		}
		if q == m {
			dst = append(dst, i-m+1)
			q = pi[q-1]
		}
	}
	return dst, nil
}
