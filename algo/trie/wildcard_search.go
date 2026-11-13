package trie

import "iter"

// Wildcard is the pattern byte that matches any one byte in WildcardSearch.
const Wildcard = '.'

// WildcardSearch returns an iterator over every stored key that matches pattern, with its value, in
// ascending byte order. A key matches when it has exactly len(pattern) bytes and each byte equals the
// pattern's byte at that position, where the pattern byte '.' (Wildcard) matches any byte.
// Semantics are per BYTE, not per rune: a character that UTF-8 encodes in 2–4 bytes needs 2–4 '.' to be
// skipped, and a '.' never matches a whole multi-byte character. '.' has no escape: a literal '.' in a
// key is matched by '.' like any other byte. To ask only whether any key matches, range once and stop.
// The key passed to yield is a new string. The trie must not be changed while the iteration runs.
//
// ID:      ALG-TRIE-008
// Time:    O(m) with no '.' (m = len(pattern)); worst case O(s), s = number of trie nodes at depth ≤ m
// (each node is entered at most once; every '.' branches into all children)
// Space:   O(m) extra for the walk's stack and key buffer; one string per key yielded
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (trie search, one branch per key position)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Tags:    trie, wildcard, pattern-match, dont-care, partial-match, iterator, range-over-func
func (t *Trie[V]) WildcardSearch(pattern string) iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		if t == nil {
			return
		}
		m := len(pattern)
		// stack[d] is the node at depth d and the next child position to try there. For a literal
		// pattern byte the only candidate is the matching child: next is set to its position and the
		// frame finishes after it.
		if m == 0 {
			if t.root.has {
				yield("", t.root.val)
			}
			return
		}
		stack := make([]frame[V], 1, m+1)
		stack[0] = frame[V]{nd: &t.root}
		key := make([]byte, 0, m)    // key[:d] is the key of the node at depth d
		limit := make([]int, 1, m+1) // limit[d]: one past the last child position to try at depth d
		limit[0] = initFrame(&stack[0], pattern[0])
		for len(stack) > 0 {
			d := len(stack) - 1
			top := &stack[d]
			if top.next >= limit[d] {
				stack = stack[:d]
				limit = limit[:d]
				continue
			}
			b, kid := top.nd.labels[top.next], top.nd.kids[top.next]
			top.next++
			key = append(key[:d], b)
			if d+1 == m { // a full-length match: report it, never descend further
				if kid.has && !yield(string(key), kid.val) {
					return
				}
				continue
			}
			stack = append(stack, frame[V]{nd: kid})
			limit = append(limit, initFrame(&stack[d+1], pattern[d+1]))
		}
	}
}

// initFrame sets f.next to the first child position to try for pattern byte p and returns one past the
// last: all children for Wildcard, the single matching child (or none) for a literal byte.
func initFrame[V any](f *frame[V], p byte) int {
	if p == Wildcard {
		f.next = 0
		return len(f.nd.kids)
	}
	if _, j := f.nd.child(p); j < len(f.nd.labels) && f.nd.labels[j] == p {
		f.next = j
		return j + 1
	}
	f.next = 0
	return 0
}
