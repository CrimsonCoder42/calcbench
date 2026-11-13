package text

import "strings"

// ReverseWords returns the words of s in reverse order, joined by single spaces, with no leading or
// trailing space. A word is a maximal run of bytes that are not ASCII white space (space, \t, \n, \v,
// \f, \r); the bytes inside each word keep their order, so UTF-8 words come out intact. A string with no
// words gives "". s is read, never modified.
//
// ID:      ALG-TEXT-013
// Time:    O(n)   n = len(s); two passes
// Space:   O(1) extra; O(n) for the result, one allocation
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2.3 (reversal: reverse the whole, then each piece, so the pieces swap order and each reads forwards)
// Tags:    string, words, reverse, reversal, tokens, whitespace
func ReverseWords(s string) string {
	size, words := 0, 0
	for i := 0; i < len(s); {
		if isASCIISpace(s[i]) {
			i++
			continue
		}
		j := i
		for j < len(s) && !isASCIISpace(s[j]) {
			j++
		}
		size += j - i
		words++
		i = j
	}
	if words == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(size + words - 1)
	for j := len(s); j > 0; {
		if isASCIISpace(s[j-1]) {
			j--
			continue
		}
		i := j
		for i > 0 && !isASCIISpace(s[i-1]) {
			i--
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(s[i:j])
		j = i
	}
	return b.String()
}

// isASCIISpace reports whether c is one of the six white-space characters of the C "C" locale
// (ISO/IEC 9899:2011 §7.4.1.10): space, \t, \n, \v, \f, \r.
func isASCIISpace(c byte) bool {
	return c == ' ' || (c >= '\t' && c <= '\r')
}
