package stack

import "cmp"

// DailyTemperatures appends to dst, for each reading i, how many positions later the first strictly
// greater reading comes (j − i for NextGreater's j), or 0 when none does, and returns the extended slice.
// It is NextGreater with each index turned into a distance.
// Units: readings in any one unit (compared, never combined); the result counts positions, which are
// days when there is one reading per day. readings is read, never modified.
//
// ID:      ALG-STACK-007
// Time:    O(n) amortized over the call   n = len(readings)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Berkman, Schieber, Vishkin, "Optimal doubly logarithmic parallel algorithms based on finding all nearest smaller values", Journal of Algorithms, 1993, 14(3):344–370, doi:10.1006/jagm.1993.1018
// Tags:    stack, monotonic-stack, next-greater, wait, days-until-warmer
func DailyTemperatures[T cmp.Ordered](dst []int, readings []T) []int {
	base := len(dst)
	out := NextGreater(dst, readings)
	for i, j := range out[base:] {
		if j == -1 {
			out[base+i] = 0
		} else {
			out[base+i] = j - i
		}
	}
	return out
}
