package text

import "slices"

// ZFunction appends to dst the Z values of s and returns the extended slice: len(s) values, where value
// i is the length of the longest common prefix of s and s[i:]. By convention value 0 is len(s) (the
// source leaves it undefined). Bytes are compared; s is read, never modified.
//
// ID:      ALG-TEXT-003
// Time:    O(n)   n = len(s); fewer than 2n byte comparisons that succeed or fail once per position
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Gusfield, Algorithms on Strings, Trees, and Sequences, Cambridge University Press, 1997, §1.3 (Z values, example aabcaabxaaz) and §1.4 (the Z algorithm)
// Tags:    string, z-function, z-algorithm, prefix, pattern-matching, period
func ZFunction[S ~string | ~[]byte](dst []int, s S) []int {
	n := len(s)
	base := len(dst)
	dst = slices.Grow(dst, n)[:base+n]
	z := dst[base:]
	if n == 0 {
		return dst
	}
	z[0] = n
	// [l, r) is the rightmost Z-box found so far: s[l:r] == s[0:r-l].
	l, r := 0, 0
	for i := 1; i < n; i++ {
		k := 0
		if i < r {
			k = min(z[i-l], r-i)
		}
		for i+k < n && s[k] == s[i+k] {
			k++
		}
		z[i] = k
		if i+k > r {
			l, r = i, i+k
		}
	}
	return dst
}
