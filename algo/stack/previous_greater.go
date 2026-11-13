package stack

import "cmp"

// PreviousGreater appends to dst, for each i, the greatest index j < i with values[j] > values[i], or -1 when
// there is none, and returns the extended slice (len(values) indices). values is read, never modified.
//
// ID:      ALG-STACK-003
// Time:    O(n) amortized over the call   n = len(values)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Berkman, Schieber, Vishkin, "Optimal doubly logarithmic parallel algorithms based on finding all nearest smaller values", Journal of Algorithms, 1993, 14(3):344–370, doi:10.1006/jagm.1993.1018
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (a stack as a linked list)
// Tags:    stack, monotonic-stack, previous-greater, nearest-greater, ansv
func PreviousGreater[T cmp.Ordered](dst []int, values []T) []int {
	return appendNearest(dst, values, false, func(candidate, v T) bool { return cmp.Less(v, candidate) })
}
