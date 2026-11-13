package text

// LongestPalindromicSubstring returns the longest substring of s that reads the same forwards and
// backwards, as s[start : start+length]. Among longest ones it returns the leftmost. An empty s gives
// (0, 0); any other s gives length >= 1. Bytes are compared (on UTF-8 text a palindrome of bytes is not
// always one of runes). s is read, never modified. It reads the centres of PalindromeRadii (ALG-TEXT-008).
//
// ID:      ALG-TEXT-009
// Time:    O(n)   n = len(s)
// Space:   O(n) extra, one allocation (the 2n+1 radii of PalindromeRadii)
// Source:  Manacher, "A New Linear-Time 'On-Line' Algorithm for Finding the Smallest Initial Palindrome of a String", Journal of the ACM, 1975, 22(3):346–351, doi:10.1145/321892.321896
// Source:  Apostolico, Breslauer, Galil, "Parallel detection of all palindromes in a string", Theoretical Computer Science, 1995, 141:163–173, doi:10.1016/0304-3975(94)00083-U
// Tags:    string, palindrome, manacher, longest, substring, linear-time
func LongestPalindromicSubstring[S ~string | ~[]byte](s S) (start, length int) {
	if len(s) == 0 {
		return 0, 0
	}
	radii := PalindromeRadii(make([]int, 0, 2*len(s)+1), s)
	for c, l := range radii {
		if l > length {
			start, length = (c+1-l)/2, l
		}
	}
	return start, length
}
