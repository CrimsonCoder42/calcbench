package prefix

// EquilibriumIndex returns the first index i at which the sum of the values before i equals the sum of
// the values after i (values[:i] and values[i+1:]; an empty side sums to 0), or -1 when there is none.
// The left sum is a running prefix sum and the right sum is the total minus the left sum minus
// values[i]; both are kept in 128-bit arithmetic, so no input overflows.
// values is read, never modified.
//
// ID:      ALG-PREFIX-011
// Time:    O(n)   n = len(values); two passes
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.1 "Static array queries", "Sum queries", p. 84 (a range sum is a difference of prefix sums), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8 "Algorithm Design Techniques", §8.2 (the cumulative array cumarr)
// Tags:    prefix-sum, equilibrium, pivot-index, balance, integer
func EquilibriumIndex[T integer](values []T) int {
	var total int128
	for _, v := range values {
		total = total.add(widen(v))
	}
	var left int128
	for i, v := range values {
		w := widen(v)
		if left == total.sub(left).sub(w) {
			return i
		}
		left = left.add(w)
	}
	return -1
}
