package hashing

// SetSymmetricDifference returns the distinct values that occur in exactly one of a and b: first those of
// a, in order of first appearance in a, then those of b, in order of first appearance in b. A NaN equals
// nothing, so every NaN of a and of b appears in the result. An empty result is a non-nil, empty slice.
// a and b are read, never modified.
//
// ID:      ALG-HASH-017
// Time:    O(|a| + |b|) expected
// Space:   O(|a| + |b|) extra for the two sets; O(|a| + |b|) for the result
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.4 (hashing)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary operations)
// Tags:    hash, set, symmetric-difference, xor, setxor, distinct, deterministic
func SetSymmetricDifference[T comparable](a, b []T) []T {
	inA := make(map[T]struct{}, len(a))
	for _, v := range a {
		inA[v] = struct{}{}
	}
	inB := make(map[T]struct{}, len(b))
	for _, v := range b {
		inB[v] = struct{}{}
	}
	out := []T{}
	emitted := make(map[T]struct{})
	for _, side := range [2]struct {
		values []T
		other  map[T]struct{}
	}{{a, inB}, {b, inA}} {
		for _, v := range side.values {
			if _, ok := side.other[v]; ok {
				continue
			}
			if _, ok := emitted[v]; ok {
				continue
			}
			emitted[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
