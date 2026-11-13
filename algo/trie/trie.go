package trie

import (
	"iter"
	"slices"
)

// Trie maps string keys to values of type V by a prefix tree over the keys' bytes: one node per
// distinct key prefix, each edge labelled by one byte. A node's children are kept in ascending byte
// order, so every walk visits keys in byte-wise lexicographic order (the order of Go's < on strings).
//
// The zero value is an empty trie ready to use. Every node left in the trie lies on the path to at least
// one key (Delete prunes branches that no longer lead to a key). Keys are bytes: a UTF-8 string and the
// same text in another normalisation are different keys. The empty string is a valid key.
// Methods on a nil *Trie behave as on an empty trie, except Insert, which stores nothing.
// A Trie is not safe for concurrent use.
type Trie[V any] struct {
	root node[V]
	n    int // number of keys
}

// node is one trie node: the value of the key ending here (when has is true) and the children, labels[j]
// being the byte on the edge to kids[j]; labels is strictly ascending.
type node[V any] struct {
	labels []byte
	kids   []*node[V]
	val    V
	has    bool
}

// child returns the child under byte b and its position, or (nil, the insertion position).
// O(log σ), σ = number of children (at most 256, so at most 8 comparisons).
func (nd *node[V]) child(b byte) (*node[V], int) {
	j, found := slices.BinarySearch(nd.labels, b)
	if !found {
		return nil, j
	}
	return nd.kids[j], j
}

// find returns the node reached by key's bytes from the root, or nil.
func (t *Trie[V]) find(key string) *node[V] {
	if t == nil {
		return nil
	}
	nd := &t.root
	for i := 0; i < len(key); i++ {
		if nd, _ = nd.child(key[i]); nd == nil {
			return nil
		}
	}
	return nd
}

// Insert stores value under key and reports whether key was already present (its old value is then
// replaced). On a nil *Trie it stores nothing and returns false.
//
// ID:      ALG-TRIE-001
// Time:    O(m)   m = len(key); at most 8 byte comparisons per level plus, for a new branch, one
// insertion into a sorted child list of at most 256 entries
// Space:   O(m) for the new nodes when key shares no prefix with a stored key; O(1) when key's path exists
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (trie insertion)
// Tags:    trie, prefix-tree, insert, dictionary, string-keys
func (t *Trie[V]) Insert(key string, value V) (replaced bool) {
	if t == nil {
		return false
	}
	nd := &t.root
	for i := 0; i < len(key); i++ {
		next, j := nd.child(key[i])
		if next == nil {
			next = &node[V]{}
			nd.labels = slices.Insert(nd.labels, j, key[i])
			nd.kids = slices.Insert(nd.kids, j, next)
		}
		nd = next
	}
	replaced = nd.has
	nd.val, nd.has = value, true
	if !replaced {
		t.n++
	}
	return replaced
}

// Search returns the value stored under key, with ok true; when key is absent it returns the zero V
// and false.
//
// ID:      ALG-TRIE-002
// Time:    O(m)   m = len(key); at most 8 byte comparisons per level
// Space:   O(1)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (Algorithm T, trie search)
// Tags:    trie, prefix-tree, search, lookup, dictionary
func (t *Trie[V]) Search(key string) (value V, ok bool) {
	if nd := t.find(key); nd != nil && nd.has {
		return nd.val, true
	}
	return value, false
}

// HasPrefix reports whether at least one stored key starts with prefix. HasPrefix("") is true exactly
// when the trie holds a key.
//
// ID:      ALG-TRIE-003
// Time:    O(m)   m = len(prefix)
// Space:   O(1)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching"
// Tags:    trie, prefix-tree, prefix, starts-with
func (t *Trie[V]) HasPrefix(prefix string) bool {
	// Every node lies on the path to a key (Delete prunes), so reaching a node other than the root
	// proves a key below it; the root needs the count.
	if prefix == "" {
		return t.Len() > 0
	}
	return t.find(prefix) != nil
}

