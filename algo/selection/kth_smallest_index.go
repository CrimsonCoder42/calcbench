package selection

import "cmp"

// KthSmallestIndex returns the index in values of the k-th smallest element (k = 1 being the minimum) and
// an index partition around it (an argpartition): perm is a permutation of 0…n-1 in which the indexes of
// the elements less than the k-th come first, then those equal to it, then those greater, each block in
// ascending index order (so the result is deterministic and stable). perm[k-1] is index, and perm[:k]
// are the indexes of k smallest elements. perm is written into dst, reused when cap(dst) >= n, else one
// new slice. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrRank when k < 1 or k > len(values). On error perm is nil.
//
// ID:      ALG-SELECT-005
// Time:    O(n) worst case   n = len(values) (introselect for the value, then one counting pass and one placing pass)
// Space:   O(n) extra (one allocation, a working copy of values for the selection); O(n) for perm
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §9.2–9.3 (selection in expected and worst-case linear time)
// Source:  Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience, 1997, 27(8):983–993, doi:10.1002/(SICI)1097-024X(199708)27:8<983::AID-SPE117>3.0.CO;2-#
// Tags:    selection, order-statistic, argpartition, index, introselect, stable-partition, dst
func KthSmallestIndex[T cmp.Ordered](dst []int, values []T, k int) (index int, perm []int, err error) {
	kth, _, err := KthSmallestInto(nil, values, k)
	if err != nil {
		return 0, nil, err
	}
	n := len(values)
	var less, equal int
	for _, v := range values {
		switch cmp.Compare(v, kth) {
		case -1:
			less++
		case 0:
			equal++
		}
	}
	if cap(dst) >= n {
		perm = dst[:n]
	} else {
		perm = make([]int, n)
	}
	nextLess, nextEqual, nextGreater := 0, less, less+equal
	for i, v := range values {
		switch cmp.Compare(v, kth) {
		case -1:
			perm[nextLess] = i
			nextLess++
		case 0:
			perm[nextEqual] = i
			nextEqual++
		default:
			perm[nextGreater] = i
			nextGreater++
		}
	}
	return perm[k-1], perm, nil
}
