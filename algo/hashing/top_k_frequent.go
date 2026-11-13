package hashing

import "cmp"

// TopKFrequent returns the k most frequent values of values with their counts, ranked by count
// descending and, for equal counts, by value ascending (cmp.Compare), or every distinct value when
// k >= the number of distinct values. k = 0 gives an empty, non-nil slice. values is read, never modified.
// It counts with FrequencyCount and ranks with TopKByCount.
// Errors: ErrNegativeK when k < 0.
//
// ID:      ALG-HASH-002
// Time:    O(n + m log k) expected   n = len(values), m = number of distinct values, k = min(k, m)
// Space:   O(m) extra for the counts; O(k) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (counting in a hash table)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4 (client TopM: a min-heap of the M largest)
// Tags:    hash, heap, top-k, frequency, mode, most-common, deterministic
func TopKFrequent[T cmp.Ordered](values []T, k int) ([]Count[T], error) {
	if k < 0 {
		return TopKByCount[T](nil, k)
	}
	return TopKByCount(FrequencyCount(values), k)
}
