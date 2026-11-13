package text

import "unicode/utf8"

// IsAnagram reports whether a and b hold the same multiset of runes: each is a rearrangement of the
// other. Comparison is exact (case, accents and spaces all count). Bytes that are not valid UTF-8 each
// read as U+FFFD. Two empty strings are anagrams. a and b are read, never modified.
//
// ID:      ALG-TEXT-010
// Time:    O(n) expected   n = len(a) + len(b) in bytes; O(n) worst case when both are ASCII
// Space:   O(1) extra when both are ASCII (a 128-entry count table, no allocation); otherwise O(d) for a count map, d = distinct runes
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2, problem C (anagrams share a signature, §2.4 and §2.8)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.2 (counting occurrences of each key) and ch. 11 (hash tables)
// Tags:    string, anagram, multiset, counting, permutation
func IsAnagram(a, b string) bool {
	if isASCII(a) && isASCII(b) {
		if len(a) != len(b) {
			return false
		}
		var count [utf8.RuneSelf]int
		for i := range len(a) {
			count[a[i]]++
			count[b[i]]--
		}
		for _, c := range count {
			if c != 0 {
				return false
			}
		}
		return true
	}
	count := make(map[rune]int)
	for _, r := range a {
		count[r]++
	}
	for _, r := range b {
		c := count[r]
		if c == 0 {
			return false
		}
		if c == 1 {
			delete(count, r)
		} else {
			count[r] = c - 1
		}
	}
	return len(count) == 0
}

// isASCII reports whether every byte of s is below 0x80.
func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}
