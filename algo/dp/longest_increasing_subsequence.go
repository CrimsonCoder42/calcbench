package dp

import (
	"cmp"

	"github.com/CrimsonCoder42/calcbench/algo/search"
)

// LongestIncreasingSubsequence returns the indices, ascending, of one longest strictly increasing
// subsequence of values; its length is len of the result. Of the longest ones it returns the one that
// patience sorting finds by following back-pointers from the top card of the last pile: its last
// element has the smallest value any longest one can end with (the latest index holding it), and each
// earlier element is the top of the previous pile at the moment the next one was dealt.
// "Increasing" is in cmp.Compare order (a NaN is below every other float). Empty input gives an empty,
// non-nil slice. values is read, never modified.
//
// ID:      ALG-DP-001
// Time:    O(n log n)   n = len(values); one binary search (search.LowerBound, ALG-SEARCH-001) per element
// Space:   O(n) extra (pile tops and back-pointers, two allocations); O(L) for the result, L = the length found
// Source:  Fredman, "On computing the length of longest increasing subsequences", Discrete Mathematics, 1975, 11(1):29–35, doi:10.1016/0012-365X(75)90103-X (O(n log n), and n log n − n log log n + O(n) comparisons are needed)
// Source:  Aldous, Diaconis, "Longest increasing subsequences: from patience sorting to the Baik–Deift–Johansson theorem", Bulletin of the AMS, 1999, 36(4):413–432, doi:10.1090/S0273-0979-99-00796-X, §1–2 (patience sorting: the number of piles is the LIS length)
// Tags:    dp, subsequence, lis, increasing, patience-sorting, binary-search
func LongestIncreasingSubsequence[T cmp.Ordered](values []T) []int {
	n := len(values)
	if n == 0 {
		return []int{}
	}
	// tops[p] is the top card of pile p: the smallest last value of an increasing subsequence of length
	// p+1 seen so far. tops is strictly increasing, so the pile for v is the first top >= v.
	tops := make([]T, 0, n)
	idx := make([]int, 2*n)
	topIndex, prev := idx[:n:n], idx[n:] // topIndex[p] = index of tops[p]; prev[i] = back-pointer of i
	for i, v := range values {
		p := search.LowerBound(tops, v)
		if p == len(tops) {
			tops = append(tops, v)
		} else {
			tops[p] = v
		}
		topIndex[p] = i
		if p > 0 {
			prev[i] = topIndex[p-1]
		} else {
			prev[i] = -1
		}
	}
	out := make([]int, len(tops))
	for k, i := len(out)-1, topIndex[len(tops)-1]; k >= 0; k-- {
		out[k] = i
		i = prev[i]
	}
	return out
}
