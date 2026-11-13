package enumerate

import (
	"cmp"
	"iter"
	"slices"
)

// SubsetsUnique returns an iterator over the distinct sub-multisets of the multiset values, each yielded
// in ascending order (cmp.Compare). With distinct values d₁ < … < d_r of multiplicities m₁ … m_r there are
// ∏(mⱼ+1) of them, generated as mixed-radix counters with d₁ the least significant digit: for [2 1 2],
// {} 1 2 12 22 122. The yielded slice is reused between steps: copy it to keep it, never modify it.
// values is copied when iteration starts and never modified.
//
// ID:      ALG-ENUM-006
// Time:    O(n log n) to sort, then O(n) per sub-multiset; ∏(mⱼ+1) of them   n = len(values), mⱼ = multiplicity of the j-th distinct value
// Space:   O(n) extra, allocated once per iteration
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.1, Algorithm M (mixed-radix generation, radix mⱼ+1 for the j-th distinct value)
// Tags:    subset, multiset, duplicates, mixed-radix, enumeration, iterator
func SubsetsUnique[T cmp.Ordered](values []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		sorted := slices.Clone(values)
		slices.Sort(sorted)
		// distinct[j] and limit[j] = its multiplicity; count[j] = how many copies are in the subset.
		var distinct []T
		var limit []int
		for i, v := range sorted {
			if i == 0 || cmp.Compare(v, sorted[i-1]) != 0 {
				distinct = append(distinct, v)
				limit = append(limit, 0)
			}
			limit[len(limit)-1]++
		}
		r := len(distinct)
		count := make([]int, r)
		buf := make([]T, 0, len(sorted))
		for {
			if !yield(buf) {
				return
			}
			j := 0
			for j < r && count[j] == limit[j] {
				count[j] = 0
				j++
			}
			if j == r {
				return
			}
			count[j]++
			buf = buf[:0]
			for d, c := range count {
				for range c {
					buf = append(buf, distinct[d])
				}
			}
		}
	}
}
