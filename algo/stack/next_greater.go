package stack

import (
	"cmp"
	"slices"
)

// NextGreater appends to dst, for each i, the least index j > i with values[j] > values[i], or -1 when
// there is none, and returns the extended slice (len(values) indices). values is read, never modified.
//
// ID:      ALG-STACK-001
// Time:    O(n) amortized over the call   n = len(values)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Berkman, Schieber, Vishkin, "Optimal doubly logarithmic parallel algorithms based on finding all nearest smaller values", Journal of Algorithms, 1993, 14(3):344–370, doi:10.1006/jagm.1993.1018
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (a stack as a linked list)
// Tags:    stack, monotonic-stack, next-greater, nearest-greater, ansv
func NextGreater[T cmp.Ordered](dst []int, values []T) []int {
	return appendNearest(dst, values, true, func(candidate, v T) bool { return cmp.Less(v, candidate) })
}

// appendNearest appends, for each i, the nearest index j on one side of i (after i when forward, before
// it otherwise) with qualifies(values[j], values[i]), or -1.
//
// The monotonic stack is kept as a linked list through the result: scanning from the far end, the
// candidates for i are its neighbour, then the neighbour's answer, then that one's answer, and so on.
// When a candidate fails, so does everything it skipped (qualifies is a strict order relation), so
// each index is passed over at most once in total: O(n) amortized and no stack to allocate.
func appendNearest[T cmp.Ordered](dst []int, values []T, forward bool, qualifies func(candidate, v T) bool) []int {
	n := len(values)
	base := len(dst)
	out := slices.Grow(dst, n)[:base+n]
	res := out[base:]
	if forward {
		for i := n - 1; i >= 0; i-- {
			j := i + 1
			if j == n {
				j = -1
			}
			for j != -1 && !qualifies(values[j], values[i]) {
				j = res[j]
			}
			res[i] = j
		}
	} else {
		for i := range n {
			j := i - 1
			for j != -1 && !qualifies(values[j], values[i]) {
				j = res[j]
			}
			res[i] = j
		}
	}
	return out
}
