package twopointer

// MoveZeroesInPlace moves every element equal to T's zero value to the end of values, keeping the
// other elements in their original relative order, and returns how many non-zero elements there are.
// Elements are swapped, never overwritten, so values keeps the same multiset: for floats, -0 == 0 counts
// as zero and keeps its sign bit; a NaN is not zero. The zeros at the end are in no promised order.
// values IS MODIFIED.
//
// ID:      ALG-TWOPTR-006
// Time:    O(n)   n = len(values); n comparisons, at most n swaps
// Space:   O(1), no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §7.1, procedure PARTITION (Lomuto's scheme: one pass, a boundary index, swap each kept element to the boundary; here the test is "not zero" and there is no pivot)
// Tags:    two-pointers, partition, stable-partition, move-zeroes, in-place
func MoveZeroesInPlace[T comparable](values []T) int {
	var zero T
	w := 0 // values[:w] are the non-zero elements seen so far, in order
	for r, v := range values {
		if v != zero {
			if r != w {
				values[w], values[r] = values[r], values[w]
			}
			w++
		}
	}
	return w
}
