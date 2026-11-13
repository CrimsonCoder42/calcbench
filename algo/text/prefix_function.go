package text

import "slices"

// PrefixFunction appends to dst the prefix function of s and returns the extended slice: len(s) values,
// where value q is the length of the longest proper prefix of s[:q+1] that is also a suffix of it (its
// longest border). This is the failure function of the Knuth–Morris–Pratt matcher. Bytes are compared;
// s is read, never modified.
//
// ID:      ALG-TEXT-001
// Time:    O(n)   n = len(s); at most 2n byte comparisons (amortised: the border grows at most n times)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Knuth, Morris, Pratt, "Fast Pattern Matching in Strings", SIAM Journal on Computing, 1977, 6(2):323–350, doi:10.1137/0206024
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.4, COMPUTE-PREFIX-FUNCTION and Figure 32.10
// Tags:    string, kmp, prefix-function, failure-function, border, period
func PrefixFunction[S ~string | ~[]byte](dst []int, s S) []int {
	n := len(s)
	base := len(dst)
	dst = slices.Grow(dst, n)[:base+n]
	pi := dst[base:]
	if n == 0 {
		return dst
	}
	pi[0] = 0
	k := 0
	for q := 1; q < n; q++ {
		for k > 0 && s[k] != s[q] {
			k = pi[k-1]
		}
		if s[k] == s[q] {
			k++
		}
		pi[q] = k
	}
	return dst
}
