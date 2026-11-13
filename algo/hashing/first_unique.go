package hashing

// FirstUnique returns the index of the first value that occurs exactly once in values, or -1 when every
// value repeats or values is empty. A NaN is never equal to anything, so a NaN is always unique.
// values is read, never modified.
//
// ID:      ALG-HASH-005
// Time:    O(n) expected   n = len(values); two passes
// Space:   O(m) extra   m = number of distinct values
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (counting in a hash table)
// Tags:    hash, map, unique, first-unique, non-repeating, count
func FirstUnique[T comparable](values []T) int {
	counts := FrequencyCount(values)
	for i, v := range values {
		// v != v only for NaN: its count is stored under a key no lookup can find, and it is unique.
		if v != v || counts[v] == 1 {
			return i
		}
	}
	return -1
}
