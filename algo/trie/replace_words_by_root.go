package trie

import "slices"

// ShortestPrefixOf returns the shortest stored key that is a prefix of word (word itself included), with
// its value and ok true; when no stored key is a prefix of word it returns "", the zero V and false.
// The returned key is word[:len(key)], a substring of word: no allocation.
//
// ID:      ALG-TRIE-009
// Time:    O(p)   p = length of the returned key (len(word) when none is found); stops at the first key
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (the search path passes every stored prefix of the argument, shortest first)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Tags:    trie, shortest-prefix, prefix-match, stemming, root-word
func (t *Trie[V]) ShortestPrefixOf(word string) (key string, value V, ok bool) {
	if t == nil {
		return "", value, false
	}
	nd := &t.root
	for i := 0; ; i++ {
		if nd.has {
			return word[:i], nd.val, true
		}
		if i == len(word) {
			return "", value, false
		}
		if nd, _ = nd.child(word[i]); nd == nil {
			return "", value, false
		}
	}
}

// ReplaceWordsByRoot appends to dst, for each word in words, the shortest root in roots that is a prefix
// of the word, or the word itself when no root is, and returns the extended slice (len(words) strings
// appended, in order). Each string appended is a substring of its word (no copying of text).
// Empty strings in roots are ignored: an empty root is a prefix of every word and would turn every word
// into "". Duplicate roots are harmless. roots and words are read, never modified.
//
// ID:      ALG-TRIE-010
// Time:    O(R + W)   R = total bytes of roots (building the trie), W = total bytes of words (worst
// case; each word stops at its first root)
// Space:   O(R) extra for the trie of roots; O(len(words)) for the result, one allocation only if dst
// lacks capacity
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching"
// Tags:    trie, replace-words, root-word, stemming, shortest-prefix, dictionary
func ReplaceWordsByRoot(dst, roots, words []string) []string {
	var dict Trie[struct{}]
	for _, r := range roots {
		if r != "" {
			dict.Insert(r, struct{}{})
		}
	}
	dst = slices.Grow(dst, len(words))
	for _, w := range words {
		if root, _, ok := dict.ShortestPrefixOf(w); ok {
			w = root
		}
		dst = append(dst, w)
	}
	return dst
}
