package enumerate

import "iter"

// Combinations returns an iterator over the k-element combinations of values, chosen by position and
// kept in input order, in lexicographic order of the positions: for [a b c d] and k = 2, ab ac ad bc bd
// cd. k = 0 yields one empty combination; k < 0 or k > len(values) yields nothing (there are none).
// Equal values are not merged. The yielded slice is reused between steps: copy it to keep it, never
// modify it. values is read during iteration and must not change while iterating; it is never modified.
//
// ID:      ALG-ENUM-004
// Time:    O(1) amortised per combination for the index step, O(k) worst step; C(n,k) combinations   n = len(values)
// Space:   O(k) extra: the yielded buffer and the index array, allocated once per iteration
// Source:  Kreher, Stinson, Combinatorial Algorithms: Generation, Enumeration, and Search, CRC Press, 1999, §2.3.1 (lexicographic ordering of k-subsets, the successor step)
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.3 (generating all combinations)
// Tags:    combination, k-subset, enumeration, iterator, lexicographic
func Combinations[T any](values []T, k int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(values)
		if k < 0 || k > n {
			return
		}
		c := make([]int, k)
		buf := make([]T, k)
		for i := range c {
			c[i] = i
			buf[i] = values[i]
		}
		for {
			if !yield(buf) {
				return
			}
			// The rightmost position that can still move right.
			i := k - 1
			for i >= 0 && c[i] == n-k+i {
				i--
			}
			if i < 0 {
				return
			}
			c[i]++
			buf[i] = values[c[i]]
			for j := i + 1; j < k; j++ {
				c[j] = c[j-1] + 1
				buf[j] = values[c[j]]
			}
		}
	}
}
