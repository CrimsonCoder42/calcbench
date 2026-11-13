package twopointer

// IsPalindrome reports whether values reads the same forwards and backwards: values[i] == values[n-1-i]
// for every i. The empty slice and a single element are palindromes. Equality is ==, so a slice with a
// floating-point NaN off the centre is never a palindrome. For strings, see algo/text.
// values is read, never modified.
//
// ID:      ALG-TWOPTR-005
// Time:    O(n)   n = len(values); at most ⌊n/2⌋ comparisons
// Space:   O(1), no allocation
// Source:  Manacher, "A New Linear-Time 'On-Line' Algorithm for Finding the Smallest Initial Palindrome of a String", Journal of the ACM, 1975, 22(3):346–351, doi:10.1145/321892.321896 (the definition: a sequence equal to its reversal)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", p. 77 (two pointers that each move in one direction only), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    two-pointers, palindrome, symmetric, reverse
func IsPalindrome[T comparable](values []T) bool {
	for lo, hi := 0, len(values)-1; lo < hi; lo, hi = lo+1, hi-1 {
		if values[lo] != values[hi] {
			return false
		}
	}
	return true
}
