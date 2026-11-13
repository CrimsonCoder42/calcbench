package combinatorics

import "fmt"

// CombinationsWithRepetition returns C(n+k−1, k) exactly: the number of size-k multisets drawn from n
// kinds (k unlabelled balls into n labelled boxes). With n = 0 the count is 1 for k = 0 and 0 otherwise.
// The sum n+k−1 is formed in uint64, so it never overflows an int.
// Errors: ErrNegative when n < 0 or k < 0; ErrOverflow when the count > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-007
// Time:    O(m log n)   m = min(k, n−1, 64), as Binomial on (n+k−1, k)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.17 Table 26.17.1 (unlabelled balls into labelled boxes, unrestricted)
// Tags:    combinatorics, combinations, repetition, multiset, counting
// Units:   n: number of kinds; k: multiset size; result: a count of multisets
// Example: CombinationsWithRepetition(5, 3) = C(7,3) = 35
func CombinationsWithRepetition(n, k int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return 0, fmt.Errorf("%w: k", ErrNegative)
	}
	if n == 0 {
		if k == 0 {
			return 1, nil
		}
		return 0, nil
	}
	r, ok := binomialU64(uint64(n)+uint64(k)-1, uint64(k))
	if !ok {
		return 0, fmt.Errorf("%w: C(n+k-1,k)", ErrOverflow)
	}
	return r, nil
}
