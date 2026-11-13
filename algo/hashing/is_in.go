package hashing

// IsIn reports, for each element of values, whether it occurs in test: found[i] is true when values[i]
// == test[j] for some j. A NaN equals nothing, so it is never found; -0 and +0 are equal. The result has
// len(values) entries (empty and non-nil for empty values). values and test are read, never modified.
//
// ID:      ALG-HASH-016
// Time:    O(n + m) expected   n = len(values), m = len(test)
// Space:   O(m) extra (the set of test); O(n) for the result
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.4 (hashing: search in O(1) expected probes)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary operations)
// Tags:    hash, set, membership, isin, contains, lookup, mask
func IsIn[T comparable](values, test []T) []bool {
	set := make(map[T]struct{}, len(test))
	for _, v := range test {
		set[v] = struct{}{}
	}
	found := make([]bool, len(values))
	for i, v := range values {
		_, found[i] = set[v]
	}
	return found
}
