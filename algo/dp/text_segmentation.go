package dp

// TextSegmentation splits text into a sequence of words from words, concatenated without separators,
// and returns the pieces in order (substrings of text, not copies); ok is false when no split exists.
// Words and text are compared as bytes. Empty text splits into no pieces (empty, non-nil, ok = true);
// empty dictionary words are ignored; a word may be used any number of times. Of several splits it
// returns the one whose last piece is longest, then (on the rest) the same rule again.
// Inputs are read, never modified.
//
// ID:      ALG-DP-015
// Time:    O(D + n·L) expected   n = len(text), L = the longest word in bytes, D = total bytes of words (building a trie of the dictionary in one hash map)
// Space:   O(D + n) extra (the trie and one back-pointer per position); O(p) for the result, p = pieces
// Source:  Erickson, Algorithms, 1st ed., self-published, 2019, ISBN 978-1-792-64483-2, §2.5 "Text Segmentation (Interpunctio Verborum)" and §3.3 "Interpunctio Verborum Redux" (Splittable(i): text[i:] splits iff some word text[i:j] is followed by a splittable text[j:]; memoised over i)
// Source:  Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, ch. 6, Exercise 5 (segmenting a string of letters into words by dynamic programming over prefixes)
// Tags:    dp, text-segmentation, word-break, trie, string, tokenise
func TextSegmentation(text string, words []string) (pieces []string, ok bool) {
	// The dictionary trie: node 0 is the root; edge (node, byte) → child lives in one map.
	type edge struct {
		node int
		b    byte
	}
	children := make(map[edge]int)
	isWord := []bool{false}
	maxLen := 0
	for _, w := range words {
		if w == "" {
			continue
		}
		maxLen = max(maxLen, len(w))
		node := 0
		for i := 0; i < len(w); i++ {
			e := edge{node, w[i]}
			child, found := children[e]
			if !found {
				child = len(isWord)
				isWord = append(isWord, false)
				children[e] = child
			}
			node = child
		}
		isWord[node] = true
	}
	n := len(text)
	// from[j] = start of the last piece of the chosen split of text[:j], or -1 if text[:j] does not split.
	from := make([]int, n+1)
	for j := 1; j <= n; j++ {
		from[j] = -1
	}
	for i := 0; i < n; i++ {
		if from[i] < 0 {
			continue
		}
		node := 0
		for j := i; j < n && j-i < maxLen; j++ {
			child, found := children[edge{node, text[j]}]
			if !found {
				break
			}
			node = child
			if isWord[node] && from[j+1] < 0 {
				from[j+1] = i // the smallest i: the longest last piece ending at j+1
			}
		}
	}
	if n > 0 && from[n] < 0 {
		return nil, false
	}
	count := 0
	for j := n; j > 0; j = from[j] {
		count++
	}
	pieces = make([]string, count)
	for j, k := n, count-1; j > 0; j, k = from[j], k-1 {
		pieces[k] = text[from[j]:j]
	}
	return pieces, true
}
