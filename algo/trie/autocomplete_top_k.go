package trie

import (
	"errors"
	"fmt"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/heap"
)

// ErrNegativeK: k is negative.
var ErrNegativeK = errors.New("trie: k must not be negative")

// Entry is one key of a Trie with its value, as returned by AutocompleteTopK.
type Entry[V any] struct {
	Key   string
	Value V
}

// AutocompleteTopK appends to dst the (at most) k stored keys that start with prefix, best first, with
// their values, and returns the extended slice.
// Ranking: rank(a, b) compares two values and returns a negative number when a ranks before b, positive
// when after, 0 when they tie (the shape of cmp.Compare; for "highest weight first" pass
// func(a, b V) int { return cmp.Compare(b.Weight, a.Weight) }). Ties go to the key that is smaller in
// byte order, so the result is fully determined. rank must be a strict weak order; one that is not
// gives an unspecified selection, never a panic. With rank nil the keys are ranked by byte order alone:
// the first k keys with the prefix, and the walk stops after the k-th.
// Errors: ErrNegativeK when k < 0 (dst returned unchanged). k = 0 or no matching key appends nothing.
//
// ID:      ALG-TRIE-007
// Time:    O(m + s log k)   m = len(prefix), s = number of trie nodes below the prefix; with rank nil,
// O(m + s') where s' is the number of nodes walked up to the k-th key
// Space:   O(k + d) extra (a heap of at most k entries; the walk's stack, d = length of the longest
// matching key); O(k) for the result. A key string is made only when its key enters the top k
// Source:  Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400 (the prefix subtree)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (a bounded priority queue keeps the k best of a stream)
// Tags:    trie, autocomplete, typeahead, top-k, prefix-search, heap
func (t *Trie[V]) AutocompleteTopK(dst []Entry[V], prefix string, k int, rank func(a, b V) int) ([]Entry[V], error) {
	if k < 0 {
		return dst, fmt.Errorf("%w: k = %d", ErrNegativeK, k)
	}
	start := t.find(prefix)
	if k == 0 || start == nil {
		return dst, nil
	}
	buf := []byte(prefix)
	if rank == nil {
		// Byte order is the walk order: the first k keys visited are the answer.
		taken := 0
		walk(start, buf, func(key []byte, nd *node[V]) bool {
			dst = append(dst, Entry[V]{string(key), nd.val})
			taken++
			return taken < k
		})
		return dst, nil
	}
	// worse(a, b): a ranks after b. The heap's top is the worst of the current best k. Ties break on
	// the key; the walk visits keys in ascending order, so a key arriving later loses every tie with a
	// key already held and enters only when strictly better than the top under rank.
	worse := func(a, b Entry[V]) bool {
		if c := rank(a.Value, b.Value); c != 0 {
			return c > 0
		}
		return a.Key > b.Key
	}
	best, _ := heap.NewHeap(worse, min(k, t.Len())) // worse is non-nil, so no error; a huge k never sizes the heap
	walk(start, buf, func(key []byte, nd *node[V]) bool {
		if best.Len() < k {
			best.Push(Entry[V]{string(key), nd.val})
			return true
		}
		if top, _ := best.Peek(); rank(nd.val, top.Value) < 0 {
			best.Pop()
			best.Push(Entry[V]{string(key), nd.val})
		}
		return true
	})
	// Pop yields worst first: fill the result back to front.
	first := len(dst)
	dst = slices.Grow(dst, best.Len())[:first+best.Len()]
	for i := len(dst) - 1; i >= first; i-- {
		dst[i], _ = best.Pop()
	}
	return dst, nil
}
