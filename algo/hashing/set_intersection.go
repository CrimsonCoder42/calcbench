package hashing

// SetIntersection returns the distinct values of a that also occur in b, in order of first appearance
// in a. The result has no duplicates; an empty result is a non-nil, empty slice. A NaN occurs in no set.
// a and b are read, never modified.
//
// ID:      ALG-HASH-006
// Time:    O(|a| + |b|) expected
// Space:   O(|b|) extra for the set of b; O(min(|a|, |b|)) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary operations)
// Source:  Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §6.4 (hashing)
// Tags:    hash, set, intersection, common, distinct, deterministic
func SetIntersection[T comparable](a, b []T) []T {
	inB := make(map[T]struct{}, len(b))
	for _, v := range b {
		inB[v] = struct{}{}
	}
	out := []T{}
	for _, v := range a {
		if _, ok := inB[v]; ok {
			out = append(out, v)
			// Deleting the emitted value removes later duplicates of it from a, with no seen-set.
			delete(inB, v)
		}
	}
	return out
}
