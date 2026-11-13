package text

import "slices"

// PalindromeRadii appends to dst the length of the longest palindrome at every centre of s and returns
// the extended slice: 2n−1 values for n = len(s) > 0, none for an empty s. Value c is the length of the
// longest palindromic substring centred on byte c/2 when c is even (an odd length, at least 1), or
// centred between bytes (c−1)/2 and (c+1)/2 when c is odd (an even length, possibly 0). That palindrome
// starts at byte (c+1−length)/2. Bytes are compared. s is read, never modified.
// This is calcbench's one home of Manacher's algorithm.
//
// ID:      ALG-TEXT-008
// Time:    O(n)   n = len(s); the rightmost palindrome boundary only moves right, at most 2n+1 steps
// Space:   O(1) extra; O(n) for the result, one allocation of 2n+1 values only if dst lacks capacity
// Source:  Manacher, "A New Linear-Time 'On-Line' Algorithm for Finding the Smallest Initial Palindrome of a String", Journal of the ACM, 1975, 22(3):346–351, doi:10.1145/321892.321896
// Source:  Apostolico, Breslauer, Galil, "Parallel detection of all palindromes in a string", Theoretical Computer Science, 1995, 141:163–173, doi:10.1016/0304-3975(94)00083-U (Manacher's method finds every maximal palindrome, odd and even centres)
// Tags:    string, palindrome, manacher, centre, radius, linear-time
func PalindromeRadii[S ~string | ~[]byte](dst []int, s S) []int {
	n := len(s)
	if n == 0 {
		return dst
	}
	// Work on the implicit string t = |s₀|s₁|…|s_(n−1)| of length 2n+1, never built: t[i] is a separator
	// when i is even and s[(i−1)/2] when i is odd. The radius of the longest palindrome of t centred at i
	// equals the length of the longest palindrome of s with the matching centre.
	size := 2*n + 1
	base := len(dst)
	dst = slices.Grow(dst, size)[:base+size]
	p := dst[base:]
	l, r := 0, -1 // t[l..r] is the palindrome reaching furthest right so far
	for i := 0; i < size; i++ {
		k := 0
		if i <= r {
			k = min(p[l+r-i], r-i)
		}
		for i-k-1 >= 0 && i+k+1 < size && ((i-k-1)&1 == 0 || s[(i-k-2)/2] == s[(i+k)/2]) {
			k++
		}
		p[i] = k
		if i+k > r {
			l, r = i-k, i+k
		}
	}
	// Drop the two outer separators (always 0): centre c of s is position c+1 of t.
	copy(p, p[1:size-1])
	return dst[:base+size-2]
}
