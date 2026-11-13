package hashing

import "errors"

// ErrNilKeyFunc: the key function is nil.
var ErrNilKeyFunc = errors.New("hashing: key function must not be nil")

// Group is one key and the values that map to it, in input order.
type Group[K comparable, T any] struct {
	Key   K
	Items []T
}

// GroupBy returns values grouped by key(value). Groups are in the order in which their key first
// appears; each group's Items keep input order. key is called exactly once per value, in order.
// Empty input gives an empty, non-nil slice. A NaN key matches nothing, so each value with a NaN key
// forms its own group. When K is an interface type, key must return comparable dynamic values (the Go
// map rule). values is read, never modified.
// Errors: ErrNilKeyFunc when key is nil.
//
// ID:      ALG-HASH-004
// Time:    O(n) expected   n = len(values), plus n calls of key
// Space:   O(g) extra for the key index   g = number of distinct keys; O(n) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary INSERT and SEARCH)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §3.4
// Tags:    hash, map, group-by, partition, bucket, stable, deterministic
func GroupBy[T any, K comparable](values []T, key func(T) K) ([]Group[K, T], error) {
	if key == nil {
		return nil, ErrNilKeyFunc
	}
	groups := []Group[K, T]{}
	index := make(map[K]int)
	for _, v := range values {
		k := key(v)
		g, ok := index[k]
		if !ok {
			g = len(groups)
			// A NaN key never matches itself, so it is not recorded; each NaN keyed value is its own group.
			if k == k {
				index[k] = g
			}
			groups = append(groups, Group[K, T]{Key: k})
		}
		groups[g].Items = append(groups[g].Items, v)
	}
	return groups, nil
}
