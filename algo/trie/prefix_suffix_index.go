package trie

import "slices"

// separator is the end-of-suffix symbol of the permuterm keys: symbol 256, outside the byte range, so it
// can never be confused with a byte of a word (the role of "$" in Manning et al. §3.2.1).
const separator = 256

// PrefixSuffixIndex answers "which words start with p and end with s" over a fixed list of words, by a
// permuterm-style trie: for every word w (at index i in the list) and every split point j, the key
// w[j:] + $ + w is stored ($ a symbol outside the byte range). A word starts with p and ends with s
// exactly when the key s + $ + w exists and has s + $ + p as a prefix, so one prefix walk answers the
// query. Prefix and suffix may overlap inside a word (prefix "ab", suffix "bc" match "abc"): the
// condition is strings.HasPrefix(w, p) && strings.HasSuffix(w, s), on bytes.
//
// Built by NewPrefixSuffixIndex; read-only after that, so safe for concurrent lookups. The zero value
// and a nil *PrefixSuffixIndex are empty indexes.
type PrefixSuffixIndex struct {
	root psNode
}

// psNode is one node: children by symbol (0–255 a byte, 256 the separator) in ascending order, the
// indexes of the words whose key ends here, and the largest word index whose key passes through.
type psNode struct {
	labels []uint16
	kids   []*psNode
	ends   []int
	last   int
}

func (nd *psNode) child(sym uint16) (*psNode, int) {
	j, found := slices.BinarySearch(nd.labels, sym)
	if !found {
		return nil, j
	}
	return nd.kids[j], j
}

// NewPrefixSuffixIndex returns the index of words; a word is identified by its position in words.
// Duplicate words are kept as separate positions. words is read, never modified, and not kept (the
// index holds its own copy of the bytes, one symbol per node).
//
// ID:      ALG-TRIE-011
// Time:    O(Σ Lᵢ²)   Lᵢ = len(words[i]); each word adds Lᵢ+1 keys of at most 2Lᵢ+1 symbols
// Space:   O(Σ Lᵢ²) nodes for the index (fewer where words share suffix-and-word paths)
// Source:  Manning, Raghavan, Schütze, Introduction to Information Retrieval, Cambridge University Press, 2008, §3.2.1 "General wildcard queries", Permuterm indexes (a query X*Y is answered by a prefix lookup of Y$X in a dictionary of rotated terms; read in the 2007 preliminary draft, §3.1.1, pp. 48–49)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Tags:    trie, permuterm, prefix-and-suffix, wildcard, word-filter, index
func NewPrefixSuffixIndex(words []string) *PrefixSuffixIndex {
	x := &PrefixSuffixIndex{}
	x.root.last = -1
	for wi, w := range words {
		for j := 0; j <= len(w); j++ {
			nd := &x.root
			nd.last = wi
			nd = nd.step(w[j:], wi)
			nd = nd.add(separator, wi)
			nd = nd.step(w, wi)
			nd.ends = append(nd.ends, wi)
		}
	}
	return x
}

// step follows (adding where missing) the bytes of s from nd, recording word index wi as the latest on
// every node passed, and returns the last node.
func (nd *psNode) step(s string, wi int) *psNode {
	for i := 0; i < len(s); i++ {
		nd = nd.add(uint16(s[i]), wi)
	}
	return nd
}

// add returns nd's child under sym, creating it if missing, with wi recorded as its latest word.
func (nd *psNode) add(sym uint16, wi int) *psNode {
	next, j := nd.child(sym)
	if next == nil {
		next = &psNode{}
		nd.labels = slices.Insert(nd.labels, j, sym)
		nd.kids = slices.Insert(nd.kids, j, next)
	}
	next.last = wi // words are added in index order, so the latest is the largest
	return next
}

// locate returns the node of the query key suffix + $ + prefix, or nil when no word matches.
func (x *PrefixSuffixIndex) locate(prefix, suffix string) *psNode {
	if x == nil {
		return nil
	}
	nd := &x.root
	for i := 0; i < len(suffix) && nd != nil; i++ {
		nd, _ = nd.child(uint16(suffix[i]))
	}
	if nd == nil {
		return nil
	}
	nd, _ = nd.child(separator)
	for i := 0; i < len(prefix) && nd != nil; i++ {
		nd, _ = nd.child(uint16(prefix[i]))
	}
	return nd
}

// Lookup appends to dst, in ascending order, the index of every word that starts with prefix and ends
// with suffix (bytes; the two may overlap), and returns the extended slice. Empty prefix and suffix
// match every word.
//
// ID:      ALG-TRIE-012
// Time:    O(p + q + s + r log r)   p = len(prefix), q = len(suffix), s = nodes below the query key
// (at most the total length of the matching words), r = number of matches (sorting them)
// Space:   O(h) extra for the walk's stack, h = number of nodes pending at once (at most s); O(r) for
// the result
// Source:  Manning, Raghavan, Schütze, Introduction to Information Retrieval, Cambridge University Press, 2008, §3.2.1 "General wildcard queries", Permuterm indexes (lookup of Y$X for the query X*Y)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (all keys with a prefix lie in one subtree)
// Tags:    trie, permuterm, prefix-and-suffix, wildcard, word-filter
func (x *PrefixSuffixIndex) Lookup(dst []int, prefix, suffix string) []int {
	start := x.locate(prefix, suffix)
	if start == nil {
		return dst
	}
	first := len(dst)
	stack := []*psNode{start}
	for len(stack) > 0 {
		nd := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		dst = append(dst, nd.ends...)
		stack = append(stack, nd.kids...)
	}
	// Each matching word has exactly one key below start (the split at len(word)-len(suffix)), so
	// there are no duplicates to remove.
	slices.Sort(dst[first:])
	return dst
}

// LookupLast returns the largest index of a word that starts with prefix and ends with suffix, with ok
// true; when no word matches it returns -1 and false.
//
// ID:      ALG-TRIE-013
// Time:    O(p + q)   p = len(prefix), q = len(suffix); at most 9 comparisons per symbol
// Space:   O(1)
// Source:  Manning, Raghavan, Schütze, Introduction to Information Retrieval, Cambridge University Press, 2008, §3.2.1 "General wildcard queries", Permuterm indexes
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Tags:    trie, permuterm, prefix-and-suffix, word-filter, latest
func (x *PrefixSuffixIndex) LookupLast(prefix, suffix string) (index int, ok bool) {
	nd := x.locate(prefix, suffix)
	if nd == nil {
		return -1, false
	}
	return nd.last, true
}
