package enumerate

import "iter"

// Subsets returns an iterator over all 2ⁿ subsets of values (by position, kept in input order), in
// binary counting order with values[0] the least significant digit: for [a b c], {} a b ab c ac bc abc.
// Equal values are not merged; use SubsetsUnique for distinct ones. The yielded slice is reused between
// steps: copy it to keep it, never modify it. values is read during iteration and must not change while
// iterating; it is never modified. There is no limit on n other than time.
//
// ID:      ALG-ENUM-005
// Time:    O(n) per subset (rebuilding the yielded slice), O(n·2ⁿ) in all   n = len(values)
// Space:   O(n) extra: the yielded buffer and the membership flags, allocated once per iteration
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.1, Algorithm M (mixed-radix generation) with every radix 2
// Tags:    subset, power-set, binary-counting, enumeration, iterator
func Subsets[T any](values []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(values)
		in := make([]bool, n)
		buf := make([]T, 0, n)
		for {
			if !yield(buf) {
				return
			}
			// Add one: clear the trailing set flags, set the first clear one.
			i := 0
			for i < n && in[i] {
				in[i] = false
				i++
			}
			if i == n {
				return
			}
			in[i] = true
			buf = buf[:0]
			for j, v := range values {
				if in[j] {
					buf = append(buf, v)
				}
			}
		}
	}
}
