package hashing

// TwoSum returns indices i < j with values[i] + values[j] == target, or (-1, -1) when no pair exists.
// Of all such pairs it returns the one with the smallest j and, for that j, the smallest i. The sum is
// exact: a pair whose true sum overflows the type never matches, and target − values[j] is
// overflow-checked. values is read, never modified.
//
// ID:      ALG-HASH-012
// Time:    O(n) expected   n = len(values)
// Space:   O(n) extra for the map of value to first index
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 2.3-7 (the pair-sum problem) and ch. 11 (hash tables)
// Tags:    hash, map, pair-sum, two-sum, complement, lookup
func TwoSum[T Integer](values []T, target T) (i, j int) {
	firstIndex := make(map[T]int, len(values))
	for k, v := range values {
		// need = target - v, exact unless it wraps; a wrapped need has no partner in the type.
		need := target - v
		if !((v > 0 && need > target) || (v < 0 && need < target)) {
			if first, ok := firstIndex[need]; ok {
				return first, k
			}
		}
		if _, ok := firstIndex[v]; !ok {
			firstIndex[v] = k
		}
	}
	return -1, -1
}
