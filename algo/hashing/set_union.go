package hashing

// SetUnion returns the distinct values of a in order of first appearance, followed by the values of b
// that do not occur in a, in order of first appearance in b. A NaN equals nothing, so every NaN of a and
// of b appears in the result. An empty result is a non-nil, empty slice. a and b are read, never modified.
//
// ID:      ALG-HASH-007
// Time:    O(|a| + |b|) expected
// Space:   O(|a| + |b|) extra for the seen-set; O(|a| + |b|) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary operations)
// Source:  Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §6.4 (hashing)
// Tags:    hash, set, union, merge, distinct, deduplicate, deterministic
func SetUnion[T comparable](a, b []T) []T {
	seen := make(map[T]struct{}, len(a))
	out := []T{}
	for _, list := range [2][]T{a, b} {
		for _, v := range list {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
