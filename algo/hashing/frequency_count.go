package hashing

// FrequencyCount returns how many times each distinct value occurs in values, as a map from value to
// count. Empty input gives an empty, non-nil map. Each NaN is its own key with count 1.
// values is read, never modified.
//
// ID:      ALG-HASH-001
// Time:    O(n) expected   n = len(values)
// Space:   O(m) for the result   m = number of distinct values
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (hash tables; Theorem 11.2)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §3.1 (client FrequencyCounter) and §3.4
// Tags:    hash, map, count, frequency, histogram, tally
func FrequencyCount[T comparable](values []T) map[T]int {
	counts := make(map[T]int)
	for _, v := range values {
		counts[v]++
	}
	return counts
}
