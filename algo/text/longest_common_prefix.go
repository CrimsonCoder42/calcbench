package text

// LongestCommonPrefix returns the length in bytes of the longest prefix shared by every string in strs:
// strs[0][:result] is that prefix. No strings gives 0; one string gives its length. Bytes are compared,
// so on UTF-8 text the prefix may end inside a multi-byte rune when two strings differ there.
// strs is read, never modified.
//
// ID:      ALG-TEXT-007
// Time:    O(k·p + k)   k = len(strs), p = the result; each string is read only up to the first column that differs
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 (digital searching: the prefix shared by a set of keys is the trie path above its first branching node, walked here column by column without building the trie)
// Tags:    string, prefix, longest-common-prefix, lcp, vertical-scan
func LongestCommonPrefix[S ~string | ~[]byte](strs []S) int {
	if len(strs) == 0 {
		return 0
	}
	first := strs[0]
	for col := 0; col < len(first); col++ {
		c := first[col]
		for _, s := range strs[1:] {
			if col >= len(s) || s[col] != c {
				return col
			}
		}
	}
	return len(first)
}
