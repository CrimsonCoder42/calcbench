package dp

// PalindromePartitionMinCuts returns the fewest cuts that split seq into contiguous pieces each of
// which is a palindrome (a minimum palindromic factorisation has this many plus one factors). Empty
// input and a palindrome need 0 cuts. Elements match by ==. seq is read, never modified.
//
// ID:      ALG-DP-017
// Time:    O(n²) worst case   n = len(seq); O(n + P) where P = the number of palindromic runs found by expanding centres
// Space:   O(n) extra (one table of n+1 ints, one allocation)
// Source:  Fici, Gagie, Kärkkäinen, Kempa, "A subquadratic algorithm for minimum palindromic factorization", Journal of Discrete Algorithms, 2014, 28:41–48, doi:10.1016/j.jda.2014.08.001, §2 "A simple quadratic algorithm" (fewest factors of seq[:j] = 1 + min over palindromes seq[i:j] of fewest factors of seq[:i]); their O(n log n) algorithm is not built here
// Tags:    dp, palindrome, partition, factorization, string, min-cuts
func PalindromePartitionMinCuts[T comparable](seq []T) int {
	n := len(seq)
	if n == 0 {
		return 0
	}
	// cuts[j] = fewest cuts for seq[:j]; cuts[0] = −1 so a palindromic prefix costs 0 cuts.
	cuts := make([]int, n+1)
	for j := range cuts {
		cuts[j] = j - 1
	}
	// Centres in increasing order: every palindrome ending at i−1 has a smaller centre than one starting
	// at i, so cuts[lo] is final before it is read.
	for c := 0; c < 2*n-1; c++ {
		for lo, hi := c/2, c/2+c%2; lo >= 0 && hi < n && seq[lo] == seq[hi]; lo, hi = lo-1, hi+1 {
			cuts[hi+1] = min(cuts[hi+1], cuts[lo]+1)
		}
	}
	return cuts[n]
}