// Delete removes key and returns its value, with ok true; when key is absent it changes nothing and
// returns the zero V and false. Nodes that no longer lead to any key are unlinked (the branch is pruned
// at the deepest ancestor that still holds a key or another child), so memory follows the keys held.
//
// ID:      ALG-TRIE-004
// Time:    O(m)   m = len(key); one walk down, no walk back up
// Space:   O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (deletion from a trie)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Tags:    trie, prefix-tree, delete, prune
func (t *Trie[V]) Delete(key string) (value V, ok bool) {
	if t == nil {
		return value, false
	}
	// cut, cutAt: the deepest node on the path that must stay (the root, a node holding a key, or a
	// node with another child) and the position of the path's child under it. Removing that one edge
	// removes exactly the nodes that lead only to key.
	nd := &t.root
	var cut *node[V]
	cutAt := 0
	for i := 0; i < len(key); i++ {
		next, j := nd.child(key[i])
		if next == nil {
			return value, false
		}
		if nd == &t.root || nd.has || len(nd.kids) > 1 {
			cut, cutAt = nd, j
		}
		nd = next
	}
	if !nd.has {
		return value, false
	}
	value = nd.val
	var zero V
	nd.val, nd.has = zero, false // drop the reference so the garbage collector can free it
	t.n--
	if len(nd.kids) == 0 && cut != nil {
		cut.labels = slices.Delete(cut.labels, cutAt, cutAt+1)
		cut.kids[cutAt] = nil
		cut.kids = slices.Delete(cut.kids, cutAt, cutAt+1)
	}
	return value, true
}

// Len returns the number of keys stored (0 for a nil *Trie).
//
// ID:      ALG-TRIE-005
// Time:    O(1)
// Space:   O(1)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Tags:    trie, size
func (t *Trie[V]) Len() int {
	if t == nil {
		return 0
	}
	return t.n
}

// WithPrefix returns an iterator over every stored key that starts with prefix, with its value, in
// ascending byte order. The key passed to yield is a new string. The trie must not be changed while
// the iteration runs.
//
// ID:      ALG-TRIE-006
// Time:    O(m + s)   m = len(prefix), s = number of trie nodes below the prefix (at most the total
// length of the matching keys' remainders); stopping early stops the walk
// Space:   O(d) extra for the walk's stack and key buffer, d = length of the longest matching key;
// one string per key yielded
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (a trie's subtree holds exactly the keys with that prefix, in lexicographic order)
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400
// Tags:    trie, prefix-tree, prefix-scan, iterator, range-over-func, lexicographic
func (t *Trie[V]) WithPrefix(prefix string) iter.Seq2[string, V] {
	return func(yield func(string, V) bool) {
		start := t.find(prefix)
		if start == nil {
			return
		}
		walk(start, []byte(prefix), func(key []byte, nd *node[V]) bool {
			return yield(string(key), nd.val)
		})
	}
}

// frame is one level of an iterative depth-first walk: the node and the next child position to visit.
type frame[V any] struct {
	nd   *node[V]
	next int
}

// walk visits, in ascending byte order, every node below and including start that holds a key, calling
// visit with the full key (buf, which visit must not keep) and the node, and stops when visit returns
// false. buf holds start's key on entry. Iterative, so key length never limits it by stack depth.
func walk[V any](start *node[V], buf []byte, visit func(key []byte, nd *node[V]) bool) bool {
	if start.has && !visit(buf, start) {
		return false
	}
	stack := []frame[V]{{nd: start}}
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.next == len(top.nd.kids) {
			stack = stack[:len(stack)-1]
			if len(stack) > 0 { // the root of the walk owns no byte of buf to remove
				buf = buf[:len(buf)-1]
			}
			continue
		}
		b, kid := top.nd.labels[top.next], top.nd.kids[top.next]
		top.next++
		buf = append(buf, b)
		if kid.has && !visit(buf, kid) {
			return false
		}
		stack = append(stack, frame[V]{nd: kid})
	}
	return true
}
