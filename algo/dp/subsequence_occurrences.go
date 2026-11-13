package dp

import "fmt"

// SubsequenceOccurrences returns the number of ways pattern occurs in text as a subsequence: the
// number of index sets i1 < i2 < … < im with text[ik] = pattern[k] (the binomial coefficient of words,
// "text choose pattern"). An empty pattern occurs once; a pattern longer than text occurs 0 times.
// Elements match by ==, so a NaN matches nothing. Inputs are read, never modified.
// Errors: ErrOverflow when the count, or the count for a prefix of pattern that the recurrence holds on
// the way (as L28), exceeds 2⁶⁴−1.
//
// ID:      ALG-DP-025
// Time:    O(n·m)   n = len(text), m = len(pattern)
// Space:   O(m) extra (one row of m+1 counts, one allocation)
// Source:  Lothaire, Combinatorics on Words, Addison-Wesley, 1983 (reissued Cambridge University Press, 1997), ch. 6 "Subwords" (Sakarovitch, Simon), §6.3 (binomial coefficients of words: (ua choose vb) = (u choose vb) + δ_{a,b}·(u choose v), (u choose ε) = 1)
// Tags:    dp, subsequence, counting, distinct-subsequences, binomial-of-words, string
func SubsequenceOccurrences[T comparable](text, pattern []T) (uint64, error) {
	m := len(pattern)
	if m > len(text) {
		return 0, nil
	}
	ways := make([]uint64, m+1) // ways[k] = occurrences of pattern[:k] in the text read so far
	ways[0] = 1
	for i, x := range text {
		// Descending k so ways[k−1] is still the count before x (each text position used once).
		for k := min(m, i+1); k >= 1; k-- {
			if pattern[k-1] == x {
				s, ok := addUint64(ways[k], ways[k-1])
				if !ok {
					return 0, fmt.Errorf("%w: occurrences of pattern[:%d] in text[:%d]", ErrOverflow, k, i+1)
				}
				ways[k] = s
			}
		}
	}
	return ways[m], nil
}
