package hashing

// SetDifference returns the distinct values of a that do not occur in b, in order of first appearance in
// a. A NaN equals nothing, so every NaN of a appears in the result. An empty result is a non-nil, empty
// slice. a and b are read, never modified.
//
// ID:      ALG-HASH-008
// Time:    O(|a| + |b|) expected
// Space:   O(|a| + |b|) extra for the set of b and the seen-set; O(|a|) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary operations)
// Source:  Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §6.4 (hashing)
// Tags:    hash, set, difference, minus, except, distinct, deterministic
func SetDifference[T comparable](a, b []T) []T {
	// One map does both jobs: values of b, and values of a already emitted, are both "skip".
	skip := make(map[T]struct{}, len(b))
	for _, v := range b {
		skip[v] = struct{}{}
	}
	out := []T{}
	for _, v := range a {
		if _, ok := skip[v]; ok {
			continue
		}
		skip[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
